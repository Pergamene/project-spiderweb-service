package mysqlstore

import (
	"database/sql"
)

// HealthcheckStore is the mysql for pages
type HealthcheckStore struct {
	db *sql.DB
}

// NewHealthcheckStore returns a HealthcheckStore
func NewHealthcheckStore(mysqldb *sql.DB) HealthcheckStore {
	return HealthcheckStore{
		db: mysqldb,
	}
}

// IsHealthy checks if the db is healthy.
func (s HealthcheckStore) IsHealthy() (bool, error) {
	rows, err := s.db.Query("SELECT `status` FROM `healthcheck`")
	if err != nil {
		return false, err
	}
	defer rows.Close()
	var status string
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			return false, err
		}
		if status == "ok" {
			return true, nil
		}
	}
	err = rows.Err()
	return false, err
}
