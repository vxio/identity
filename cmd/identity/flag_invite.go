package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/moov-io/identity/pkg/api"
	"github.com/moov-io/identity/pkg/identity"
	"github.com/moov-io/identity/pkg/zerotrust"
)

func sendInvite(env identity.Environment) error {
	env.Logger.Log("level", "info", "msg", "Sending invite", "email", *fEmail, "tenant", *fTenantID, "caller", *fCallerID)

	tID, err := uuid.Parse(*fTenantID)
	if err != nil {
		env.Logger.Log("level", "fatal", "msg", "Unable to parse tenantID: "+*fTenantID)
		return err
	}

	cID, err := uuid.Parse(*fCallerID)
	if err != nil {
		env.Logger.Log("level", "fatal", "msg", "Unable to parse callerID: "+*fCallerID)
		return err
	}

	email := fEmail
	if fEmail == nil || !strings.Contains(*email, "@") {
		env.Logger.Log("level", "fatal", "msg", "Is not a valid email: "+*fEmail)
		return errors.New("email is required and isn't valid")
	}

	session := zerotrust.Session{
		CallerID: zerotrust.IdentityID(cID),
		TenantID: zerotrust.TenantID(tID),
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

	env.Logger.Log("level", "info", "msg", fmt.Sprintf("%s\n", string(prettyJSON)))
	return nil
}
