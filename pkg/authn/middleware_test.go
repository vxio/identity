package authn_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/moov-io/identity/pkg/authn"
	"github.com/square/go-jose/jwt"
)

const Host = "http://local.moov.io"

func Test_Handler(t *testing.T) {
	s := Setup(t)

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	s.assert.Nil(err)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)

	ls := s.AddSession(req, func(ls *authn.LoginSession) {})

	endpoint := newEndpoint(s, nil, func(loginSession authn.LoginSession) {
		s.assert.Equal(ls, loginSession)
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.assert.Equal(200, recorder.Result().StatusCode)
}

func Test_NoAuthnCookie(t *testing.T) {
	s := Setup(t)

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	s.assert.Nil(err)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)

	endpoint := newEndpoint(s, nil, func(_ authn.LoginSession) {
		s.assert.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.assert.Equal(404, recorder.Result().StatusCode)
}

func Test_Expired(t *testing.T) {
	s := Setup(t)

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	s.assert.Nil(err)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)

	s.AddSession(req, func(ls *authn.LoginSession) {
		ls.Expiry = jwt.NewNumericDate(s.stime.Now().Add(time.Hour * -1))
	})

	endpoint := newEndpoint(s, nil, func(loginSession authn.LoginSession) {
		s.assert.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.assert.Equal(404, recorder.Result().StatusCode)
}

func Test_NotBefore(t *testing.T) {
	s := Setup(t)

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	s.assert.Nil(err)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)

	s.AddSession(req, func(ls *authn.LoginSession) {
		ls.NotBefore = jwt.NewNumericDate(s.stime.Now().Add(time.Hour))
	})

	endpoint := newEndpoint(s, nil, func(loginSession authn.LoginSession) {
		s.assert.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.assert.Equal(404, recorder.Result().StatusCode)
}

func Test_Scope(t *testing.T) {
	s := Setup(t)

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	s.assert.Nil(err)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)

	s.AddSession(req, func(ls *authn.LoginSession) {
		ls.NotBefore = jwt.NewNumericDate(s.stime.Now().Add(time.Hour))
		ls.Scopes = []string{"must_have1", "must_have2", "other_scope"}
	})

	scopes := []string{"must_have1", "must_have2"}
	endpoint := newEndpoint(s, &scopes, func(loginSession authn.LoginSession) {
		s.assert.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.assert.Equal(404, recorder.Result().StatusCode)
}

func Test_Scope_Missing(t *testing.T) {
	s := Setup(t)

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	s.assert.Nil(err)
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)

	s.AddSession(req, func(ls *authn.LoginSession) {
		ls.NotBefore = jwt.NewNumericDate(s.stime.Now().Add(time.Hour))
	})

	scopes := []string{"must_have"}
	endpoint := newEndpoint(s, &scopes, func(loginSession authn.LoginSession) {
		s.assert.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	s.assert.Equal(404, recorder.Result().StatusCode)
}

func newEndpoint(s Scope, scopes *[]string, run func(loginSession authn.LoginSession)) http.Handler {
	if scopes == nil {
		scopes = new([]string)
	}

	mw, err := authn.NewMiddleware(s.logger, s.stime, s.authnJwe)
	if err != nil {
		panic(err)
	}

	endpoint := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authn.WithLoginSessionFromRequest(s.logger, w, r, *scopes, func(loginSession authn.LoginSession) {
			run(loginSession)
			w.WriteHeader(200)
		})
	}))

	return endpoint
}
