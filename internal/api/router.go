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

// RouterHandler is the information needed to establish a handle for the route.
type RouterHandler struct {
	Method   string
	Endpoint string
	Handle   httprouter.Handle
}

// NewRouter adds the routes to a new handler and returns the handler with non-auth routes.
func NewRouter(apiPath, staticPath string, routerHandlers []RouterHandler) Router {
	handler := httprouter.New()
	handleAuthRoutes(handler, routerHandlers)
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

func handleAuthRoutes(handler *httprouter.Router, routerHandlers []RouterHandler) {
	for _, routerHandler := range routerHandlers {
		handleAuthRoute(handler, routerHandler)
	}
}

func handleAuthRoute(handler *httprouter.Router, routerHandler RouterHandler) {
	if routerHandler.Method == http.MethodGet {
		handler.GET(routerHandler.Endpoint, routerHandler.Handle)
	} else if routerHandler.Method == http.MethodPut {
		handler.PUT(routerHandler.Endpoint, routerHandler.Handle)
	} else if routerHandler.Method == http.MethodPost {
		handler.POST(routerHandler.Endpoint, routerHandler.Handle)
	} else if routerHandler.Method == http.MethodDelete {
		handler.DELETE(routerHandler.Endpoint, routerHandler.Handle)
	} else if routerHandler.Method == http.MethodPatch {
		handler.PATCH(routerHandler.Endpoint, routerHandler.Handle)
	} else if routerHandler.Method == http.MethodOptions {
		handler.OPTIONS(routerHandler.Endpoint, routerHandler.Handle)
	}
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
		RespondWith(r, w, http.StatusInternalServerError, &InternalErr{}, errors.Errorf("panicked\n%+v", e))
	}
}
