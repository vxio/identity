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
	tID, err := uuid.Parse(*fTenantID)
	if err != nil {
		return err
	}

	cID, err := uuid.Parse(*fCallerID)
	if err != nil {
		return err
	}

	email := fEmail
	if fEmail == nil || strings.Contains(*email, "@") {
		return errors.New("email is required and isn't valid")
	}

	session := zerotrust.Session{
		CallerID: zerotrust.IdentityID(cID),
		TenantID: zerotrust.TenantID(tID),
	}

	invite := api.SendInvite{
		Email: *email,
	}

	sent, err := env.InviteService.SendInvite(session, invite)
	if err != nil {
		return err
	}

	prettyJSON, err := json.MarshalIndent(sent, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(prettyJSON))
	return nil
}
