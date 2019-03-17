package healthcheckhandler

import (
	"fmt"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
)

// HealthcheckRouterHandlers returns the requests for the associated routes.
func HealthcheckRouterHandlers(apiPath string, healthcheckService HealthcheckService) []api.RouterHandler {
	handler := HealthcheckHandler{
		HealthcheckService: healthcheckService,
	}
	var routerHandlers []api.RouterHandler
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/%v/healthcheck", apiPath),
		Handle:   handler.IsHealthy,
	})
	return routerHandlers
}
