package pagehandler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
	"github.com/pkg/errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/handlertestutils"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/page/mocks"
	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	pageservice "github.com/Pergamene/project-spiderweb-service/internal/services/page"
)

func getPage(guid string, versionID string, permissionType permission.Type) page.Page {
	return page.Page{
		GUID:    guid,
		Title:   "test title",
		Summary: "test summary",
		Version: version.Version{
			GUID: versionID,
		},
		PermissionType: permissionType,
	}
}

type createPageCall struct {
	pageParams   pageservice.CreatePageParams
	returnRecord page.Page
	returnErr    error
}

func TestCreatePage(t *testing.T) {
	cases := []struct {
		name                 string
		headers              map[string]string
		params               url.Values
		requestBody          string
		authN                api.AuthN
		authZ                api.AuthZ
		datacenter           string
		expectedResponseBody string
		expectedStatusCode   int
		createPageCalls      []createPageCall
	}{
		{
			name:                 "not authenticated",
			authN:                handlertestutils.DefaultAuthN("PROD"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
		{
			name: "happy page creation, local",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody: "{\"title\":\"test title\",\"summary\":\"test summary\",\"versionId\":1,\"permission\":\"PR\"}",
			// params: url.Values{
			// 	"title":      []string{"test title"},
			// 	"summary":    []string{"test summary"},
			// 	"versionId":  []string{"1"},
			// 	"permission": []string{"PR"},
			// },
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"id\":\"PG_1\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			createPageCalls: []createPageCall{
				{
					pageParams: pageservice.CreatePageParams{
						Page:    getPage("", "VR_1", permission.TypePrivate),
						OwnerID: "UR_1",
					},
					returnRecord: getPage("PG_1", "VR_1", permission.TypePrivate),
				},
			},
		},
		{
			name: "missing title for the page",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:          "{\"summary\":\"test summary\",\"versionId\":1,\"permission\":\"PR\"}",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"400 - Bad Request\",\"message\":\"must provide title\"}}\n",
			expectedStatusCode:   400,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageService := new(mocks.PageService)
			for index := range tc.createPageCalls {
				pageService.On("CreatePage", mock.Anything, tc.createPageCalls[index].pageParams).Return(tc.createPageCalls[index].returnRecord, tc.createPageCalls[index].returnErr)
			}
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         http.MethodPost,
				Endpoint:       "pages",
				Params:         tc.params,
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "CreatePage", len(tc.createPageCalls))
		})
	}
}

type updatePageCall struct {
	pageParams pageservice.SetPageParams
	returnErr  error
}

func TestUpdatePage(t *testing.T) {
	cases := []struct {
		name                 string
		pageID               string
		headers              map[string]string
		requestBody          string
		authN                api.AuthN
		authZ                api.AuthZ
		datacenter           string
		expectedResponseBody string
		expectedStatusCode   int
		updatePageCalls      []updatePageCall
	}{
		{
			name:                 "not authenticated",
			pageID:               "PG_1",
			authN:                handlertestutils.DefaultAuthN("PROD"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
		{
			name:   "happy page creation, local",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:          "{\"title\":\"test title\",\"summary\":\"test summary\",\"versionId\":1,\"permission\":\"PR\"}",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			updatePageCalls: []updatePageCall{
				{
					pageParams: pageservice.SetPageParams{
						Page:   getPage("PG_1", "VR_1", ""),
						UserID: "UR_1",
					},
				},
			},
		},
		{
			name:   "trying to edit a page that you don't have permission to update",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:          "{\"title\":\"test title\",\"summary\":\"test summary\",\"versionId\":1,\"permission\":\"PR\"}",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authorized\"}}\n",
			expectedStatusCode:   401,
			updatePageCalls: []updatePageCall{
				{
					pageParams: pageservice.SetPageParams{
						Page:   getPage("PG_1", "VR_1", ""),
						UserID: "UR_1",
					},
					returnErr: &storeerror.NotAuthorized{
						UserID:  "UR_1",
						TableID: "PG_1",
						Err:     errors.New("failure"),
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageService := new(mocks.PageService)
			for index := range tc.updatePageCalls {
				pageService.On("UpdatePage", mock.Anything, tc.updatePageCalls[index].pageParams).Return(tc.updatePageCalls[index].returnErr)
			}
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         http.MethodPut,
				Endpoint:       fmt.Sprintf("pages/%v", tc.pageID),
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "UpdatePage", len(tc.updatePageCalls))
		})
	}
}
