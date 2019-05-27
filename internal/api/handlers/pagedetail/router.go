package pagedetailhandler

import (
	"fmt"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
)

// HTTP path fragments keys
const (
	PageIDRouteKey       = "pageID"
	PageDetailIDRouteKey = "detailID"
)

// PageDetailRouterHandlers returns the requests for the associated routes.
func PageDetailRouterHandlers(apiPath string, pageDetailService PageDetailService) []api.RouterHandler {
	handler := PageDetailHandler{
		PageDetailService: pageDetailService,
	}
	var routerHandlers []api.RouterHandler
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPut,
		Endpoint: fmt.Sprintf("/%v/pages/:%v/details/:%v", apiPath, PageIDRouteKey, PageDetailIDRouteKey),
		Handle:   handler.UpdatePageDetail,
	})
	return routerHandlers
}
