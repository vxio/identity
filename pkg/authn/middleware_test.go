package authn_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/moov-io/identity/pkg/authn"
)

func Test_Handler(t *testing.T) {
	a, s, f := Setup(t)

	ls := authn.LoginSession{}
	f.Fuzz(&ls)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	a.Nil(err)
	req.AddCookie(s.Cookie(ls))

	endpoint := newEndpoint(s, func(loginSession authn.LoginSession) {
		a.Equal(ls, loginSession)
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	a.Equal(200, recorder.Result().StatusCode)
}

func Test_NoAuthnCookie(t *testing.T) {
	a, s, _ := Setup(t)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	a.Nil(err)

	endpoint := newEndpoint(s, func(_ authn.LoginSession) {
		a.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	a.Equal(404, recorder.Result().StatusCode)
}

func Test_Expired(t *testing.T) {
	a, s, f := Setup(t)

	ls := authn.LoginSession{}
	f.Fuzz(&ls)
	ls.ExpiresAt = s.stime.Now().Add(time.Hour * -1).Unix()

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	a.Nil(err)
	req.AddCookie(s.Cookie(ls))

	endpoint := newEndpoint(s, func(loginSession authn.LoginSession) {
		a.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	a.Equal(404, recorder.Result().StatusCode)
}

func Test_NotBefore(t *testing.T) {
	a, s, f := Setup(t)

	ls := authn.LoginSession{}
	f.Fuzz(&ls)
	ls.NotBefore = s.stime.Now().Add(time.Hour).Unix()

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	a.Nil(err)
	req.AddCookie(s.Cookie(ls))

	endpoint := newEndpoint(s, func(loginSession authn.LoginSession) {
		a.Fail("Should not have ran")
	})

	recorder := httptest.NewRecorder()
	endpoint.ServeHTTP(recorder, req)

	a.Equal(404, recorder.Result().StatusCode)
}

func newEndpoint(s Scope, run func(loginSession authn.LoginSession)) http.Handler {
	mw, err := authn.NewMiddleware(s.stime, s.authnKeys)
	if err != nil {
		panic(err)
	}

	endpoint := mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authn.WithLoginSessionFromRequest(s.logger, w, r, func(loginSession authn.LoginSession) {
			run(loginSession)
			w.WriteHeader(200)
		})
	}))

	return endpoint
}
