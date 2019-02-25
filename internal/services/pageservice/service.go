package pageservice

import (
	"context"

	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/store"
	"github.com/pkg/errors"
)

// PageService is the service for handling page-related APIs
type PageService struct {
	PageStore store.PageStore
}

// CreatePageParams params for CreatePage
type CreatePageParams struct {
	Page    page.Page
	OwnerID string
}

// CreatePage creates a new page.
func (s PageService) CreatePage(ctx context.Context, params CreatePageParams) (page.Page, error) {
	// @TODO:
	params.Page.Version.ID = 1
	params.Page.GUID = "PG_123456789012"
	params.Page.PermissionType = permission.TypePrivate
	page, err := s.PageStore.CreatePage(params.Page, params.OwnerID)
	if err != nil {
		return page, errors.Wrapf(err, "failed to create page: %+v", params)
	}
	return page, nil
}
