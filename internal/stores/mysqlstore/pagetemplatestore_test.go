package mysqlstore

import (
	"database/sql"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"

	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func testPageTemplateStoreClearAllTables(db *sql.DB) error {
	tables := []string{"PageTemplate"}
	for _, table := range tables {
		err := clearTableForTest(db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestGetPageTemplate(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		returnPageTemplate     pagetemplate.PageTemplate
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
			},
			paramGUID: "PGT_1",
			returnPageTemplate: pagetemplate.PageTemplate{
				ID:   1,
				GUID: "PGT_1",
				Name: "TEST_TEMPLATE",
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageTemplateStore := PageTemplateStore{
				db: mysqldb,
			}
			err := testPageTemplateStoreClearAllTables(pageTemplateStore.db)
			require.NoError(t, err)
			err = execPreTestQueries(pageTemplateStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageTemplateStore.db = nil
			}
			result, err := pageTemplateStore.GetPageTemplate(tc.paramGUID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnPageTemplate, result)
		})
	}
}
