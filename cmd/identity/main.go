/*
 * Moov Identity API
 *
 * Handles all identities for tracking the users of the Moov platform.
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/moov-io/base/admin"
	config "github.com/moov-io/identity/pkg/config"
	"github.com/moov-io/identity/pkg/database"
	identityserver "github.com/moov-io/identity/pkg/server"
)

var logger log.Logger

func main() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	// Listen for application termination.
	terminationListener := newTerminationListener()

	ConfigService := config.NewConfigService(logger)

	config := &Config{}
	if err := ConfigService.Load(config); err != nil {
		return
	}

	fmt.Printf("%+v\n", config)

	//db setup
	db, close := initializeDatabase(logger, config.Database)
	defer close()

	TimeService := identityserver.NewSystemTimeService()
	NotificationsService := identityserver.NewNotificationsService("some.mail.server.com", 443, "username", "password", "noreply@moov.io")

	IdentityRepository := identityserver.NewIdentityRepository(db)
	IdentitiesService := identityserver.NewIdentitiesService(TimeService, IdentityRepository)

	CredentialRepository := identityserver.NewCredentialRepository(db)
	CredentialsService := identityserver.NewCredentialsService(TimeService, CredentialRepository)

	InvitesRepository := identityserver.NewInvitesRepository(db)
	InvitesService := identityserver.NewInvitesService(TimeService, InvitesRepository, NotificationsService)

	InternalService := identityserver.NewInternalService(*CredentialsService, *IdentitiesService)

	// internal admin server
	InternalController := identityserver.NewInternalAPIController(InternalService)
	adminRouter := identityserver.NewRouter(InternalController)
	adminServer := bootAdminServer(adminRouter, terminationListener, logger, config.Admin)
	defer adminServer.Shutdown()

	// public server

	// debug api
	WhoAmIController := identityserver.NewWhoAmIController()

	IdentitiesController := identityserver.NewIdentitiesApiController(IdentitiesService)
	CredentialsController := identityserver.NewCredentialsApiController(CredentialsService)
	InvitesController := identityserver.NewInvitesController(InvitesService)
	publicRouter := identityserver.NewRouter(IdentitiesController, CredentialsController, InvitesController, InternalController, WhoAmIController)
	_, shutdownServer := bootPublicServer(publicRouter, terminationListener, logger, config.HTTP)
	defer shutdownServer()

	awaitTermination(terminationListener)
}

func newTerminationListener() chan error {
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return errs
}

func awaitTermination(terminationListener chan error) {
	if err := <-terminationListener; err != nil {
		logger.Log("exit", err)
	}
}

func bootAdminServer(adminRouter *mux.Router, errs chan<- error, logger log.Logger, config HTTPConfig) *admin.Server {
	adminServer := admin.NewServer(config.Bind.Address)
	adminServer.AddHandler("/", adminRouter.ServeHTTP)

	go func() {
		logger.Log("admin", fmt.Sprintf("listening on %s", adminServer.BindAddr()))
		if err := adminServer.Listen(); err != nil {
			err = fmt.Errorf("problem starting admin http: %v", err)
			logger.Log("admin", err)
			errs <- err // send err to shutdown channel
		}
	}()

	return adminServer
}

func bootPublicServer(routes *mux.Router, errs chan<- error, logger log.Logger, config HTTPConfig) (*http.Server, func()) {

	// Create main HTTP server
	serve := &http.Server{
		Addr:    config.Bind.Address,
		Handler: routes,
		TLSConfig: &tls.Config{
			InsecureSkipVerify:       false,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
		},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start main HTTP server
	go func() {
		logger.Log("http", fmt.Sprintf("listening on %s", config.Bind.Address))
		if err := serve.ListenAndServe(); err != nil {
			err = fmt.Errorf("problem starting http: %v", err)
			logger.Log("http", err)
			errs <- err // send err to shutdown channel
		}
	}()

	shutdownServer := func() {
		if err := serve.Shutdown(context.TODO()); err != nil {
			logger.Log("exit", err)
		}
	}

	return serve, shutdownServer
}

func initializeDatabase(logger log.Logger, config database.DatabaseConfig) (*sql.DB, func()) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// migrate database
	db, err := database.New(ctx, logger, config)
	if err != nil {
		panic(fmt.Sprintf("error creating database: %v", err))
	}

	shutdown := func() {
		cancelFunc()
		if err := db.Close(); err != nil {
			logger.Log("exit", err)
		}
	}

	return db, shutdown
}
