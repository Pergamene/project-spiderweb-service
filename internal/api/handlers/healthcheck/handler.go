package healthcheckhandler

import (
	"context"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// HealthcheckService see Service for more details
type HealthcheckService interface {
	IsHealthy(ctx context.Context) (bool, error)
}

// HealthcheckHandler is the handler for page API
type HealthcheckHandler struct {
	HealthcheckService HealthcheckService
}

// IsHealthy see Service for more details
func (h HealthcheckHandler) IsHealthy(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	if !authData.IsAdmin() {
		api.RespondWith(r, w, http.StatusUnauthorized, &api.InternalErr{}, errors.Wrapf(err, "user not authorized for healthcheck: %v", authData.UserID))
		return
	}
	isHealthy, err := h.HealthcheckService.IsHealthy(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	statusString := "ok"
	if !isHealthy {
		statusString = "error"
	}
	api.RespondWith(r, w, http.StatusOK, map[string]string{"status": statusString}, nil)
}
