package session

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moov-io/identity/pkg/client"
	"github.com/moov-io/identity/pkg/logging"
	tmw "github.com/moov-io/tumbler/pkg/middleware"
)

type SessionController interface {
	AppendRoutes(router *mux.Router) *mux.Router
}

func NewSessionController(logger logging.Logger, service SessionService) SessionController {
	return &sessionController{
		logger:  logger,
		service: service,
	}
}

type sessionController struct {
	logger  logging.Logger
	service SessionService
}

func (c sessionController) AppendRoutes(router *mux.Router) *mux.Router {

	router.
		Name("Identity.getSession").
		Methods("GET").
		Path("/session").
		HandlerFunc(c.getSessionHandler)

	router.
		Name("Identity.changeSession").
		Methods("PUT").
		Path("/session").
		HandlerFunc(c.changeTenantHandler)

	return router
}

func (c *sessionController) jsonResponse(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(value)
}

func (c *sessionController) errorResponse(w http.ResponseWriter, err error) {
	switch err {
	case ErrIdentityNotFound:
		w.WriteHeader(404)
	case sql.ErrNoRows:
		w.WriteHeader(404)
	default:
		w.WriteHeader(500)
	}
}

func (c *sessionController) getSessionHandler(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		details, err := c.service.GetDetails(claims)
		if err != nil {
			c.errorResponse(w, err)
			return
		}

		c.jsonResponse(w, &details)
	})
}

func (c *sessionController) changeTenantHandler(w http.ResponseWriter, r *http.Request) {
	tmw.WithClaimsFromRequest(w, r, func(claims tmw.TumblerClaims) {
		update := &client.ChangeSessionDetails{}
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			w.WriteHeader(400)
			return
		}

		details, cookie, err := c.service.ChangeDetails(r, claims, *update)
		if err != nil {
			c.errorResponse(w, err)
			return
		}

		http.SetCookie(w, cookie)
		c.jsonResponse(w, &details)
	})
}
