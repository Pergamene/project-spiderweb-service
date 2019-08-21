package pageservice

import (
	"context"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/property"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/store"
	"github.com/pkg/errors"
)

// PageService is the service for handling page-related APIs
type PageService struct {
	PageStore         store.PageStore
	PageTemplateStore store.PageTemplateStore
	VersionStore      store.VersionStore
	UserStore         store.UserStore
}

// CreatePageParams params for CreatePage
type CreatePageParams struct {
	Page    page.Page
	OwnerID string
}

// CreatePage creates a new page.
func (s PageService) CreatePage(ctx context.Context, params CreatePageParams) (page.Page, error) {
	err := s.populatePageIDs(ctx, &params.Page)
	if err != nil {
		return page.Page{}, err
	}
	pageGUID, err := s.PageStore.GetUniquePageGUID(params.Page.GUID)
	if err != nil {
		return page.Page{}, err
	}
	params.Page.GUID = pageGUID
	u, err := s.UserStore.GetUser(params.OwnerID)
	page, err := s.PageStore.CreatePage(params.Page, u.ID)
	if err != nil {
		return page, errors.Wrapf(err, "failed to create page: %+v", params)
	}
	return page, nil
}

func (s PageService) populatePageIDs(ctx context.Context, p *page.Page) error {
	if p.PageTemplate.GUID != "" {
		pt, err := s.PageTemplateStore.GetPageTemplate(p.PageTemplate.GUID)
		if err != nil {
			return err
		}
		p.PageTemplate = pt
	}
	if p.Version.GUID != "" {
		v, err := s.VersionStore.GetVersion(p.Version.GUID)
		if err != nil {
			return err
		}
		p.Version = v
	}
	return nil
}

// UpdatePageParams params for UpdatePage
type UpdatePageParams struct {
	Page   page.Page
	UserID string
}

// UpdatePage sets a page to what is provided.
func (s PageService) UpdatePage(ctx context.Context, params UpdatePageParams) error {
	_, err := s.PageStore.CanEditPage(params.Page.GUID, params.UserID)
	if err != nil {
		return err
	}
	err = s.populatePageIDs(ctx, &params.Page)
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

// GetPage returns just the page entity.
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

// GetEntirePageParams params for GetEntirePage
type GetEntirePageParams struct {
	Page   page.Page
	UserID string
}

// GetEntirePage returns a full page object, with properties, details, etc.
func (s PageService) GetEntirePage(ctx context.Context, params GetEntirePageParams) (page.Page, error) {
	_, err := s.PageStore.CanReadPage(params.Page.GUID, params.UserID)
	if err != nil {
		return page.Page{}, err
	}
	p, err := s.GetPage(ctx, GetPageParams{
		Page:   params.Page,
		UserID: params.UserID,
	})
	if err != nil {
		return p, err
	}
	err = s.populatePageIDs(ctx, &p)
	if err != nil {
		return p, errors.Wrapf(err, "failed to populate page with ids: %+v", params)
	}
	return p, nil
}

// GetPagesParams params for GetPages
type GetPagesParams struct {
	NextBatchID string
	UserID      string
}

// GetPages returns a list of pages filtered and ordered as specified.
func (s PageService) GetPages(ctx context.Context, params GetPagesParams) ([]page.Page, int, string, error) {
	ps, total, nextBatchID, err := s.PageStore.GetPages(params.UserID, params.NextBatchID, 10)
	if err != nil {
		return ps, total, nextBatchID, errors.Wrapf(err, "failed to get pages: %+v", params)
	}
	return ps, total, nextBatchID, nil
}

// RemovePageParams params for RemovePage
type RemovePageParams struct {
	Page   page.Page
	UserID string
}

// RemovePage marks the page as removed.
func (s PageService) RemovePage(ctx context.Context, params RemovePageParams) error {
	_, err := s.PageStore.CanEditPage(params.Page.GUID, params.UserID)
	if err != nil {
		return err
	}
	err = s.PageStore.RemovePage(params.Page.GUID)
	if err != nil {
		return errors.Wrapf(err, "failed to remove page: %+v", params)
	}
	return nil
}

// GetPagePropertiesParams params for GetPageProperties
type GetPagePropertiesParams struct {
	Page   page.Page
	UserID string
}

// GetPageProperties returns the page's properties.
func (s PageService) GetPageProperties(ctx context.Context, params GetPagePropertiesParams) ([]property.Property, error) {
	ps := make([]property.Property, 0)
	_, err := s.PageStore.CanReadPage(params.Page.GUID, params.UserID)
	if err != nil {
		return ps, err
	}
	ps, err = s.PageStore.GetPageProperties(params.Page.GUID)
	if err != nil {
		return ps, errors.Wrapf(err, "failed to get page properties: %+v", params)
	}
	return ps, nil
}

// ReplacePagePropertiesParams params for ReplacePageProperties
type ReplacePagePropertiesParams struct {
	Page       page.Page
	Properties []property.Property
	UserID     string
}

// ReplacePageProperties replaces the current page's properties with the new properties.
func (s PageService) ReplacePageProperties(ctx context.Context, params ReplacePagePropertiesParams) error {
	_, err := s.PageStore.CanEditPage(params.Page.GUID, params.UserID)
	if err != nil {
		return err
	}
	err = s.PageStore.ReplacePageProperties(params.Page.GUID, params.Properties)
	if err != nil {
		return errors.Wrapf(err, "failed to replace page properties: %+v", params)
	}
	return nil
}
