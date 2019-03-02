package healthcheckrouter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
)

// HealthcheckHandler is the handler for page API
type HealthcheckHandler struct {
	HealthcheckService HealthcheckService
}

// HealthcheckService see Service for more details
type HealthcheckService interface {
	IsHealthy(ctx context.Context) (bool, error)
}

// Handlers returns the requests for the associated routes.
func Handlers(apiPath string, healthcheckService HealthcheckService) []api.RouterHandler {
	healthcheckHandler := HealthcheckHandler{
		HealthcheckService: healthcheckService,
	}
	var routerHandlers []api.RouterHandler
	routerHandlers = append(routerHandlers, api.RouterHandler{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("/%v/healthcheck", apiPath),
		Handle:   healthcheckHandler.IsHealthy,
	})
	return routerHandlers
}
