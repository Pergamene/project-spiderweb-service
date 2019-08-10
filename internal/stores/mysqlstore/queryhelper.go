package mysqlstore

import (
	"database/sql"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
)

func getSingleRow(guid string, rows *sql.Rows, queryErr error, columns ...interface{}) error {
	if queryErr != nil {
		return queryErr
	}
	if err := rows.Err(); err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(columns...)
		if err != nil {
			return err
		}
		return nil
	}
	return &storeerror.NotFound{
		ID: guid,
	}
}
