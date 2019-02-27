package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// Router handles the API routes
type Router struct {
	http.Handler
	NonAuthRoutes []NonAuthRoute
}

// NonAuthRoute a route that does not require authentication
type NonAuthRoute struct {
	Method  string
	Path    string
	Handler httprouter.Handle
}

// RouterHandlers are the handlers that are used for the routes.
type RouterHandlers struct {
	PageHandler        PageHandler
	HealthcheckHandler HealthcheckHandler
}

// PageHandler see handlers for more details.
type PageHandler interface {
	CreatePage(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

// HealthcheckHandler see handlers for more details.
type HealthcheckHandler interface {
	IsHealthy(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

// NewRouter adds the routes to a new handler and returns the handler with non-auth routes.
func NewRouter(apiPath, staticPath string, routerHandlers RouterHandlers) Router {
	handler := httprouter.New()
	handleAuthRoutes(handler, apiPath, routerHandlers)
	nonAuthRoutes := newNonAuthRoutes()
	handleNonAuthRoutes(handler, nonAuthRoutes)
	serveFiles(handler, apiPath, staticPath)
	handler.NotFound = http.HandlerFunc(handleNotFound)
	handler.MethodNotAllowed = http.HandlerFunc(handleMethodNotAllowed)
	handler.PanicHandler = panicHandler()
	return Router{
		Handler:       handler,
		NonAuthRoutes: nonAuthRoutes,
	}
}

func newNonAuthRoutes() []NonAuthRoute {
	return []NonAuthRoute{}
}

func handleAuthRoutes(handler *httprouter.Router, apiPath string, routerHandlers RouterHandlers) {
	handler.POST(fmt.Sprintf("/%v/page", apiPath), routerHandlers.PageHandler.CreatePage)
	handler.GET(fmt.Sprintf("/%v/healthcheck", apiPath), routerHandlers.HealthcheckHandler.IsHealthy)
}

func handleNonAuthRoutes(handler *httprouter.Router, nonAuthRoutes []NonAuthRoute) {
	for _, route := range nonAuthRoutes {
		handler.Handle(route.Method, route.Path, route.Handler)
	}
}

func serveFiles(handler *httprouter.Router, apiPath, staticPath string) {
	handler.ServeFiles(fmt.Sprintf("/%v/docs/*filepath", apiPath), http.Dir(fmt.Sprintf("%v/docs", staticPath)))
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	RespondWith(r, w, http.StatusNotFound, errors.New("not found"), nil)
}

func handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	RespondWith(r, w, http.StatusMethodNotAllowed, errors.New("method not allowed"), nil)
}

func panicHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	return func(w http.ResponseWriter, r *http.Request, e interface{}) {
		RespondWith(r, w, http.StatusInternalServerError, &InternalErr{}, errors.New("panicked"))
	}
}
