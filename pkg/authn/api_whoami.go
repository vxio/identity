package authn

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	api "github.com/moov-io/identity/pkg/api"
)

type whoAmIController struct{}

// NewWhoAmIController - Router for the Who Am I api routes.
func NewWhoAmIController() api.Router {
	return &whoAmIController{}
}

// Routes returns all of the api route for the InternalApiController
func (c *whoAmIController) Routes() api.Routes {
	return api.Routes{
		{
			Name:        "WhoAmI",
			Method:      strings.ToUpper("Get"),
			Pattern:     "/whoami",
			HandlerFunc: c.WhoAmI,
		},
	}
}

// WhoAmI - Responds back with information about the authenticated session
func (c *whoAmIController) WhoAmI(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}

	api.EncodeJSONResponse(user.(*jwt.Token).Claims.(jwt.MapClaims), nil, w)
}
