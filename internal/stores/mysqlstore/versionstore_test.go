package mysqlstore

import (
	"database/sql"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/version"

	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func testVersionStoreClearAllTables(db *sql.DB) error {
	tables := []string{"Version"}
	for _, table := range tables {
		err := clearTableForTest(db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestGetVersion(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		returnVersion          version.Version
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
			},
			paramGUID: "VR_1",
			returnVersion: version.Version{
				ID:   1,
				GUID: "VR_1",
				Name: "TEST_VERSION",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			versionStore := VersionStore{
				db: mysqldb,
			}
			err := testVersionStoreClearAllTables(versionStore.db)
			require.NoError(t, err)
			err = execPreTestQueries(versionStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				versionStore.db = nil
			}
			result, err := versionStore.GetVersion(tc.paramGUID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnVersion, result)
		})
	}
}
