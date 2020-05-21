package zerotrust

import (
	"errors"
	"fmt"
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
		fmt.Printf("%+v\n", r.Context())
		return nil, errors.New("Unable to find Session in context")
	}
	return session, nil
}

func WithSession(w http.ResponseWriter, r *http.Request, run func(Session)) {
	session, err := SessionFromRequest(r)
	if err != nil {
		fmt.Println("Session not found", err)
		w.WriteHeader(500)
		return
	}

	run(*session)
	return
}
