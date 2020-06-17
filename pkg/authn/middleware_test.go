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
	a, s, f := Setup(t)

	ls := authn.LoginSession{}
	f.Fuzz(&ls)

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)
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

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)
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
	ls.Expiry = jwt.NewNumericDate(s.stime.Now().Add(time.Hour * -1))

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)
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
	ls.NotBefore = jwt.NewNumericDate(s.stime.Now().Add(time.Hour))

	req, err := http.NewRequest("GET", Host+"/", strings.NewReader(""))
	req.Header.Add("X-Forwarded-For", "1.2.3.4")
	req.Header.Add("Origin", Host)
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
	mw, err := authn.NewMiddleware(s.logger, s.stime, s.authnJwe)
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
