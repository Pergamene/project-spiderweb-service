package store

import "github.com/Pergamene/project-spiderweb-service/internal/models/page"

// PageStore defines the required functionality for any associated store.
type PageStore interface {
	GetUniquePageGUID(proposedPageGUID string) (string, error)
	CanEditPage(pageGUID, userID string) (bool, error)
	CanReadPage(pageGUID, userID string) (bool, error)
	SetPage(record page.Page) error
	CreatePage(record page.Page, ownerID string) (page.Page, error)
	GetPage(pageGUID string) (page.Page, error)
	GetPages(userID string, nextBatchID string, limit int) ([]page.Page, int, string, error)
	RemovePage(pageGUID string) error
}
