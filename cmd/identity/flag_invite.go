package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/gateway"
	"github.com/moov-io/identity/pkg/identity"
)

func sendInvite(env identity.Environment) error {
	logCtx := env.Logger.WithMap(map[string]string{
		"email":  *fEmail,
		"tenant": *fTenantID,
		"caller": *fCallerID,
	})

	logCtx.Info().Log("Sending invite")

	tID, err := uuid.Parse(*fTenantID)
	if err != nil {
		return logCtx.Fatal().LogError("Unable to parse tenantID: "+*fTenantID, err)
	}

	cID, err := uuid.Parse(*fCallerID)
	if err != nil {
		return logCtx.Fatal().LogError("Unable to parse callerID: "+*fCallerID, err)
	}

	email := fEmail
	if fEmail == nil || !strings.Contains(*email, "@") {
		return logCtx.Fatal().LogError("Is not a valid email: "+*fEmail, errors.New("email is required and isn't valid"))
	}

	session := gateway.Session{
		CallerID: gateway.IdentityID(cID),
		TenantID: gateway.TenantID(tID),
	}

	invite := api.SendInvite{
		Email: *email,
	}

	sent, _, err := env.InviteService.SendInvite(session, invite)
	if err != nil {
		return err
	}

	prettyJSON, err := json.MarshalIndent(sent, "", "  ")
	if err != nil {
		return err
	}

	env.Logger.Info().Log(fmt.Sprintf("%s\n", string(prettyJSON)))
	return nil
}
