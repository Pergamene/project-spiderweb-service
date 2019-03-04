package mysqlstore

import (
	"fmt"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"
	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func TestPageUpdatePage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramRecord            page.Page
		pageGUID               string
		expectedDPPage         page.Page
		returnErr              error
	}{
		{
			name:           "good update",
			preTestQueries: []string{"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )"},
			paramRecord: page.Page{
				GUID:    "PG_1",
				Title:   "new title",
				Summary: "new summary",
			},
			pageGUID: "PG_1",
			expectedDPPage: page.Page{
				Version: version.Version{
					ID: 1,
				},
				PageTemplate: pagetemplate.PageTemplate{
					ID: 1,
				},
				GUID:           "PG_1",
				Title:          "new title",
				Summary:        "new summary",
				PermissionType: permission.TypePrivate,
			},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			err := clearTableForTest(pageStore.db, "Page")
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			err = pageStore.UpdatePage(tc.paramRecord)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			p, err := pageStore.GetPage(tc.pageGUID)
			require.NoError(t, err)
			p.CreatedAt = nil
			p.UpdatedAt = nil
			p.DeletedAt = nil
			require.Equal(t, tc.expectedDPPage, p)
		})
	}
}
