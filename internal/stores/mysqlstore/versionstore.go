package mysqlstore

import (
	"database/sql"
	"errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
	"github.com/Pergamene/project-spiderweb-service/internal/util/wrapsql"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
)

// VersionStore is the mysql for versions
type VersionStore struct {
	db *sql.DB
}

// NewVersionStore returns a VersionStore
func NewVersionStore(mysqldb *sql.DB) VersionStore {
	return VersionStore{
		db: mysqldb,
	}
}

// GetVersion returns the given version.
func (s VersionStore) GetVersion(guid string) (version.Version, error) {
	if guid == "" {
		return version.Version{}, errors.New("must provide guid to get the version")
	}
	if s.db == nil {
		return version.Version{}, &storeerror.DBNotSetUp{}
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"ID", "guid", "name"},
		FromTable: "Version",
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "guid", Operator: "= ?"},
				{LeftSide: "deletedAt", Operator: "IS NULL"},
			},
		},
		Limit: 1,
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), guid)
	var v version.Version
	err = wrapsql.GetSingleRow(guid, rows, err, &v.ID, &v.GUID, &v.Name)
	return v, err
}
