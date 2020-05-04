package identityserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type WhoAmiIController struct{}

func NewWhoAmIController() Router {
	return &WhoAmiIController{}
}

// Routes returns all of the api route for the InternalApiController
func (c *WhoAmiIController) Routes() Routes {
	return Routes{
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

	EncodeJSONResponse(user.(*jwt.Token).Claims.(jwt.MapClaims), nil, w)
}
