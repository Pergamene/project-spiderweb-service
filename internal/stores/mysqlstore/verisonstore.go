package mysqlstore

import (
	"database/sql"
	"errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/version"

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
	rows, err := s.db.Query("SELECT `ID`, `guid`, `name` FROM `Version` WHERE `guid` = ? AND `deletedAt` IS NULL LIMIT 1", guid)
	var v version.Version
	err = getSingleRow(guid, rows, err, &v.ID, &v.GUID, &v.Name)
	return v, err
}
