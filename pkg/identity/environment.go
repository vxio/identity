package identity

import (
	"context"
	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/moov-io/identity" // need to import the embedded files

	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/authn"
	configpkg "github.com/moov-io/identity/pkg/config"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/logging"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/session"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
)

// Environment - Contains everything thats been instantiated for this service.
type Environment struct {
	Logger logging.Logger
	Config Config

	InviteService      api.InvitesApiServicer
	IdentitiesService  identities.Service
	CredentialsService credentials.CredentialsService

	PublicRouter mux.Router

	Shutdown func()
}

// NewEnvironment - Generates a new default environment. Overrides can be specified via configs.
func NewEnvironment(logger logging.Logger, configOverride *Config) (*Environment, error) {
	var config *Config
	if configOverride != nil {
		config = configOverride
	} else {
		ConfigService := configpkg.NewConfigService(logger)

		config = &Config{}
		if err := ConfigService.Load(config); err != nil {
			return nil, err
		}
	}

	//db setup
	db, close, err := initializeDatabase(logger, config.Database)
	if err != nil {
		close()
		return nil, err
	}

	TimeService := stime.NewSystemTimeService()

	AuthnPublicKeys, err := webkeys.NewWebKeysService(logger, config.Authentication.Keys)
	if err != nil {
		return nil, logger.Fatal().LogErrorF("Unable to load up the Authentication JSON Web Key Set - %w", err)
	}

	GatewayPublicKeys, err := webkeys.NewWebKeysService(logger, config.Gateway.Keys)
	if err != nil {
		return nil, logger.Fatal().LogErrorF("Unable to load up the Gateway JSON Web Key Set - %w", err)
	}

	SessionKeys, err := webkeys.NewWebKeysService(logger, config.Session.Keys)
	if err != nil {
		return nil, logger.Fatal().LogErrorF("Unable to load up up the Session JSON Web Key Set - %w", err)
	}

	SessionService := session.NewSessionService(TimeService, SessionKeys, config.Session)

	templateService, err := notifications.NewTemplateRepository(logger)
	if err != nil {
		return nil, err
	}

	NotificationsService, err := notifications.NewNotificationsService(logger, config.Notifications, templateService)
	if err != nil {
		return nil, err
	}

	IdentityRepository := identities.NewIdentityRepository(db)
	IdentitiesService := identities.NewIdentitiesService(TimeService, IdentityRepository)

	CredentialRepository := credentials.NewCredentialRepository(db)
	CredentialsService := credentials.NewCredentialsService(TimeService, CredentialRepository)

	InvitesRepository := invites.NewInvitesRepository(db)
	InvitesService, err := invites.NewInvitesService(config.Invites, TimeService, InvitesRepository, NotificationsService)
	if err != nil {
		return nil, err
	}

	AuthnService := authn.NewAuthnService(logger, *CredentialsService, *IdentitiesService, SessionService, InvitesService, config.Authentication.LandingURL)

	// router
	router := mux.NewRouter()

	// public endpoint
	jwksController := webkeys.NewJWKSController(SessionKeys)
	jwksRouter := router.NewRoute().Subrouter()
	jwksController.AppendRoutes(jwksRouter)

	// authn endpoints

	AuthnMiddleware, err := authn.NewMiddleware(logger, TimeService, AuthnPublicKeys)
	if err != nil {
		return nil, logger.Fatal().LogErrorF("Can't startup the Authn middleware - %w", err)
	}

	AuthnController := authn.NewAuthnAPIController(logger, AuthnService)

	authnRouter := router.NewRoute().Subrouter()
	authnRouter = api.AppendRouters(authnRouter, AuthnController)
	authnRouter.Use(AuthnMiddleware.Handler)

	// authed server

	// auth middleware for the tokens coming from the gateway
	GatewayMiddleware, err := gateway.NewMiddleware(logger, TimeService, GatewayPublicKeys)
	if err != nil {
		return nil, logger.Fatal().LogErrorF("Can't startup the Gateway middleware - %w", err)
	}

	WhoAmIController := session.NewWhoAmIController(logger, SessionService, *IdentitiesService)
	IdentitiesController := identities.NewIdentitiesController(IdentitiesService)
	CredentialsController := credentials.NewCredentialsApiController(CredentialsService)
	InvitesController := invites.NewInvitesController(InvitesService)

	authedRouter := router.NewRoute().Subrouter()
	authedRouter = api.AppendRouters(authedRouter, IdentitiesController, CredentialsController, InvitesController, WhoAmIController)
	authedRouter.Use(GatewayMiddleware.Handler)

	env := Environment{
		Logger: logger,
		Config: *config,

		InviteService:      InvitesService,
		IdentitiesService:  *IdentitiesService,
		CredentialsService: *CredentialsService,

		PublicRouter: *router,

		Shutdown: func() {
			close()
		},
	}

	return &env, nil
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
