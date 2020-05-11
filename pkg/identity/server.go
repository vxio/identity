package identity

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/moov-io/base/admin"
	_ "github.com/moov-io/identity" // need to import the embedded files

	"github.com/go-kit/kit/log"
)

func (env *Environment) RunServers() {

	// Listen for application termination.
	terminationListener := newTerminationListener()

	adminServer := bootAdminServer(terminationListener, env.Logger, env.Config.Servers.Admin)
	defer adminServer.Shutdown()

	_, shutdownPrivateServer := bootHTTPServer("private", &env.PrivateRouter, terminationListener, env.Logger, env.Config.Servers.Private)
	defer shutdownPrivateServer()

	_, shutdownPublicServer := bootHTTPServer("public", &env.PublicRouter, terminationListener, env.Logger, env.Config.Servers.Public)
	defer shutdownPublicServer()

	awaitTermination(env.Logger, terminationListener)
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

func awaitTermination(logger log.Logger, terminationListener chan error) {
	if err := <-terminationListener; err != nil {
		logger.Log("exit", err)
	}
}

func bootHTTPServer(name string, routes *mux.Router, errs chan<- error, logger log.Logger, config HTTPConfig) (*http.Server, func()) {

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
		logger.Log(name, fmt.Sprintf("listening on %s", config.Bind.Address))
		if err := serve.ListenAndServe(); err != nil {
			err = fmt.Errorf("problem starting http: %v", err)
			logger.Log(name, err)
			errs <- err // send err to shutdown channel
		}
	}()

	shutdownServer := func() {
		if err := serve.Shutdown(context.TODO()); err != nil {
			logger.Log(name, err)
		}
	}

	return serve, shutdownServer
}

func bootAdminServer(errs chan<- error, logger log.Logger, config HTTPConfig) *admin.Server {
	adminServer := admin.NewServer(config.Bind.Address)

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
