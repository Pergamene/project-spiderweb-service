package pageservice

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/store/mocks"
)

var pageService PageService
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()
	result := m.Run()
	os.Exit(result)
}

func getPage(guid, title, summary string) page.Page {
	return page.Page{
		GUID:    guid,
		Title:   title,
		Summary: summary,
	}
}

func getStoreUnauthorizedErr(userID, tableID string, err error) error {
	return &storeerror.NotAuthorized{
		UserID:  userID,
		TableID: tableID,
		Err:     err,
	}
}

type getPageTemplateCall struct {
	paramPageTemplateGUID string
	returnPageTemplate    pagetemplate.PageTemplate
	returnErr             error
}

type getVersionCall struct {
	paramVersionGUID string
	returnVersion    version.Version
	returnErr        error
}

type canEditPageCall struct {
	paramPageGUID   string
	paramPageUserID string
	returnIsOwner   bool
	returnErr       error
}

type setPageCall struct {
	paramPage page.Page
	returnErr error
}

func TestSetPage(t *testing.T) {
	cases := []struct {
		name                 string
		params               SetPageParams
		canEditPageCalls     []canEditPageCall
		getPageTemplateCalls []getPageTemplateCall
		getVersionCalls      []getVersionCall
		setPageCalls         []setPageCall
		returnErr            error
	}{
		{
			name: "test happy path",
			params: SetPageParams{
				Page: page.Page{
					GUID:  "PG_1",
					Title: "New Title",
				},
				UserID: "UR_1",
			},
			canEditPageCalls: []canEditPageCall{
				{
					paramPageGUID:   "PG_1",
					paramPageUserID: "UR_1",
				},
			},
			setPageCalls: []setPageCall{{paramPage: page.Page{
				GUID:  "PG_1",
				Title: "New Title",
			}}},
		},
		{
			name: "test update of version and page template",
			params: SetPageParams{
				Page: page.Page{
					GUID:         "PG_1",
					Title:        "New Title",
					PageTemplate: pagetemplate.PageTemplate{GUID: "PGT_1"},
					Version:      version.Version{GUID: "VR_1"},
				},
				UserID: "UR_1",
			},
			canEditPageCalls: []canEditPageCall{
				{
					paramPageGUID:   "PG_1",
					paramPageUserID: "UR_1",
				},
			},
			getPageTemplateCalls: []getPageTemplateCall{
				{
					paramPageTemplateGUID: "PGT_1",
					returnPageTemplate:    pagetemplate.PageTemplate{GUID: "PGT_1", ID: 1, Name: "TEST_NAME_TEMPLATE"},
				},
			},
			getVersionCalls: []getVersionCall{
				{
					paramVersionGUID: "VR_1",
					returnVersion:    version.Version{GUID: "VR_1", ID: 1, Name: "TEST_NAME_VERSION"},
				},
			},
			setPageCalls: []setPageCall{{paramPage: page.Page{
				GUID:         "PG_1",
				Title:        "New Title",
				PageTemplate: pagetemplate.PageTemplate{GUID: "PGT_1", ID: 1, Name: "TEST_NAME_TEMPLATE"},
				Version:      version.Version{GUID: "VR_1", ID: 1, Name: "TEST_NAME_VERSION"},
			}}},
		},
		{
			name: "test unauthorized call",
			params: SetPageParams{
				Page: page.Page{
					GUID:  "PG_1",
					Title: "New Title",
				},
				UserID: "UR_1",
			},
			canEditPageCalls: []canEditPageCall{
				{
					paramPageGUID:   "PG_1",
					paramPageUserID: "UR_1",
					returnErr:       getStoreUnauthorizedErr("UR_1", "PG_1", nil),
				},
			},
			returnErr: errors.New("User UR_1 is not authorized to perform the action on the ID PG_1"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := new(mocks.PageStore)
			pageTemplateStore := new(mocks.PageTemplateStore)
			versionStore := new(mocks.VersionStore)
			for index := range tc.getPageTemplateCalls {
				pageTemplateStore.On("GetPageTemplate", tc.getPageTemplateCalls[index].paramPageTemplateGUID).Return(tc.getPageTemplateCalls[index].returnPageTemplate, tc.getPageTemplateCalls[index].returnErr)
			}
			for index := range tc.getVersionCalls {
				versionStore.On("GetVersion", tc.getVersionCalls[index].paramVersionGUID).Return(tc.getVersionCalls[index].returnVersion, tc.getVersionCalls[index].returnErr)
			}
			for index := range tc.canEditPageCalls {
				pageStore.On("CanEditPage", tc.canEditPageCalls[index].paramPageGUID, tc.canEditPageCalls[index].paramPageUserID).Return(tc.canEditPageCalls[index].returnIsOwner, tc.canEditPageCalls[index].returnErr)
			}
			for index := range tc.setPageCalls {
				pageStore.On("SetPage", tc.setPageCalls[index].paramPage).Return(tc.setPageCalls[index].returnErr)
			}
			pageService = PageService{
				PageStore:         pageStore,
				PageTemplateStore: pageTemplateStore,
				VersionStore:      versionStore,
			}
			err := pageService.SetPage(ctx, tc.params)
			pageTemplateStore.AssertNumberOfCalls(t, "GetPageTemplate", len(tc.getPageTemplateCalls))
			versionStore.AssertNumberOfCalls(t, "GetVersion", len(tc.getVersionCalls))
			pageStore.AssertNumberOfCalls(t, "CanEditPage", len(tc.canEditPageCalls))
			pageStore.AssertNumberOfCalls(t, "SetPage", len(tc.setPageCalls))
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
		})
	}
}

