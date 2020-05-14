package main

import (
	"flag"
	"fmt"

	"github.com/moov-io/identity/pkg/identity"
)

var (
	fCallerID = flag.String("caller", "00000000-0000-0000-0000-000000000000", "UUID of the caller")
	fTenantID = flag.String("tenant", "409189e3-b2f8-4646-93f8-3d622c3b8418", "UUID of the tenant")

	fInvite = flag.Bool("invite", false, "Flag to invite a user and exist")
	fEmail  = flag.String("email", "", "Email of the user to invite")
)

func main() {
	env, err := identity.NewEnvironment(nil)
	if err != nil {
		fmt.Println("Error loading up environment. ", err)
		return
	}
	defer env.Shutdown()

	flag.Parse()

	err = AddJJ(env)
	if err != nil {
		fmt.Println("Unable to add JJ\n  ", err)
		return
	}

	if *fInvite {
		env.Logger.Log("main", "Sending invite")
		sendInvite(*env)
	} else {
		env.Logger.Log("main", "Starting services")
		shutdown := env.RunServers()
		defer shutdown()
	}
}
