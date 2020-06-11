package authn

import (
	"net/http"
	"time"
)

func GetAuthnCookie(request *http.Request) (*http.Cookie, error) {
	cookie, err := request.Cookie("moov-authn")
	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func DeleteAuthnCookie(response http.ResponseWriter) {
	http.SetCookie(response, &http.Cookie{
		Name:     "moov-authn",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteDefaultMode,
		Secure:   false,
		HttpOnly: true,
	})
}
