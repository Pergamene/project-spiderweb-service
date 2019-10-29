package page

import (
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagedetail"
	"github.com/Pergamene/project-spiderweb-service/internal/models/pageproperty"
	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"

	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
)

// ReducedPage is a page object as it is realized from the GetPage API rather than the GetEntirePage API (see ExpandedPage).
type ReducedPage struct {
	ID             int64           `json:"-"`
	VersionID      string          `json:"versionId"`
	PageTemplateID string          `json:"pageTemplateId"`
	GUID           string          `json:"id"`
	Title          string          `json:"title"`
	Summary        string          `json:"summary"`
	PermissionType permission.Type `json:"permission"`
	CreatedAt      *time.Time      `json:"createdAt"`
	UpdatedAt      *time.Time      `json:"updatedAt"`
	DeletedAt      *time.Time      `json:"deletedAt,omitempty"`
}

// Page is the entire page object that aggregates all its information.
type Page struct {
	ID             int64                       `json:"-"`
	Version        version.Version             `json:"version"`
	PageTemplate   pagetemplate.PageTemplate   `json:"pageTemplate"`
	GUID           string                      `json:"id"`
	Title          string                      `json:"title"`
	Summary        string                      `json:"summary"`
	PermissionType permission.Type             `json:"permission"`
	PageProperties []pageproperty.PageProperty `json:"properties"`
	PageDetails    []pagedetail.PageDetail     `json:"details"`
	CreatedAt      *time.Time                  `json:"createdAt"`
	UpdatedAt      *time.Time                  `json:"updatedAt"`
	DeletedAt      *time.Time                  `json:"deletedAt,omitempty"`
}

// Expand returns an Page version of the reference ReducedPage.
func (p ReducedPage) Expand() Page {
	return Page{
		ID: p.ID,
		Version: version.Version{
			GUID: p.VersionID,
		},
		PageTemplate: pagetemplate.PageTemplate{
			GUID: p.PageTemplateID,
		},
		GUID:           p.GUID,
		Title:          p.Title,
		Summary:        p.Summary,
		PermissionType: p.PermissionType,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		DeletedAt:      p.DeletedAt,
	}
}

// Reduce returns an ReducedPage version of the reference Page.
func (p Page) Reduce() ReducedPage {
	return ReducedPage{
		ID:             p.ID,
		VersionID:      p.Version.GUID,
		PageTemplateID: p.PageTemplate.GUID,
		GUID:           p.GUID,
		Title:          p.Title,
		Summary:        p.Summary,
		PermissionType: p.PermissionType,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		DeletedAt:      p.DeletedAt,
	}
}

// GetJSONConformed conforms the expanded page to be ready for JSON marshelling.
func (p Page) GetJSONConformed() interface{} {
	// see: https://stackoverflow.com/questions/33183071/golang-serialize-deserialize-an-empty-array-not-as-null
	if p.PageProperties == nil {
		p.PageProperties = []pageproperty.PageProperty{}
	}
	if p.PageDetails == nil {
		p.PageDetails = []pagedetail.PageDetail{}
	}
	return p
}

// GetJSONConformed conforms the page to be ready for JSON marshelling.
func (p ReducedPage) GetJSONConformed() interface{} {
	return p
}
