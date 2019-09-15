package mysqlstore

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"
	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
	"github.com/Pergamene/project-spiderweb-service/internal/models/property"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func testPageStoreClearAllTables(db *sql.DB) error {
	tables := []string{"Page", "PageOwner", "PageTemplate", "User", "Version"}
	for _, table := range tables {
		err := clearTableForTest(db, table)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCreatePage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramRecord            page.Page
		paramOwnerID           int64
		expectedPageGUID       string
		expectedDPPage         page.Page
		expectedOwnerGUID      string
		returnPage             page.Page
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
			},
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
				ID:             1,
				GUID:           "PG_1",
				Title:          "new title",
				Summary:        "new summary",
				Version:        version.Version{GUID: "VR_1"},
				PermissionType: permission.TypePrivate,
				PageTemplate:   pagetemplate.PageTemplate{GUID: "PGT_1"},
			},
			expectedOwnerGUID: "UR_1",
			returnPage: page.Page{
				ID:             1,
				GUID:           "PG_1",
				Title:          "new title",
				Summary:        "new summary",
				Version:        version.Version{ID: 1, GUID: "VR_1"},
				PermissionType: permission.TypePrivate,
				PageTemplate:   pagetemplate.PageTemplate{ID: 1, GUID: "PGT_1"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := testPageStoreClearAllTables(pageStore.db)
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
			result.CreatedAt = nil
			result.UpdatedAt = nil
			result.DeletedAt = nil
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
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
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
			err := testPageStoreClearAllTables(pageStore.db)
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
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
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
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_2\", \"bob2@test.com\", NOW(), NOW())",
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
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_2\", \"bob2@test.com\", NOW(), NOW())",
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
			err := testPageStoreClearAllTables(pageStore.db)
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

func TestUpdatePage(t *testing.T) {
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
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
			},
			paramRecord: page.Page{
				GUID:    "PG_1",
				Title:   "new title",
				Summary: "new summary",
			},
			expectedPageGUID: "PG_1",
			expectedDPPage: page.Page{
				ID: 1,
				Version: version.Version{
					GUID: "VR_1",
				},
				PageTemplate: pagetemplate.PageTemplate{
					GUID: "PGT_1",
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
			err := testPageStoreClearAllTables(pageStore.db)
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			err = pageStore.UpdatePage(tc.paramRecord)
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

func TestGetPage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		returnPage             page.Page
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
			},
			paramGUID: "PG_1",
			returnPage: page.Page{
				ID: 1,
				Version: version.Version{
					GUID: "VR_1",
				},
				PageTemplate: pagetemplate.PageTemplate{
					GUID: "PGT_1",
				},
				GUID:           "PG_1",
				Title:          "original title",
				PermissionType: permission.TypePrivate,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := testPageStoreClearAllTables(pageStore.db)
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

func TestGetPages(t *testing.T) {
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
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_2\", \"bob2@test.com\", NOW(), NOW())",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"test title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_2\", \"test title 2\", \"some kind of summary\", \"PU\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_3\", \"test title 3\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_4\", \"test title 4\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 2, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 3, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 4, 2, true)",
			},
			paramUserID:      "UR_1",
			paramThisBatchID: "",
			paramLimit:       2,
			returnPages: []page.Page{
				{
					ID:             1,
					GUID:           "PG_1",
					Version:        version.Version{GUID: "VR_1"},
					PageTemplate:   pagetemplate.PageTemplate{GUID: "PGT_1"},
					Title:          "test title",
					PermissionType: permission.TypePrivate,
				},
				{
					ID:             2,
					GUID:           "PG_2",
					Version:        version.Version{GUID: "VR_1"},
					PageTemplate:   pagetemplate.PageTemplate{GUID: "PGT_1"},
					Title:          "test title 2",
					Summary:        "some kind of summary",
					PermissionType: permission.TypePublic,
				},
			},
			returnNextBatchID: "PG_3",
			returnTotal:       3,
		},
		{
			name: "happy path, but with a next batch id used to offset request",
			preTestQueries: []string{
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_2\", \"bob2@test.com\", NOW(), NOW())",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"test title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_2\", \"test title 2\", \"some kind of summary\", \"PU\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_3\", \"test title 3\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_4\", \"test title 4\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 2, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 3, 1, true)",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 4, 2, true)",
			},
			paramUserID:      "UR_1",
			paramThisBatchID: "PG_3",
			paramLimit:       2,
			returnPages: []page.Page{
				{
					ID:             3,
					Version:        version.Version{GUID: "VR_1"},
					PageTemplate:   pagetemplate.PageTemplate{GUID: "PGT_1"},
					GUID:           "PG_3",
					Title:          "test title 3",
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
			err := testPageStoreClearAllTables(pageStore.db)
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
			for i := range pages {
				pages[i].CreatedAt = nil
				pages[i].UpdatedAt = nil
				pages[i].DeletedAt = nil
			}
			require.Equal(t, tc.returnPages, pages)
			require.Equal(t, tc.returnTotal, total)
			require.Equal(t, tc.returnNextBatchID, nextBatchID)
		})
	}
}

func TestRemovePage(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramGUID              string
		expectedPageGUID       string
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"original title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
			},
			paramGUID:        "PG_1",
			expectedPageGUID: "PG_1",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := testPageStoreClearAllTables(pageStore.db)
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			err = pageStore.RemovePage(tc.paramGUID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			if tc.expectedPageGUID == "" {
				return
			}
			_, err = pageStore.GetPage(tc.expectedPageGUID)
			if _, ok := err.(*storeerror.NotFound); !ok {
				t.Fatalf("Page %v was not deleted", tc.expectedPageGUID)
			}
		})
	}
}

