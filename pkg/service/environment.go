package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/authn"
	"github.com/moov-io/identity/pkg/config"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/tumbler/pkg/jwe"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
	"github.com/moov-io/tumbler/pkg/webkeys"
)

// Environment - Contains everything thats been instantiated for this service.
type Environment struct {
	Logger       logging.Logger
	Config       *Config
	TimeService  stime.TimeService
	GatewayKeys  webkeys.WebKeysService
	PublicRouter *mux.Router
	Shutdown     func()

	InviteService      api.InvitesApiServicer
	IdentitiesService  identities.Service
	CredentialsService credentials.CredentialsService
}

// NewEnvironment - Generates a new default environment. Overrides can be specified via configs.
func NewEnvironment(env *Environment) (*Environment, error) {
	if env == nil {
		env = &Environment{}
	}

	if env.Logger == nil {
		env.Logger = logging.NewDefaultLogger()
	}

	if env.Config == nil {
		ConfigService := config.NewConfigService(env.Logger)

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
		env.TimeService = stime.NewSystemTimeService()
	}

	AuthnKeys, err := webkeys.NewWebKeysService(env.Logger, env.Config.Authentication.Keys)
	if err != nil {
		return nil, env.Logger.Fatal().LogErrorF("Unable to load up the Authentication JSON Web Key Set - %w", err)
	}

	AuthnTokenService := jwe.NewJWEService(env.TimeService, time.Second, AuthnKeys)

	SessionKeys, err := webkeys.NewWebKeysService(env.Logger, env.Config.Session.Keys)
	if err != nil {
		return nil, env.Logger.Fatal().LogErrorF("Unable to load up up the Session JSON Web Key Set - %w", err)
	}

	SessionJwe := jwe.NewJWEService(env.TimeService, env.Config.Session.Expiration, SessionKeys)

	SessionService := session.NewSessionService(env.TimeService, SessionJwe, env.Config.Session)

	templateService, err := notifications.NewTemplateRepository(env.Logger)
	if err != nil {
		return nil, err
	}

	NotificationsService, err := notifications.NewNotificationsService(env.Logger, env.Config.Notifications, templateService)
	if err != nil {
		return nil, err
	}

	IdentityRepository := identities.NewIdentityRepository(db)
	IdentitiesService := identities.NewIdentitiesService(env.TimeService, IdentityRepository)

	CredentialRepository := credentials.NewCredentialRepository(db)
	CredentialsService := credentials.NewCredentialsService(env.TimeService, CredentialRepository)

	AuthnClient, err := authn.NewAuthnClient(env.Logger, env.Config.Services.Authn)
	if err != nil {
		return nil, err
	}

	InvitesRepository := invites.NewInvitesRepository(db)
	InvitesService, err := invites.NewInvitesService(env.Config.Invites, env.TimeService, InvitesRepository, NotificationsService, AuthnClient, IdentitiesService)
	if err != nil {
		return nil, err
	}

	AuthnService := authn.NewAuthnService(env.Logger, *CredentialsService, IdentitiesService, SessionService, InvitesService)

	// router
	if env.PublicRouter == nil {
		env.PublicRouter = mux.NewRouter()
	}

	// public endpoint
	jwksController := webkeys.NewJWKSController(SessionKeys)
	jwksRouter := env.PublicRouter.NewRoute().Subrouter()
	jwksController.AppendRoutes(jwksRouter)

	// authn endpoints

	AuthnMiddleware, err := authn.NewMiddleware(env.Logger, env.TimeService, AuthnTokenService)
	if err != nil {
		return nil, env.Logger.Fatal().LogErrorF("Can't startup the Authn middleware - %w", err)
	}

	AuthnController := authn.NewAuthnAPIController(env.Logger, AuthnService)

	authnRouter := env.PublicRouter.NewRoute().Subrouter()
	authnRouter = api.AppendRouters(env.Logger, authnRouter, AuthnController)
	authnRouter.Use(AuthnMiddleware.Handler)

	// auth middleware for the tokens coming from the gateway
	GatewayMiddleware, err := tmw.NewServerFromConfig(env.Logger, env.TimeService, env.Config.Gateway)
	if err != nil {
		return nil, env.Logger.Fatal().LogErrorF("Can't startup the Gateway middleware - %w", err)
	}

	SessionController := session.NewSessionController(env.Logger, IdentitiesService, env.TimeService)
	IdentitiesController := identities.NewIdentitiesController(IdentitiesService)
	CredentialsController := credentials.NewCredentialsApiController(CredentialsService)
	InvitesController := invites.NewInvitesController(env.Logger, InvitesService)

	authedRouter := env.PublicRouter.NewRoute().Subrouter()
	authedRouter = api.AppendRouters(env.Logger, authedRouter, IdentitiesController, CredentialsController, InvitesController)
	SessionController.AppendRoutes(authedRouter)
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
