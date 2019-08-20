package store

import (
	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"
)

// PageTemplateStore defines the required functionality for any associated store.
type PageTemplateStore interface {
	GetPageTemplate(pageTemplateGUID string) (pagetemplate.PageTemplate, error)
}