func TestGetUniquePageGUID(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramProposedPageGUID  string
		returnPageGUIDPrefix   string
		returnPageGUIDLength   int
		returnErr              error
	}{
		{
			name:                 "happy path",
			preTestQueries:       []string{},
			returnPageGUIDPrefix: "PG_",
			returnPageGUIDLength: 15,
		},
		{
			name:                  "happy path, but with a proposal",
			preTestQueries:        []string{},
			paramProposedPageGUID: "PG_123456789012",
			returnPageGUIDPrefix:  "PG_",
			returnPageGUIDLength:  15,
		},
		{
			name:                  "proposal is bad",
			preTestQueries:        []string{},
			paramProposedPageGUID: "PG_123456789",
			returnErr:             errors.New("proposed guid must be 15 characters"),
		},
		{
			name: "proposal already exists",
			preTestQueries: []string{
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_123456789012\", \"original title\", \"\", \"PR\", NOW(), NOW() )",
			},
			paramProposedPageGUID: "PG_123456789012",
			returnErr:             errors.New("the proposed guid PG_123456789012 already exists"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := testPageStoreClearAllTables(pageStore.db)
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			result, err := pageStore.GetUniquePageGUID(tc.paramProposedPageGUID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnPageGUIDLength, len(result))
			require.Equal(t, tc.returnPageGUIDPrefix, result[:len(tc.returnPageGUIDPrefix)])
		})
	}
}

func TestGetPageProperties(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		paramPageGUID          string
		returnProperties       []property.Property
		returnErr              error
	}{
		{
			name: "happy path",
			preTestQueries: []string{
				"INSERT INTO Version (`guid`, `name`, `createdAt`, `updatedAt`) VALUES( \"VR_1\", \"TEST_VERSION\", NOW(), NOW())",
				"INSERT INTO PageTemplate (`Version_ID`, `guid`, `name`, `hasProperties`, `hasDetails`, `hasRelations`, `createdAt`, `updatedAt`) VALUES(1, \"PGT_1\", \"TEST_TEMPLATE\", true, true, true, NOW(), NOW())",
				"INSERT INTO User (`guid`, `email`, `createdAt`, `updatedAt`) VALUES( \"UR_1\", \"bob@test.com\", NOW(), NOW())",
				"INSERT INTO Page (`Version_ID`, `PageTemplate_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, \"PG_1\", \"test title\", \"\", \"PR\", NOW(), NOW() )",
				"INSERT INTO PageOwner (`Page_ID`, `User_ID`, `isOwner`) VALUES( 1, 1, true)",
				"INSERT INTO Property (`Version_ID`, `type`, `key`, `createdAt`, `updatedAt`) VALUES( 1, \"NU\", \"population\", NOW(), NOW())",
				"INSERT INTO Property (`Version_ID`, `type`, `key`, `createdAt`, `updatedAt`) VALUES( 1, \"ST\", \"banner\", NOW(), NOW())",
				"INSERT INTO Property (`Version_ID`, `type`, `key`, `createdAt`, `updatedAt`) VALUES( 1, \"ST\", \"color\", NOW(), NOW())",
				"INSERT INTO Property (`Version_ID`, `type`, `key`, `createdAt`, `updatedAt`) VALUES( 1, \"ST\", \"symbol\", NOW(), NOW())",
				"INSERT INTO PagePropertyNumber (`Page_ID`, `Property_ID`, `Version_ID`, `value`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 1, 1, 100000, \"PR\", NOW(), NOW())",
				"INSERT INTO PagePropertyString (`Page_ID`, `Property_ID`, `Version_ID`, `value`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 2, 1, \"lion heads\", \"PR\", NOW(), NOW())",
				"INSERT INTO PagePropertyString (`Page_ID`, `Property_ID`, `Version_ID`, `value`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 3, 1, \"blue\", \"PR\", NOW(), NOW())",
				"INSERT INTO PagePropertyString (`Page_ID`, `Property_ID`, `Version_ID`, `value`, `permission`, `createdAt`, `updatedAt`) VALUES( 1, 4, 1, \"lion\", \"PR\", NOW(), NOW())",
				"INSERT INTO PagePropertyOrder (`Page_ID`, `Property_ID`, `order`) VALUES( 1, 1, 1)",
				"INSERT INTO PagePropertyOrder (`Page_ID`, `Property_ID`, `order`) VALUES( 1, 2, 2)",
				"INSERT INTO PagePropertyOrder (`Page_ID`, `Property_ID`, `order`) VALUES( 1, 3, 0)",
				"INSERT INTO PagePropertyOrder (`Page_ID`, `Property_ID`, `order`) VALUES( 1, 4, 3)",
			},
			paramPageGUID: "PG_1",
			returnProperties: []property.Property{
				property.Property{
					ID:    1,
					Key:   "test",
					Type:  property.TypeString,
					Value: "test",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := PageStore{
				db: mysqldb,
			}
			err := testPageStoreClearAllTables(pageStore.db)
			require.NoError(t, err)
			err = execPreTestQueries(pageStore.db, tc.preTestQueries)
			require.NoError(t, err)
			if tc.shouldReplaceDBWithNil {
				pageStore.db = nil
			}
			pageProperties, err := pageStore.GetPageProperties(tc.paramPageGUID)
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnProperties, pageProperties)
		})
	}
}
