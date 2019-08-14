package mysqlstore

import (
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"
	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func TestCreatePage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramRecord            page.Page
		paramOwnerID           int
		expectedPageGUID       string
		expectedDPPage         page.Page
		expectedOwnerGUID      string
		returnPage             page.Page
		returnErr              error
	}{
		{
			name:           "happy path",
			preTestQueries: []string{"INSERT INTO User (`guid`, `email`) VALUES( \"UR_1\", \"bob@test.com\")"},
			paramRecord: page.Page{
				GUID:           "PG_1",
				Title:          "new title",
				Summary:        "new summary",
				Version:        version.Version{ID: 1, GUID: "VR_1"},
				PermissionType: permission.TypePrivate,
				PageTemplate:   pagetemplate.PageTemplate{ID: 1, GUID: "PGT_1"},
			},
			paramOwnerID:     1,
			expectedPageGUID: "PG_1",
			expectedDPPage: page.Page{
				GUID:           "PG_1",
				Title:          "new title",
				Summary:        "new summary",
				Version:        version.Version{ID: 1},
				PermissionType: permission.TypePrivate,
				PageTemplate:   pagetemplate.PageTemplate{ID: 1},
			},
			expectedOwnerGUID: "UR_1",
			returnPage: page.Page{
				GUID:           "PG_1",
				Title:          "new title",
				Summary:        "new summary",
				Version:        version.Version{ID: 1},
				PermissionType: permission.TypePrivate,
				PageTemplate:   pagetemplate.PageTemplate{ID: 1},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := clearTableForTest(pageStore.db, "Page")
			require.NoError(t, err)
			err = clearTableForTest(pageStore.db, "User")
			require.NoError(t, err)
			err = clearTableForTest(pageStore.db, "PageOwner")
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			result, err := pageStore.CreatePage(tc.paramRecord, tc.paramOwnerID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnPage, result)
			p, err := pageStore.GetPage(tc.expectedPageGUID)
			require.NoError(t, err)
			p.CreatedAt = nil
			p.UpdatedAt = nil
			p.DeletedAt = nil
			require.Equal(t, tc.expectedDPPage, p)
			isOwner, err := pageStore.CanEditPage(tc.expectedPageGUID, tc.expectedOwnerGUID)
			require.NoError(t, err)
			require.Equal(t, true, isOwner)
		})
	}
}

func TestCanEditPage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		paramUserID            string
		returnCanEdit          bool
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_1\", \"bob@test.com\")",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 1, true)",
			},
			paramGUID:     "PG_1",
			paramUserID:   "UR_1",
			returnCanEdit: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := clearTableForTest(pageStore.db, "Page")
			require.NoError(t, err)
			err = clearTableForTest(pageStore.db, "User")
			require.NoError(t, err)
			err = clearTableForTest(pageStore.db, "PageOwner")
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			canEdit, err := pageStore.CanEditPage(tc.paramGUID, tc.paramUserID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnCanEdit, canEdit)
		})
	}
}

func TestCanReadPage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		paramUserID            string
		returnCanRead          bool
		returnErr              error
	}{
		{
			name: "happy path, is owner",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_1\", \"bob@test.com\")",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 1, true)",
			},
			paramGUID:     "PG_1",
			paramUserID:   "UR_1",
			returnCanRead: true,
		},
		{
			name: "happy path, not owner but public",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PU\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_1\", \"bob@test.com\")",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_2\", \"bob2@test.com\")",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 2, true)",
			},
			paramGUID:     "PG_1",
			paramUserID:   "UR_1",
			returnCanRead: true,
		},
		{
			name: "happy path, not owner and private: can't read",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_1\", \"bob@test.com\")",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_2\", \"bob2@test.com\")",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 2, true)",
			},
			paramGUID:     "PG_1",
			paramUserID:   "UR_1",
			returnCanRead: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := clearTableForTest(pageStore.db, "Page")
			require.NoError(t, err)
			err = clearTableForTest(pageStore.db, "User")
			require.NoError(t, err)
			err = clearTableForTest(pageStore.db, "PageOwner")
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			canRead, err := pageStore.CanReadPage(tc.paramGUID, tc.paramUserID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnCanRead, canRead)
		})
	}
}

