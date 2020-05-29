package webkeys

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/square/go-jose.v2"
)

type JWKSController struct {
	keyService WebKeysService
}

func NewJWKSController(keyService WebKeysService) *JWKSController {
	return &JWKSController{keyService: keyService}
}

func (c *JWKSController) AppendRoutes(router *mux.Router) *mux.Router {
	router.Name("well-known-jwks").Methods("GET").Path(c.WellKnownJwksPath()).HandlerFunc(c.WellKnownJwks)
	return router
}

func (c *JWKSController) WellKnownJwksPath() string {
	return "/.well-known/jwks.json"
}

func (c *JWKSController) WellKnownJwks(w http.ResponseWriter, r *http.Request) {
	keys, err := c.keyService.Keys()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// We only ever want to output keys that are public so heres a guard that checks that.
	publicKeys := jose.JSONWebKeySet{}
	for _, v := range keys.Keys {
		if v.IsPublic() {
			publicKeys.Keys = append(publicKeys.Keys, v)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(publicKeys); err != nil {
		return
	}
}
