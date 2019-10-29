package mysqlstore

import (
	"database/sql"
	"errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"
	"github.com/Pergamene/project-spiderweb-service/internal/util/wrapsql"

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
	statement := wrapsql.SelectStatement{
		Selectors: []string{"ID", "guid", "name"},
		FromTable: "PageTemplate",
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "guid", Operator: "= ?"},
				{LeftSide: "deletedAt", Operator: "IS NULL"},
			},
		},
		Limit: 1,
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), guid)
	var pageTemplate pagetemplate.PageTemplate
	err = wrapsql.GetSingleRow(guid, rows, err, &pageTemplate.ID, &pageTemplate.GUID, &pageTemplate.Name)
	return pageTemplate, err
}