func TestSetPage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramRecord            page.Page
		expectedPageGUID       string
		expectedDPPage         page.Page
		returnErr              error
	}{
		{
			name:           "happy path",
			preTestQueries: []string{"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )"},
			paramRecord: page.Page{
				GUID:    "PG_1",
				Title:   "new title",
				Summary: "new summary",
			},
			expectedPageGUID: "PG_1",
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
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := clearTableForTest(pageStore.db, "Page")
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			err = pageStore.SetPage(tc.paramRecord)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			p, err := pageStore.GetPage(tc.expectedPageGUID)
			require.NoError(t, err)
			p.CreatedAt = nil
			p.UpdatedAt = nil
			p.DeletedAt = nil
			require.Equal(t, tc.expectedDPPage, p)
		})
	}
}

func GetPage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		returnPage             page.Page
		returnErr              error
	}{
		{
			name:           "happy path",
			preTestQueries: []string{"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )"},
			paramGUID:      "PG_1",
			returnPage: page.Page{
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
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := clearTableForTest(pageStore.db, "Page")
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			p, err := pageStore.GetPage(tc.paramGUID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			p.CreatedAt = nil
			p.UpdatedAt = nil
			p.DeletedAt = nil
			require.Equal(t, tc.returnPage, p)
		})
	}
}

func GetPages(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramUserID            string
		paramThisBatchID       string
		paramLimit             int
		returnPages            []page.Page
		returnTotal            int
		returnNextBatchID      string
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"test title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_2\", \"test title 2\", \"some kind of summary\", \"PU\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_3\", \"test title 3\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_4\", \"test title 4\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_1\", \"bob@test.com\")",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_2\", \"bob2@test.com\")",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 2, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 3, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 4, 2, true)",
			},
			paramUserID:      "PG_1",
			paramThisBatchID: "",
			paramLimit:       2,
			returnPages: []page.Page{
				{
					Version:        version.Version{ID: 1},
					PageTemplate:   pagetemplate.PageTemplate{ID: 1},
					GUID:           "PG_1",
					Title:          "test title",
					PermissionType: permission.TypePrivate,
				},
				{
					Version:        version.Version{ID: 1},
					PageTemplate:   pagetemplate.PageTemplate{ID: 1},
					GUID:           "PG_2",
					Title:          "test title 2",
					PermissionType: permission.TypePrivate,
				},
			},
			returnNextBatchID: "PG_3",
			returnTotal:       3,
		},
		{
			name: "happy path, but with a next batch id used to offset request",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"test title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_2\", \"test title 2\", \"some kind of summary\", \"PU\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_3\", \"test title 3\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_4\", \"test title 4\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_1\", \"bob@test.com\")",
				"INSERT INTO User (`guid`, `email`) VALUES( \"UR_2\", \"bob2@test.com\")",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 2, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 3, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 4, 2, true)",
			},
			paramUserID:      "PG_1",
			paramThisBatchID: "PG_3",
			paramLimit:       2,
			returnPages: []page.Page{
				{
					Version:        version.Version{ID: 1},
					PageTemplate:   pagetemplate.PageTemplate{ID: 1},
					GUID:           "PG_3",
					Title:          "test title",
					PermissionType: permission.TypePrivate,
				},
			},
			returnTotal: 3,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := clearTableForTest(pageStore.db, "Page")
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			pages, total, nextBatchID, err := pageStore.GetPages(tc.paramUserID, tc.paramThisBatchID, tc.paramLimit)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			for _, p := range pages {
				p.CreatedAt = nil
				p.UpdatedAt = nil
				p.DeletedAt = nil
			}
			require.Equal(t, tc.returnPages, pages)
			require.Equal(t, tc.returnTotal, total)
			require.Equal(t, tc.returnNextBatchID, nextBatchID)
		})
	}
}
