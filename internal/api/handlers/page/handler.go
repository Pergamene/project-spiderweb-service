package pagehandler

import (
	"context"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
	pageservice "github.com/Pergamene/project-spiderweb-service/internal/services/page"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// PageService see Service for more details
type PageService interface {
	CreatePage(ctx context.Context, params pageservice.CreatePageParams) (page.Page, error)
	UpdatePage(ctx context.Context, params pageservice.UpdatePageParams) error
}

// PageHandler is the handler for page API
type PageHandler struct {
	PageService PageService
}

// CreatePage see Service for more details
func (h PageHandler) CreatePage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewCreatePageRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	record, err := h.PageService.CreatePage(ctx, pageservice.CreatePageParams{
		Page: page.Page{
			Title:   request.Title,
			Summary: request.Summary,
			Version: version.Version{
				ID: request.VersionID,
			},
			PermissionType: request.PermissionType,
		},
		OwnerID: authData.UserID,
	})
	if castErr, ok := err.(*storeerror.DupEntry); ok {
		api.RespondWith(r, w, http.StatusBadRequest, castErr, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	api.RespondWith(r, w, http.StatusOK, map[string]string{"id": record.GUID}, nil)
}

// UpdatePage see Service for more details
func (h PageHandler) UpdatePage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewUpdatePageRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	err = h.PageService.UpdatePage(ctx, pageservice.UpdatePageParams{
		Page: page.Page{
			GUID:    request.GUID,
			Title:   request.Title,
			Summary: request.Summary,
		},
		UserID: authData.UserID,
	})
	if _, ok := err.(*storeerror.NotAuthorized); ok {
		api.RespondWith(r, w, http.StatusUnauthorized, &api.FailedAuthorization{}, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	api.RespondWith(r, w, http.StatusOK, nil, nil)
}
