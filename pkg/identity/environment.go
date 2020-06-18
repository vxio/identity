package identity

import (
	"context"
	"database/sql"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/moov-io/identity" // need to import the embedded files
	"github.com/moov-io/tumbler/pkg/jwe"

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/authn"
	configpkg "github.com/moov-io/identity/pkg/config"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	"github.com/moov-io/tumbler/pkg/webkeys"
)

// Environment - Contains everything thats been instantiated for this service.
type Environment struct {
	Logger      logging.Logger
	Config      *Config
	TimeService *stime.TimeService

	AuthnKeys   webkeys.WebKeysService
	SessionKeys webkeys.WebKeysService

	InviteService      api.InvitesApiServicer
	IdentitiesService  *identities.Service
	CredentialsService *credentials.CredentialsService

	PublicRouter *mux.Router

	Shutdown func()
}

// NewEnvironment - Generates a new default environment. Overrides can be specified via configs.
func NewEnvironment(env *Environment) (*Environment, error) {
	if env.Logger == nil {
		env.Logger = logging.NewDefaultLogger()
	}

	if env.Config == nil {
		ConfigService := configpkg.NewConfigService(env.Logger)

		global := &GlobalConfig{}
		if err := ConfigService.Load(global); err != nil {
			return nil, err
		}

		env.Config = &global.Identity
	}

	//db setup
	db, close, err := initializeDatabase(env.Logger, env.Config.Database)
	if err != nil {
		close()
		return nil, err
	}

	if env.TimeService == nil {
		t := stime.NewSystemTimeService()
		env.TimeService = &t
	}

	if env.AuthnKeys == nil {
		ak, err := webkeys.NewWebKeysService(env.Logger, env.Config.Authentication.Keys)
		if err != nil {
			return nil, env.Logger.Fatal().LogErrorF("Unable to load up the Authentication JSON Web Key Set - %w", err)
		}
		env.AuthnKeys = ak
	}

	AuthnTokenService := jwe.NewJWEService(*env.TimeService, time.Second, env.AuthnKeys)

	if env.SessionKeys == nil {
		SessionKeys, err := webkeys.NewWebKeysService(env.Logger, env.Config.Session.Keys)
		if err != nil {
			return nil, env.Logger.Fatal().LogErrorF("Unable to load up up the Session JSON Web Key Set - %w", err)
		}
		env.SessionKeys = SessionKeys
	}

	SessionJwe := jwe.NewJWEService(*env.TimeService, env.Config.Session.Expiration, env.SessionKeys)

	SessionService := session.NewSessionService(*env.TimeService, SessionJwe, env.Config.Session)

	templateService, err := notifications.NewTemplateRepository(env.Logger)
	if err != nil {
		return nil, err
	}

	NotificationsService, err := notifications.NewNotificationsService(env.Logger, env.Config.Notifications, templateService)
	if err != nil {
		return nil, err
	}

	if env.IdentitiesService == nil {
		IdentityRepository := identities.NewIdentityRepository(db)
		IdentitiesService := identities.NewIdentitiesService(*env.TimeService, IdentityRepository)
		env.IdentitiesService = IdentitiesService
	}

	if env.CredentialsService == nil {
		CredentialRepository := credentials.NewCredentialRepository(db)
		CredentialsService := credentials.NewCredentialsService(*env.TimeService, CredentialRepository)
		env.CredentialsService = CredentialsService
	}

	if env.InviteService == nil {
		InvitesRepository := invites.NewInvitesRepository(db)
		InvitesService, err := invites.NewInvitesService(env.Config.Invites, *env.TimeService, InvitesRepository, NotificationsService)
		if err != nil {
			return nil, err
		}

		env.InviteService = InvitesService
	}

	AuthnService := authn.NewAuthnService(env.Logger, *env.CredentialsService, *env.IdentitiesService, SessionService, env.InviteService, env.Config.Authentication.LandingURL)

	// router
	if env.PublicRouter == nil {
		env.PublicRouter = mux.NewRouter()
	}

	// public endpoint
	jwksController := webkeys.NewJWKSController(env.SessionKeys)
	jwksRouter := env.PublicRouter.NewRoute().Subrouter()
	jwksController.AppendRoutes(jwksRouter)

	// authn endpoints

	AuthnMiddleware, err := authn.NewMiddleware(env.Logger, *env.TimeService, AuthnTokenService)
	if err != nil {
		return nil, env.Logger.Fatal().LogErrorF("Can't startup the Authn middleware - %w", err)
	}

	AuthnController := authn.NewAuthnAPIController(env.Logger, AuthnService)

	authnRouter := env.PublicRouter.NewRoute().Subrouter()
	authnRouter = api.AppendRouters(env.Logger, authnRouter, AuthnController)
	authnRouter.Use(AuthnMiddleware.Handler)

	// authed server

	// auth middleware for the tokens coming from the gateway
	GatewayMiddleware, err := tmw.NewTumblerMiddlewareFromConfig(env.Logger, *env.TimeService, env.Config.Gateway)
	if err != nil {
		return nil, env.Logger.Fatal().LogErrorF("Can't startup the Gateway middleware - %w", err)
	}

	WhoAmIController := session.NewWhoAmIController(env.Logger, SessionService, *env.IdentitiesService)
	IdentitiesController := identities.NewIdentitiesController(env.IdentitiesService)
	CredentialsController := credentials.NewCredentialsApiController(env.CredentialsService)
	InvitesController := invites.NewInvitesController(env.InviteService)

	authedRouter := env.PublicRouter.NewRoute().Subrouter()
	authedRouter = api.AppendRouters(env.Logger, authedRouter, IdentitiesController, CredentialsController, InvitesController, WhoAmIController)
	authedRouter.Use(GatewayMiddleware.Handler)

	env.Shutdown = func() {
		close()
	}

	return env, nil
}

func initializeDatabase(logger logging.Logger, config database.DatabaseConfig) (*sql.DB, func(), error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// migrate database
	db, err := database.New(ctx, logger, config)
	if err != nil {
		return nil, cancelFunc, logger.Fatal().LogError("Error creating database", err)
	}

	shutdown := func() {
		logger.Info().Log("Shutting down the db")
		cancelFunc()
		if err := db.Close(); err != nil {
			logger.Fatal().LogError("Error closing DB", err)
		}
	}

	if err := database.RunMigrations(logger, db, config); err != nil {
		return nil, shutdown, logger.Fatal().LogError("Error running migrations", err)
	}

	logger.Info().Log("finished initializing db")

	return db, shutdown, err
}