type getUniquePageGUIDCall struct {
	paramProposedGUID string
	returnGUID        string
	returnErr         error
}

type createPageCall struct {
	paramPage    page.Page
	paramOwnerID string
	returnPage   page.Page
	returnErr    error
}

func TestCreatePage(t *testing.T) {
	cases := []struct {
		name                   string
		params                 CreatePageParams
		getPageTemplateCalls   []getPageTemplateCall
		getVersionCalls        []getVersionCall
		getUniquePageGUIDCalls []getUniquePageGUIDCall
		createPageCalls        []createPageCall
		returnPage             page.Page
		returnErr              error
	}{
		{
			name: "test happy path",
			params: CreatePageParams{
				Page: page.Page{
					Title:        "New Title",
					PageTemplate: pagetemplate.PageTemplate{GUID: "PGT_1"},
					Version:      version.Version{GUID: "VR_1"},
				},
				OwnerID: "UR_1",
			},
			getPageTemplateCalls: []getPageTemplateCall{
				{
					paramPageTemplateGUID: "PGT_1",
					returnPageTemplate:    pagetemplate.PageTemplate{GUID: "PGT_1", ID: 1, Name: "TEST_NAME_TEMPLATE"},
				},
			},
			getVersionCalls: []getVersionCall{
				{
					paramVersionGUID: "VR_1",
					returnVersion:    version.Version{GUID: "VR_1", ID: 1, Name: "TEST_NAME_VERSION"},
				},
			},
			getUniquePageGUIDCalls: []getUniquePageGUIDCall{
				{
					returnGUID: "PG_NEW",
				},
			},
			createPageCalls: []createPageCall{
				{
					paramPage: page.Page{
						GUID:         "PG_NEW",
						Title:        "New Title",
						PageTemplate: pagetemplate.PageTemplate{GUID: "PGT_1", ID: 1, Name: "TEST_NAME_TEMPLATE"},
						Version:      version.Version{GUID: "VR_1", ID: 1, Name: "TEST_NAME_VERSION"},
					},
					paramOwnerID: "UR_1",
					returnPage: page.Page{
						ID:           1,
						GUID:         "PG_NEW",
						Title:        "New Title",
						PageTemplate: pagetemplate.PageTemplate{GUID: "PGT_1", ID: 1, Name: "TEST_NAME_TEMPLATE"},
						Version:      version.Version{GUID: "VR_1", ID: 1, Name: "TEST_NAME_VERSION"},
					},
				},
			},
			returnPage: page.Page{
				ID:           1,
				GUID:         "PG_NEW",
				Title:        "New Title",
				PageTemplate: pagetemplate.PageTemplate{GUID: "PGT_1", ID: 1, Name: "TEST_NAME_TEMPLATE"},
				Version:      version.Version{GUID: "VR_1", ID: 1, Name: "TEST_NAME_VERSION"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageStore := new(mocks.PageStore)
			pageTemplateStore := new(mocks.PageTemplateStore)
			versionStore := new(mocks.VersionStore)
			for index := range tc.getPageTemplateCalls {
				pageTemplateStore.On("GetPageTemplate", tc.getPageTemplateCalls[index].paramPageTemplateGUID).Return(tc.getPageTemplateCalls[index].returnPageTemplate, tc.getPageTemplateCalls[index].returnErr)
			}
			for index := range tc.getVersionCalls {
				versionStore.On("GetVersion", tc.getVersionCalls[index].paramVersionGUID).Return(tc.getVersionCalls[index].returnVersion, tc.getVersionCalls[index].returnErr)
			}
			for index := range tc.getUniquePageGUIDCalls {
				pageStore.On("GetUniquePageGUID", tc.getUniquePageGUIDCalls[index].paramProposedGUID).Return(tc.getUniquePageGUIDCalls[index].returnGUID, tc.getUniquePageGUIDCalls[index].returnErr)
			}
			for index := range tc.createPageCalls {
				pageStore.On("CreatePage", tc.createPageCalls[index].paramPage, tc.createPageCalls[index].paramOwnerID).Return(tc.createPageCalls[index].returnPage, tc.createPageCalls[index].returnErr)
			}
			pageService = PageService{
				PageStore:         pageStore,
				PageTemplateStore: pageTemplateStore,
				VersionStore:      versionStore,
			}
			result, err := pageService.CreatePage(ctx, tc.params)
			pageTemplateStore.AssertNumberOfCalls(t, "GetPageTemplate", len(tc.getPageTemplateCalls))
			versionStore.AssertNumberOfCalls(t, "GetVersion", len(tc.getVersionCalls))
			pageStore.AssertNumberOfCalls(t, "GetUniquePageGUID", len(tc.getUniquePageGUIDCalls))
			pageStore.AssertNumberOfCalls(t, "CreatePage", len(tc.createPageCalls))
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnPage, result)
		})
	}
}
