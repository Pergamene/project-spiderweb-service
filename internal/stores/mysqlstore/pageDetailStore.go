package mysqlstore

import (
	"database/sql"
	"errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagedetail"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
)

// PageDetailStore is the mysql for a page detail
type PageDetailStore struct {
	db *sql.DB
}

// NewPageDetailStore returns a PageDetailStore
func NewPageDetailStore(mysqldb *sql.DB) PageDetailStore {
	return PageDetailStore{
		db: mysqldb,
	}
}

// UpdatePageDetail updates the given page.
func (s PageDetailStore) UpdatePageDetail(record pagedetail.PageDetail) error {
	if record.GUID == "" {
		return errors.New("must provide record.GUID to update the page")
	}
	if record.Title == "" {
		return errors.New("must provide record.Title to update the page")
	}
	if s.db == nil {
		return &storeerror.DBNotSetUp{}
	}
	statement, err := s.db.Prepare("UPDATE `Page` SET `title` = ?, `summary` = ? WHERE `guid` = ?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(record.Title, record.Summary, record.GUID)
	return err
}
