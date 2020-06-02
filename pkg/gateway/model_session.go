package gateway

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type Session struct {
	CallerID IdentityID
	TenantID TenantID
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

func NewRandomSession() Session {
	return Session{
		CallerID: IdentityID(uuid.New()),
		TenantID: TenantID(uuid.New()),
	}
}

func (s *Session) LogContext() map[string]string {
	return map[string]string{
		"identity_id": s.CallerID.String(),
		"tenant_id":   s.TenantID.String(),
	}
}
