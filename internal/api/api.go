package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// LocalEnv should be set on DATACENTER when running locally.
const LocalEnv = "LOCAL"

// Handler are the details for handling the API
type Handler struct {
	AuthN      Authenticator
	AuthZ      Authorizer
	Router     Router
	Datacenter string
	APIPath    string
}

// Authenticator inteface for authenticating.
type Authenticator interface {
	Authenticate(r *http.Request) (AuthData, error)
}

// Authorizer inteface for authorizing.
type Authorizer interface {
	Authorize(path string, authData AuthData) (AuthData, error)
}

// ServeHTTP handles responding to HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r = r.WithContext(ctx)
	if h.requiresNoAuth(w, r) {
		h.Router.ServeHTTP(w, r)
		return
	}
	r, authData, responded := h.authenticate(w, r)
	if responded {
		return
	}
	r = r.WithContext(ctx)
	r, authData, responded = h.authorize(w, r, authData)
	if responded {
		return
	}
	ctx = SetDataOnContext(ctx, authData)
	r = r.WithContext(ctx)
	h.Router.ServeHTTP(w, r)
}

func (h *Handler) requiresNoAuth(w http.ResponseWriter, r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, fmt.Sprintf("/%v/docs", h.APIPath)) {
		return true
	}
	if strings.HasPrefix(r.URL.Path, fmt.Sprintf("/%v/static/img", h.APIPath)) {
		return true
	}
	for _, nonAuthRoute := range h.Router.NonAuthRoutes {
		if r.URL.Path == nonAuthRoute.Path && r.Method == nonAuthRoute.Method {
			h.Router.ServeHTTP(w, r)
			return true
		}
	}
	return false
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) (*http.Request, AuthData, bool) {
	authData, err := h.AuthN.Authenticate(r)
	if castErr, ok := err.(*FailedAuthentication); ok {
		RespondWith(r, w, http.StatusUnauthorized, castErr, err)
		return r, authData, true
	}
	if err != nil {
		RespondWith(r, w, http.StatusInternalServerError, &InternalErr{}, errors.Wrap(err, "failed to determine authn"))
		return r, authData, true
	}
	return r, authData, false
}

func (h *Handler) authorize(w http.ResponseWriter, r *http.Request, authData AuthData) (*http.Request, AuthData, bool) {
	authData, err := h.AuthZ.Authorize(r.URL.Path, authData)
	if castErr, ok := err.(*FailedAuthorization); ok {
		RespondWith(r, w, http.StatusForbidden, castErr, err)
		return r, authData, true
	}
	if err != nil {
		RespondWith(r, w, http.StatusInternalServerError, &InternalErr{}, errors.Wrap(err, "failed to determine authz"))
		return r, authData, true
	}
	return r, authData, false
}
