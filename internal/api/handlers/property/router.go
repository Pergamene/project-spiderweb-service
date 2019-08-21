package propertyhandler

import (
	"fmt"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
)

// HTTP path fragments keys
const (
	PropertyIDRouteKey = "propertyID"
)

// PropertyRouterHandlers returns the requests for the associated routes.
func PropertyRouterHandlers(apiPath string, propertyService PropertyService) []api.RouterHandler {
	handler := PropertyHandler{
		PropertyService: propertyService,
	}
	var routerHandlers []api.RouterHandler
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPost,
		Endpoint: fmt.Sprintf("/%v/properties", apiPath),
		Handle:   handler.CreateProperty,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPatch,
		Endpoint: fmt.Sprintf("/%v/properties/:%v", apiPath, PropertyIDRouteKey),
		Handle:   handler.UpdateProperty,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodDelete,
		Endpoint: fmt.Sprintf("/%v/properties/:%v", apiPath, PropertyIDRouteKey),
		Handle:   handler.DisableProperty,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodPost,
		Endpoint: fmt.Sprintf("/%v/properties/:%v", apiPath, PropertyIDRouteKey),
		Handle:   handler.EnableProperty,
	})
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/%v/properties", apiPath),
		Handle:   handler.GetProperties,
	})
	return routerHandlers
}
