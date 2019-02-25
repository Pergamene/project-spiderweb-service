package mysqlstore

import (
	"database/sql"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
)

// PageStore is the mysql for pages
type PageStore struct {
	db *sql.DB
}

// NewPageStore returns a PageStore
func NewPageStore(mysqldb *sql.DB) PageStore {
	return PageStore{
		db: mysqldb,
	}
}

// CreatePage creates a new page.
func (s PageStore) CreatePage(record page.Page, ownerID string) (page.Page, error) {
	// Prepare statement for inserting data
	statement, err := s.db.Prepare("INSERT INTO Page (`Version_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		return record, err
	}
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	defer statement.Close()
	result, err := statement.Exec(record.Version.ID, record.GUID, record.Title, record.Summary, record.PermissionType, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return record, err
	}
	id, err := result.LastInsertId()
	record.ID = id
	return record, nil
}
