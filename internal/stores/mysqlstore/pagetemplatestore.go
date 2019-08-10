package mysqlstore

import (
	"database/sql"
	"errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
)

// PageTemplateStore is the mysql for pagetemplates
type PageTemplateStore struct {
	db *sql.DB
}

// NewPageTemplateStore returns a PageTemplateStore
func NewPageTemplateStore(mysqldb *sql.DB) PageTemplateStore {
	return PageTemplateStore{
		db: mysqldb,
	}
}

// GetPageTemplate returns the given pagetemplate.
func (s PageTemplateStore) GetPageTemplate(guid string) (pagetemplate.PageTemplate, error) {
	if guid == "" {
		return pagetemplate.PageTemplate{}, errors.New("must provide guid to get the pageTemplate")
	}
	if s.db == nil {
		return pagetemplate.PageTemplate{}, &storeerror.DBNotSetUp{}
	}
	clause := whereClause{
		operator: "AND", whereOperations: []whereOperation{
			{leftSide: "guid", operator: "= ?"},
			{leftSide: "guid", operator: "IS NULL"},
		},
	}
	statement := newSelectStatement([]string{"ID", "guid", "name"}, "PageTemplate", clause, 1)
	rows, err := s.db.Query(statement, guid)
	var pageTemplate pagetemplate.PageTemplate
	err = getSingleRow(guid, rows, err, &pageTemplate.ID, &pageTemplate.GUID, &pageTemplate.Name)
	return pageTemplate, err
}
