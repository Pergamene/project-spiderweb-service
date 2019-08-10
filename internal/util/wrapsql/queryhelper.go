package wrapsql

import (
	"database/sql"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
)

func GetSingleRow(guid string, rows *sql.Rows, queryErr error, columns ...interface{}) error {
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

func ExecSingleInsert(db *sql.DB, iq InsertQuery) (lastInsertID int64, err error) {
	var statement *sql.Stmt
	var result sql.Result
	queryString, orderedValues := GetInsertString(iq)
	statement, err = db.Prepare(queryString)
	if err != nil {
		return
	}
	defer statement.Close()
	result, err = statement.Exec(orderedValues...)
	if err != nil {
		return
	}
	lastInsertID, err = result.LastInsertId()
	return
}
