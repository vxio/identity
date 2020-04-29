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
	"github.com/moov-io/identity/pkg/database"
	identityserver "github.com/moov-io/identity/pkg/server"

	"github.com/moov-io/base/admin"
)

var logger log.Logger

func main() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	// Listen for application termination.
	terminationListener := newTerminationListener()

	//db setup
	_, close := initializeDatabase(logger)
	defer close()

	//internal admin server

	InternalApiService := identityserver.NewInternalApiService()
	InternalApiController := identityserver.NewInternalApiController(InternalApiService)

	adminRouter := identityserver.NewRouter(InternalApiController)

	adminServer := bootAdminServer(adminRouter, terminationListener, logger)
	defer adminServer.Shutdown()

	// public server

	IdentitiesApiService := identityserver.NewIdentitiesApiService()
	IdentitiesApiController := identityserver.NewIdentitiesApiController(IdentitiesApiService)

	CredentialsApiService := identityserver.NewCredentialsApiService()
	CredentialsApiController := identityserver.NewCredentialsApiController(CredentialsApiService)

	InvitesApiService := identityserver.NewInvitesApiService()
	InvitesApiController := identityserver.NewInvitesApiController(InvitesApiService)

	publicRouter := identityserver.NewRouter(IdentitiesApiController, CredentialsApiController, InvitesApiController, InternalApiController)

	_, shutdownServer := bootPublicServer(publicRouter, terminationListener, logger)
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

func bootAdminServer(adminRouter *mux.Router, errs chan<- error, logger log.Logger) *admin.Server {
	adminAddr := os.Getenv("HTTP_ADMIN_BIND_ADDRESS")
	if adminAddr == "" {
		adminAddr = ":8201"
	}

	adminServer := admin.NewServer(adminAddr)
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

func bootPublicServer(routes *mux.Router, errs chan<- error, logger log.Logger) (*http.Server, func()) {
	httpAddr := os.Getenv("HTTP_BIND_ADDRESS")
	if httpAddr == "" {
		httpAddr = ":8200"
	}

	// Create main HTTP server
	serve := &http.Server{
		Addr:    httpAddr,
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
		logger.Log("http", fmt.Sprintf("listening on %s", httpAddr))
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

func initializeDatabase(logger log.Logger) (*sql.DB, func()) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// migrate database
	db, err := database.New(ctx, logger, database.Type())
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
