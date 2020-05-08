package authn

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	api "github.com/moov-io/identity/pkg/api"
)

type WhoAmiIController struct{}

func NewWhoAmIController() api.Router {
	return &WhoAmiIController{}
}

// Routes returns all of the api route for the InternalApiController
func (c *WhoAmiIController) Routes() api.Routes {
	return api.Routes{
		{
			"WhoAmI",
			strings.ToUpper("Get"),
			"/whoami",
			c.WhoAmI,
		},
	}
}

func (c *WhoAmiIController) WhoAmI(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}

	api.EncodeJSONResponse(user.(*jwt.Token).Claims.(jwt.MapClaims), nil, w)
}
