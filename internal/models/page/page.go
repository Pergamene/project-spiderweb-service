package page

import (
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagedetail"
	"github.com/Pergamene/project-spiderweb-service/internal/models/pageproperty"
	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"

	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
)

// Page is the main object that houses most information.
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

// GetJSONConformed conforms the page to be ready for JSON marshelling.
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
