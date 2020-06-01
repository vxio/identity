package main

import (
	"flag"
	"os"

	"github.com/moov-io/identity/pkg/identity"
	"github.com/moov-io/identity/pkg/logging"
)

var (
	fCallerID = flag.String("caller", "00000000-0000-0000-0000-000000000000", "UUID of the caller")
	fTenantID = flag.String("tenant", "409189e3-b2f8-4646-93f8-3d622c3b8418", "UUID of the tenant")

	fInvite = flag.Bool("invite", false, "Flag to invite a user and exist")
	fEmail  = flag.String("email", "", "Email of the user to invite")
)

func main() {
	logger := logging.NewDefaultLogger().WithKeyValue("app", "identity")

	env, err := identity.NewEnvironment(logger, nil)
	if err != nil {
		logger.Fatal().LogError("Error loading up environment.", err)
		os.Exit(1)
	}
	defer env.Shutdown()

	env.Logger.Info().Log("Environment built")

	flag.Parse()

	if *fInvite {
		env.Logger.Info().Log("Sending invite")
		if err := sendInvite(*env); err != nil {
			env.Logger.Fatal().LogError("Unable to send invite", err)
			os.Exit(1)
		}
	} else {
		env.Logger.Info().Log("Starting services")
		shutdown := env.RunServers(true)
		defer shutdown()
	}

	os.Exit(0)
}
