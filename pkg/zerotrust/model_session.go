package zerotrust

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
	CallerID IdentityID `json:"iid"`
	TenantID TenantID   `json:"tid"`
}

type SessionJwt struct {
	jwt.StandardClaims

	Session
}

func SessionFromRequest(r *http.Request) (*Session, error) {
	session, ok := r.Context().Value(SessionContextKey).(*Session)
	if !ok || session == nil {
		return nil, errors.New("Unable to find Session in context")
	}
	return session, nil
}

func WithSession(w http.ResponseWriter, r *http.Request, run func(Session)) {
	session, err := SessionFromRequest(r)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	run(*session)
}
