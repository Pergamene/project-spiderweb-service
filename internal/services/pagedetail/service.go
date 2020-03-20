package pagedetailservice

import (
	"context"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagedetail"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/store"
	"github.com/pkg/errors"
)

// PageDetailService is the service for handling page detail-related APIs
type PageDetailService struct {
	PageDetailStore store.PageDetailStore
}

// UpdatePageDetailParams params for UpdatePageDetail
type UpdatePageDetailParams struct {
	Detail pagedetail.PageDetail
	PageID string
	UserID string
}

// UpdatePageDetail Updates a page detail.
func (s PageDetailService) UpdatePageDetail(ctx context.Context, params UpdatePageDetailParams) error {
	err := s.PageDetailStore.UpdatePageDetail(params.Detail)
	if err != nil {
		return errors.Wrapf(err, "failed to update detail: %v", params)
	}
	return nil
}
