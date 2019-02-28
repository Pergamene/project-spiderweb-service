package store

import "github.com/Pergamene/project-spiderweb-service/internal/models/page"

// PageStore defines the required functionality for any associated store.
type PageStore interface {
	AssertCanModifyPage(pageGUID, userID string) (bool, error)
	UpdatePage(record page.Page) error
	CreatePage(record page.Page, ownerID string) (page.Page, error)
}
