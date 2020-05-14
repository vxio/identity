/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package identity

import (
	"context"
	"database/sql"
	"fmt"
	"os"

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

func NewEnvironment(configOverride *IdentityConfig) (*Environment, error) {

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

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
	db, close := initializeDatabase(logger, config.Database)

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

	SessionPrivateKeys, err := webkeys.NewWebKeysService(logger, config.Keys.SessionPrivate)
	if err != nil {
		logger.Log("main", "Unable to load up up the Session JSON Web Key Set")
		return nil, err
	}

	SessionService := authn.NewSessionService(TimeService, SessionPrivateKeys, config.Session)

	NotificationsService, err := notifications.NewNotificationsService(config.Notifications)
	if err != nil {
		return nil, err
	}

	IdentityRepository := identities.NewIdentityRepository(db)
	IdentitiesService := identities.NewIdentitiesService(TimeService, IdentityRepository)

	CredentialRepository := credentials.NewCredentialRepository(db)
	CredentialsService := credentials.NewCredentialsService(TimeService, CredentialRepository)

	InvitesRepository := invites.NewInvitesRepository(db)
	InvitesService := invites.NewInvitesService(config.Invites, TimeService, InvitesRepository, NotificationsService)

	AuthnService := authn.NewAuthnService(*CredentialsService, *IdentitiesService, SessionService)

	// router
	router := mux.NewRouter()

	// authn endpoints

	AuthnMiddleware, err := authn.NewAuthnMiddleware(TimeService, AuthnPublicKeys)
	if err != nil {
		logger.Log("main", fmt.Sprintf("Can't startup the Authn middleware - %s", err))
		return nil, err
	}

	AuthnController := authn.NewAuthnAPIController(AuthnService)

	authnRouter := router.NewRoute().Subrouter()
	authnRouter = api.AppendRouters(authnRouter, AuthnController)
	authnRouter.Use(AuthnMiddleware.Handler)

	// public server

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

	publicRouter := router.NewRoute().Subrouter()
	publicRouter = api.AppendRouters(publicRouter, IdentitiesController, CredentialsController, InvitesController, WhoAmIController)
	publicRouter.Use(GatewayMiddleware.Handler)

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

func initializeDatabase(logger log.Logger, config database.DatabaseConfig) (*sql.DB, func()) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// migrate database
	db, err := database.New(ctx, logger, config)
	if err != nil {
		panic(fmt.Sprintf("error creating database: %v", err))
	}

	shutdown := func() {
		fmt.Println("Shutting down the db")
		cancelFunc()
		if err := db.Close(); err != nil {
			logger.Log("exit", err)
		}
	}

	if err := database.RunMigrations(db, config); err != nil {
		panic(fmt.Sprintf("Error running migrations: %s", err))
	}

	return db, shutdown
}
