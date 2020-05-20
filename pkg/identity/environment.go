package identity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gorilla/mux"
	_ "github.com/moov-io/identity" // need to import the embedded files

	"github.com/go-kit/kit/log"
	api "github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/authn"
	configpkg "github.com/moov-io/identity/pkg/config"
	"github.com/moov-io/identity/pkg/credentials"
	"github.com/moov-io/identity/pkg/database"
	"github.com/moov-io/identity/pkg/identities"
	"github.com/moov-io/identity/pkg/invites"
	"github.com/moov-io/identity/pkg/notifications"
	"github.com/moov-io/identity/pkg/stime"
	"github.com/moov-io/identity/pkg/webkeys"
	"github.com/moov-io/identity/pkg/zerotrust"
)

type Environment struct {
	Logger log.Logger
	Config IdentityConfig

	InviteService      api.InvitesApiServicer
	IdentitiesService  identities.IdentitiesService
	CredentialsService credentials.CredentialsService

	PublicRouter mux.Router

	Shutdown func()
}

func NewEnvironment(logger log.Logger, configOverride *IdentityConfig) (*Environment, error) {
	var config *IdentityConfig
	if configOverride != nil {
		config = configOverride
	} else {
		ConfigService := configpkg.NewConfigService(logger)

		config = &IdentityConfig{}
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

	AuthnPublicKeys, err := webkeys.NewWebKeysService(logger, config.Keys.AuthnPublic)
	if err != nil {
		logger.Log("main", "Unable to load up the Authentication JSON Web Key Set")
		return nil, err
	}

	GatewayPublicKeys, err := webkeys.NewWebKeysService(logger, config.Keys.GatewayPublic)
	if err != nil {
		logger.Log("main", "Unable to load up the Gateway JSON Web Key Set")
		return nil, err
	}

	SessionPublicKeys, err := webkeys.NewWebKeysService(logger, config.Keys.SessionPublic)
	if err != nil {
		logger.Log("main", "Unable to load up up the Session Public JSON Web Key Set")
		return nil, err
	}

	SessionPrivateKeys, err := webkeys.NewWebKeysService(logger, config.Keys.SessionPrivate)
	if err != nil {
		logger.Log("main", "Unable to load up up the Session Private JSON Web Key Set")
		return nil, err
	}

	SessionService := authn.NewSessionService(TimeService, SessionPrivateKeys, config.Session)

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

	AuthnService := authn.NewAuthnService(*CredentialsService, *IdentitiesService, SessionService, InvitesService, config.Session.LandingURL)

	// router
	router := mux.NewRouter()

	// public endpoint
	jwksController := webkeys.NewJWKSController(SessionPublicKeys)
	jwksRouter := router.NewRoute().Subrouter()
	jwksRouter = jwksController.AppendRoutes(jwksRouter)

	// authn endpoints

	AuthnMiddleware, err := authn.NewAuthnMiddleware(TimeService, AuthnPublicKeys)
	if err != nil {
		logger.Log("main", fmt.Sprintf("Can't startup the Authn middleware - %s", err))
		return nil, err
	}

	AuthnController := authn.NewAuthnAPIController(logger, AuthnService)

	authnRouter := router.NewRoute().Subrouter()
	authnRouter = api.AppendRouters(authnRouter, AuthnController)
	authnRouter.Use(AuthnMiddleware.Handler)

	// authed server

	// auth middleware for the tokens coming from the gateway
	GatewayMiddleware, err := zerotrust.NewJWTMiddleware(GatewayPublicKeys)
	if err != nil {
		logger.Log("main", fmt.Sprintf("Can't startup the Gateway middleware - %s", err))
		return nil, err
	}

	WhoAmIController := authn.NewWhoAmIController()
	IdentitiesController := identities.NewIdentitiesApiController(IdentitiesService)
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

func initializeDatabase(logger log.Logger, config database.DatabaseConfig) (*sql.DB, func(), error) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// migrate database
	db, err := database.New(ctx, logger, config)
	if err != nil {
		msg := fmt.Sprintf("error creating database: %v", err)
		logger.Log("msg", msg)
		return nil, func() {}, err
	}

	shutdown := func() {
		logger.Log("msg", "Shutting down the db")
		cancelFunc()
		if err := db.Close(); err != nil {
			logger.Log("exit", err)
		}
	}

	if err := database.RunMigrations(db, config); err != nil {
		msg := fmt.Sprintf("Error running migrations: %s", err)
		logger.Log("msg", msg)
		return nil, shutdown, err
	}

	logger.Log("msg", "finished....")

	return db, shutdown, err
}
