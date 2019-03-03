package storeerror

import "fmt"

// DBNotSetUp is an error that signifies that the db is not configured correctly to query.
type DBNotSetUp struct {
	Err error
}

func (e *DBNotSetUp) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("DB is not configured\n%v", e.Err)
	}
	return "DB is not configured"
}
