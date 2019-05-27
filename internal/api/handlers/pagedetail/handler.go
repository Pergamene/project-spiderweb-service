package pagedetailhandler

import (
	"context"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagedetail"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	pagedetailservice "github.com/Pergamene/project-spiderweb-service/internal/services/pagedetail"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// PageDetailService see Service for more details
type PageDetailService interface {
	UpdatePageDetail(ctx context.Context, params pagedetailservice.UpdatePageDetailParams) error
}

// PageDetailHandler is the handler for the associated API
type PageDetailHandler struct {
	PageDetailService PageDetailService
}

// UpdatePageDetail see Service for more details
func (h PageDetailHandler) UpdatePageDetail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewUpdatePageDetailRequest(r, p)
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
	err = h.PageDetailService.UpdatePageDetail(ctx, pagedetailservice.UpdatePageDetailParams{
		Detail: pagedetail.PageDetail{
			Title:   request.Title,
			Summary: request.Summary,
			//@TODO:
			// Partitions: request.Partitions,
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
