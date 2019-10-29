package mysqlstore

import (
	"database/sql"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
	"github.com/Pergamene/project-spiderweb-service/internal/util/guidgen"
	"github.com/Pergamene/project-spiderweb-service/internal/util/wrapsql"
	"github.com/pkg/errors"
)

func getUniqueGUID(db *sql.DB, prefix string, length int, table, proposedGUID string, retry int) (string, error) {
	guid := proposedGUID
	if guid == "" {
		guid = guidgen.GenerateGUID("PG", 15)
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"guid"},
		FromTable: table,
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "guid", Operator: "= ?"},
			},
		},
		Limit: 1,
	}
	rows, err := db.Query(wrapsql.GetSelectString(statement), guid)
	var resGUID string
	err = wrapsql.GetSingleRow(guid, rows, err, &resGUID)
	if err == nil {
		if proposedGUID != "" {
			return "", errors.Errorf("the proposed guid %v already exists", proposedGUID)
		}
		if retry >= guidgen.MaxGUIDRetryAttempts {
			return "", guidgen.ErrMaxGUIDRetryAttempts
		}
		return getUniqueGUID(db, prefix, length, table, proposedGUID, retry+1)
	}
	if _, ok := err.(*storeerror.NotFound); ok {
		return guid, nil
	}
	return "", err
}
