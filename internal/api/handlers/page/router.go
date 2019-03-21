package pagehandler

import (
	"fmt"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
)

// HTTP path fragments keys
const (
	PageIDRouteKey = "pageID"
)

// PageRouterHandlers returns the requests for the associated routes.
func PageRouterHandlers(apiPath string, pageService PageService) []api.RouterHandler {
	handler := PageHandler{
		PageService: pageService,
	}
	var routerHandlers []api.RouterHandler
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPost,
		Endpoint: fmt.Sprintf("/%v/page", apiPath),
		Handle:   handler.CreatePage,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPut,
		Endpoint: fmt.Sprintf("/%v/page/:%v", apiPath, PageIDRouteKey),
		Handle:   handler.UpdatePage,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/%v/page/:%v", apiPath, PageIDRouteKey),
		Handle:   handler.GetPage,
	})
	return routerHandlers
}
