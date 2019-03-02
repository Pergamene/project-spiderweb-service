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

	"github.com/stretchr/testify/require"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/handlerutils"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/page/mocks"
	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
	pageservice "github.com/Pergamene/project-spiderweb-service/internal/services/page"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/mock"
)

func getPage(guid string, versionID int64, permissionType permission.Type) page.Page {
	return page.Page{
		GUID:    guid,
		Title:   "test title",
		Summary: "test summary",
		Version: version.Version{
			ID: versionID,
		},
		PermissionType: permissionType,
	}
}

func TestCreatePage(t *testing.T) {
	cases := []struct {
		name                          string
		headers                       map[string]string
		params                        url.Values
		requestBody                   string
		authN                         api.AuthN
		authZ                         api.AuthZ
		datacenter                    string
		expectedResponseBody          string
		expectedStatusCode            int
		serviceCreatePageCalled       bool
		serviceCreatePageParams       pageservice.CreatePageParams
		serviceCreatePageReturnRecord page.Page
		serviceCreatePageReturnErr    error
	}{
		{
			name:                 "not authenticated",
			authN:                handlerutils.DefaultAuthN("PROD"),
			authZ:                handlerutils.DefaultAuthZ(),
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
			authN:                   handlerutils.DefaultAuthN("LOCAL"),
			authZ:                   handlerutils.DefaultAuthZ(),
			expectedResponseBody:    "{\"result\":{\"id\":\"PG_1\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:      200,
			serviceCreatePageCalled: true,
			serviceCreatePageParams: pageservice.CreatePageParams{
				Page:    getPage("", 1, permission.TypePrivate),
				OwnerID: "UR_1",
			},
			serviceCreatePageReturnRecord: getPage("PG_1", 1, permission.TypePrivate),
		},
		{
			name: "missing title for the page",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:          "{\"summary\":\"test summary\",\"versionId\":1,\"permission\":\"PR\"}",
			authN:                handlerutils.DefaultAuthN("LOCAL"),
			authZ:                handlerutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"400 - Bad Request\",\"message\":\"must provide title\"}}\n",
			expectedStatusCode:   400,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageService := new(mocks.PageService)
			pageService.On("CreatePage", mock.Anything, tc.serviceCreatePageParams).Return(tc.serviceCreatePageReturnRecord, tc.serviceCreatePageReturnErr)
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlerutils.HandleTestRequest(handlerutils.HandleTestRequestParams{
				Method:         http.MethodPost,
				Endpoint:       "page",
				Params:         tc.params,
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "CreatePage", testutils.GetExpectedNumberOfCalls(tc.serviceCreatePageCalled))
		})
	}
}

func TestUpdatePage(t *testing.T) {
	cases := []struct {
		name                       string
		pageID                     string
		headers                    map[string]string
		requestBody                string
		authN                      api.AuthN
		authZ                      api.AuthZ
		datacenter                 string
		expectedResponseBody       string
		expectedStatusCode         int
		serviceUpdatePageCalled    bool
		serviceUpdatePageParams    pageservice.UpdatePageParams
		serviceUpdatePageReturnErr error
	}{
		{
			name:                 "not authenticated",
			pageID:               "PG_1",
			authN:                handlerutils.DefaultAuthN("PROD"),
			authZ:                handlerutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
		{
			name:   "happy page creation, local",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:             "{\"title\":\"test title\",\"summary\":\"test summary\",\"versionId\":1,\"permission\":\"PR\"}",
			authN:                   handlerutils.DefaultAuthN("LOCAL"),
			authZ:                   handlerutils.DefaultAuthZ(),
			expectedResponseBody:    "{\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:      200,
			serviceUpdatePageCalled: true,
			serviceUpdatePageParams: pageservice.UpdatePageParams{
				Page:   getPage("PG_1", 0, ""),
				UserID: "UR_1",
			},
		},
		{
			name:   "trying to edit a page that you don't have permission to update",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:             "{\"title\":\"test title\",\"summary\":\"test summary\",\"versionId\":1,\"permission\":\"PR\"}",
			authN:                   handlerutils.DefaultAuthN("LOCAL"),
			authZ:                   handlerutils.DefaultAuthZ(),
			expectedResponseBody:    "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"unauthorized error\"}}\n",
			expectedStatusCode:      401,
			serviceUpdatePageCalled: true,
			serviceUpdatePageParams: pageservice.UpdatePageParams{
				Page:   getPage("PG_1", 0, ""),
				UserID: "UR_1",
			},
			serviceUpdatePageReturnErr: &storeerror.NotAuthorized{
				UserID:  "UR_1",
				TableID: "PG_1",
				Err:     errors.New("failure"),
			},
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageService := new(mocks.PageService)
			pageService.On("UpdatePage", mock.Anything, tc.serviceUpdatePageParams).Return(tc.serviceUpdatePageReturnErr)
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlerutils.HandleTestRequest(handlerutils.HandleTestRequestParams{
				Method:         http.MethodPut,
				Endpoint:       fmt.Sprintf("page/%v", tc.pageID),
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "UpdatePage", testutils.GetExpectedNumberOfCalls(tc.serviceUpdatePageCalled))
		})
	}
}
