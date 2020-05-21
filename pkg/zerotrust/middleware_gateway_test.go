package zerotrust

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func Test_TenantID(t *testing.T) {
	uid := uuid.New()
	tid := TenantID(uid)
	assert.Equal(t, uid.String(), tid.String())
}

func Test_IdentityID(t *testing.T) {
	uid := uuid.New()
	tid := IdentityID(uid)
	assert.Equal(t, uid.String(), tid.String())
}

func Test_Handler(t *testing.T) {
	s := NewScope(t)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	s.a.Nil(err)

	tokenString := s.SignedString(s.NewSessionJwt())
	req.Header.Set("Authorization", "Bearer "+tokenString)

	endpoint := s.mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithSession(w, r, func(session Session) {
			s.a.Nil(err)
			s.a.Equal(s.identityID, session.CallerID)
			s.a.Equal(s.tenantID, session.TenantID)
		})
	}))

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)
}

func Test_NoAuthorizationHeader(t *testing.T) {
	s := NewScope(t)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	s.a.Nil(err)

	endpoint := s.mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithSession(w, r, func(session Session) {
			s.a.Fail("Should not have passed authentication")
		})
	}))

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)
}

func Test_NonBearerAuthorizationHeader(t *testing.T) {
	s := NewScope(t)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	s.a.Nil(err)

	req.SetBasicAuth("dne", "dne")

	endpoint := s.mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithSession(w, r, func(session Session) {
			s.a.Fail("Should not have passed authentication")
		})
	}))

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.a.Equal(404, recorder.Code)
}

func Test_ExpiredRequest(t *testing.T) {
	s := NewScope(t)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	s.a.Nil(err)

	session := s.NewSessionJwt()
	session.ExpiresAt = s.time.Now().Add(time.Hour * -1).Unix()

	tokenString := s.SignedString(session)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	endpoint := s.mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithSession(w, r, func(session Session) {
			s.a.Fail("Should not have passed authentication")
		})
	}))

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.a.Equal(404, recorder.Code)
}

func Test_NotBeforeRequest(t *testing.T) {
	s := NewScope(t)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	s.a.Nil(err)

	session := s.NewSessionJwt()
	session.NotBefore = s.time.Now().Add(time.Minute * 2).Unix()

	tokenString := s.SignedString(session)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	endpoint := s.mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WithSession(w, r, func(session Session) {
			s.a.Fail("Should not have passed authentication")
		})
	}))

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.a.Equal(404, recorder.Code)
}
