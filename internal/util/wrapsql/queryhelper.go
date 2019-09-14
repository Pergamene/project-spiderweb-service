package wrapsql

import (
	"database/sql"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
)

// GetSingleRow extracts the given sql.Rows to return a single row scanned into the given columns
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

// ExecSingleInsert executes a single INSERT command and returns the lastInsertID
func ExecSingleInsert(db *sql.DB, query InsertQuery) (lastInsertID int64, err error) {
	var statement *sql.Stmt
	var result sql.Result
	queryString, orderedValues := GetInsertString(query)
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

// ExecBatchInsert executes a batch INSERT command
func ExecBatchInsert(db *sql.DB, query BatchInsertQuery) (err error) {
	var statement *sql.Stmt
	queryString, orderedValues := GetBatchInsertString(query)
	statement, err = db.Prepare(queryString)
	if err != nil {
		return
	}
	defer statement.Close()
	_, err = statement.Exec(orderedValues...)
	if err != nil {
		return
	}
	return
}

// ExecSingleUpdate executes a single UPDATE command
func ExecSingleUpdate(db *sql.DB, query UpdateQuery, whereClauseInjectedValues ...interface{}) (err error) {
	var statement *sql.Stmt
	queryString, orderedValues := GetUpdateString(query, whereClauseInjectedValues...)
	statement, err = db.Prepare(queryString)
	if err != nil {
		return
	}
	defer statement.Close()
	_, err = statement.Exec(orderedValues...)
	return
}

// ExecDelete executes a DELETE command
func ExecDelete(db *sql.DB, query DeleteQuery, whereClauseInjectedValues ...interface{}) (err error) {
	var statement *sql.Stmt
	queryString, orderedValues := GetDeleteString(query, whereClauseInjectedValues...)
	statement, err = db.Prepare(queryString)
	if err != nil {
		return
	}
	defer statement.Close()
	_, err = statement.Exec(orderedValues...)
	return
}
