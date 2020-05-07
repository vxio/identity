package identityserver

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type TenantID uuid.UUID

func (id TenantID) String() string {
	return uuid.UUID(id).String()
}

type IdentityID uuid.UUID

func (id IdentityID) String() string {
	return uuid.UUID(id).String()
}

type Session struct {
	CallerID IdentityID
	TenantID TenantID
}

func NewSessionFromRequest(r *http.Request) (*Session, error) {
	token, ok := r.Context().Value("user").(*jwt.Token)
	if !ok {
		return nil, errors.New("Unable to cast context `user` into jwt.Token")
	}

	m, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Unable to cast context `Claims` into `jwt.MapClaims`")
	}

	sub, ok := m[`sub`].(string)
	if !ok {
		return nil, errors.New("Unable to cast context `m[sub]` into `string``")
	}

	caller, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}

	tenant, err := uuid.Parse(r.Header.Get("x-tenant-id"))
	if err != nil {
		return nil, err
	}

	session := Session{
		CallerID: IdentityID(caller),
		TenantID: TenantID(tenant),
	}

	return &session, nil
}

func WithSession(w http.ResponseWriter, r *http.Request, run func(Session)) {
	session, err := NewSessionFromRequest(r)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	run(*session)
	return
}
