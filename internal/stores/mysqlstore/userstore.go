package mysqlstore

import (
	"database/sql"
	"errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/appuser"
	"github.com/Pergamene/project-spiderweb-service/internal/util/wrapsql"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
)

// UserStore is the mysql for versions
type UserStore struct {
	db *sql.DB
}

// NewUserStore returns a UserStore
func NewUserStore(mysqldb *sql.DB) UserStore {
	return UserStore{
		db: mysqldb,
	}
}

// GetUser returns the given appuser.
func (s UserStore) GetUser(guid string) (appuser.User, error) {
	if guid == "" {
		return appuser.User{}, errors.New("must provide guid to get the user")
	}
	if s.db == nil {
		return appuser.User{}, &storeerror.DBNotSetUp{}
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"ID", "guid", "email"},
		FromTable: "User",
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "guid", Operator: "= ?"},
				{LeftSide: "deletedAt", Operator: "IS NULL"},
			},
		},
		Limit: 1,
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), guid)
	var u appuser.User
	err = wrapsql.GetSingleRow(guid, rows, err, &u.ID, &u.GUID, &u.Email)
	return u, err
}
