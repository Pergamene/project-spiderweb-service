package mysqlstore

import (
	"database/sql"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/appuser"

	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func testUserStoreClearAllTables(db *sql.DB) error {
	tables := []string{"User"}
	for _, table := range tables {
		err := clearTableForTest(db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestGetUser(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		returnUser             appuser.User
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
			},
			paramGUID: "UR_1",
			returnUser: appuser.User{
				ID:    1,
				GUID:  "UR_1",
				Email: "bob@test.com",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userStore := UserStore{
				db: mysqldb,
			}
			err := testUserStoreClearAllTables(userStore.db)
			require.NoError(t, err)
			err = execPreTestQueries(userStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				userStore.db = nil
			}
			result, err := userStore.GetUser(tc.paramGUID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnUser, result)
		})
	}
}
