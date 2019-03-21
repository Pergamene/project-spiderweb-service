package pageservice

import (
	"context"

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
	params.Page.GUID = "PG_123456789012"
	page, err := s.PageStore.CreatePage(params.Page, params.OwnerID)
	if err != nil {
		return page, errors.Wrapf(err, "failed to create page: %+v", params)
	}
	return page, nil
}

// UpdatePageParams params for UpdatePage
type UpdatePageParams struct {
	Page   page.Page
	UserID string
}

// UpdatePage Updates a new page.
func (s PageService) UpdatePage(ctx context.Context, params UpdatePageParams) error {
	_, err := s.PageStore.CanEditPage(params.Page.GUID, params.UserID)
	if err != nil {
		return err
	}
	err = s.PageStore.UpdatePage(params.Page)
	if err != nil {
		return errors.Wrapf(err, "failed to update page: %+v", params)
	}
	return nil
}

// GetPageParams params for GetPage
type GetPageParams struct {
	Page   page.Page
	UserID string
}

// GetPage Updates a new page.
func (s PageService) GetPage(ctx context.Context, params GetPageParams) (page.Page, error) {
	_, err := s.PageStore.CanReadPage(params.Page.GUID, params.UserID)
	if err != nil {
		return page.Page{}, err
	}
	p, err := s.PageStore.GetPage(params.Page.GUID)
	if err != nil {
		return p, errors.Wrapf(err, "failed to get page: %+v", params)
	}
	return p, nil
}
