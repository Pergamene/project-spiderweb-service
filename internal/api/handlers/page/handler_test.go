package pagehandler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/pagetemplate"

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

func getPage(guid, title, summary, versionID, pageTemplateID string, permissionType permission.Type) page.Page {
	return page.Page{
		GUID:    guid,
		Title:   title,
		Summary: summary,
		Version: version.Version{
			GUID: versionID,
		},
		PermissionType: permissionType,
		PageTemplate: pagetemplate.PageTemplate{
			GUID: pageTemplateID,
		},
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
			name: "happy page, local",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody: "{\"title\":\"test title\",\"summary\":\"test summary\",\"versionId\":\"VR_1\",\"pageTemplateId\":\"PGT_1\",\"permission\":\"PR\"}",
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
						Page:    getPage("", "test title", "test summary", "VR_1", "PGT_1", permission.TypePrivate),
						OwnerID: "UR_1",
					},
					returnRecord: getPage("PG_1", "test title", "test summary", "VR_1", "PGT_1", permission.TypePrivate),
				},
			},
		},
		{
			name: "missing title for the page",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:          "{\"summary\":\"test summary\",\"versionId\":\"VR_1\",\"pageTemplateId\":\"PGT_1\",\"permission\":\"PR\"}",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"400 - Bad Request\",\"message\":\"must provide title\"}}\n",
			expectedStatusCode:   400,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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

type setPageCall struct {
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
		setPageCalls         []setPageCall
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
			name:   "happy page, local",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:          "{\"title\":\"test title\",\"summary\":\"test summary\"}",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			setPageCalls: []setPageCall{
				{
					pageParams: pageservice.SetPageParams{
						Page:   getPage("PG_1", "test title", "test summary", "", "", ""),
						UserID: "UR_1",
					},
				},
			},
		},
		{
			name:   "trying to edit a page that you don't have permission to edit",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			requestBody:          "{\"title\":\"test title\",\"summary\":\"test summary\"}",
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authorized\"}}\n",
			expectedStatusCode:   401,
			setPageCalls: []setPageCall{
				{
					pageParams: pageservice.SetPageParams{
						Page:   getPage("PG_1", "test title", "test summary", "", "", ""),
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
		t.Run(tc.name, func(t *testing.T) {
			pageService := new(mocks.PageService)
			for index := range tc.setPageCalls {
				pageService.On("SetPage", mock.Anything, tc.setPageCalls[index].pageParams).Return(tc.setPageCalls[index].returnErr)
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
			pageService.AssertNumberOfCalls(t, "SetPage", len(tc.setPageCalls))
		})
	}
}

type getEntirePageCall struct {
	pageParams pageservice.GetEntirePageParams
	returnPage page.Page
	returnErr  error
}

func TestGetEntirePage(t *testing.T) {
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
		getEntirePageCalls   []getEntirePageCall
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
			name:   "happy page, local",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"version\":{\"id\":\"VR_1\",\"name\":\"\",\"parentId\":\"\"},\"pageTemplate\":{\"name\":\"\",\"guid\":\"PGT_1\"},\"id\":\"PG_1\",\"title\":\"test title\",\"summary\":\"test summary\",\"permission\":\"PR\",\"properties\":[],\"details\":[],\"createdAt\":null,\"updatedAt\":null},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			getEntirePageCalls: []getEntirePageCall{
				{
					pageParams: pageservice.GetEntirePageParams{
						Page:   getPage("PG_1", "", "", "", "", ""),
						UserID: "UR_1",
					},
					returnPage: getPage("PG_1", "test title", "test summary", "VR_1", "PGT_1", permission.TypePrivate),
				},
			},
		},
		{
			name:   "trying to get a page that you don't have permission to read",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authorized\"}}\n",
			expectedStatusCode:   401,
			getEntirePageCalls: []getEntirePageCall{
				{
					pageParams: pageservice.GetEntirePageParams{
						Page:   getPage("PG_1", "", "", "", "", ""),
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
		t.Run(tc.name, func(t *testing.T) {
			pageService := new(mocks.PageService)
			for index := range tc.getEntirePageCalls {
				pageService.On("GetEntirePage", mock.Anything, tc.getEntirePageCalls[index].pageParams).Return(tc.getEntirePageCalls[index].returnPage, tc.getEntirePageCalls[index].returnErr)
			}
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         http.MethodGet,
				Endpoint:       fmt.Sprintf("pages/%v/full", tc.pageID),
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "GetEntirePage", len(tc.getEntirePageCalls))
		})
	}
}

type getPageCall struct {
	pageParams pageservice.GetPageParams
	returnPage page.Page
	returnErr  error
}

func TestGetPage(t *testing.T) {
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
		getPageCalls         []getPageCall
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
			name:   "happy page, local",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"versionId\":\"VR_1\",\"pageTemplateId\":\"PGT_1\",\"id\":\"PG_1\",\"title\":\"test title\",\"summary\":\"test summary\",\"permission\":\"PR\",\"createdAt\":null,\"updatedAt\":null},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			getPageCalls: []getPageCall{
				{
					pageParams: pageservice.GetPageParams{
						Page:   getPage("PG_1", "", "", "", "", ""),
						UserID: "UR_1",
					},
					returnPage: getPage("PG_1", "test title", "test summary", "VR_1", "PGT_1", permission.TypePrivate),
				},
			},
		},
		{
			name:   "trying to get a page that you don't have permission to read",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authorized\"}}\n",
			expectedStatusCode:   401,
			getPageCalls: []getPageCall{
				{
					pageParams: pageservice.GetPageParams{
						Page:   getPage("PG_1", "", "", "", "", ""),
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
		t.Run(tc.name, func(t *testing.T) {
			pageService := new(mocks.PageService)
			for index := range tc.getPageCalls {
				pageService.On("GetPage", mock.Anything, tc.getPageCalls[index].pageParams).Return(tc.getPageCalls[index].returnPage, tc.getPageCalls[index].returnErr)
			}
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         http.MethodGet,
				Endpoint:       fmt.Sprintf("pages/%v", tc.pageID),
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "GetPage", len(tc.getPageCalls))
		})
	}
}

type getPagesCall struct {
	pageParams        pageservice.GetPagesParams
	returnPages       []page.Page
	returnTotal       int
	returnNextBatchID string
	returnErr         error
}

func TestGetPages(t *testing.T) {
	cases := []struct {
		name                 string
		headers              map[string]string
		requestBody          string
		authN                api.AuthN
		authZ                api.AuthZ
		datacenter           string
		expectedResponseBody string
		expectedStatusCode   int
		getPagesCalls        []getPagesCall
	}{
		{
			name:                 "not authenticated",
			authN:                handlertestutils.DefaultAuthN("PROD"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
		{
			name: "happy page, local",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"batch\":[{\"versionId\":\"VR_1\",\"pageTemplateId\":\"PGT_1\",\"id\":\"PG_1\",\"title\":\"test title\",\"summary\":\"test summary\",\"permission\":\"PR\",\"createdAt\":null,\"updatedAt\":null},{\"versionId\":\"VR_1\",\"pageTemplateId\":\"PGT_1\",\"id\":\"PG_2\",\"title\":\"test title 2 \",\"summary\":\"test summary 2\",\"permission\":\"PR\",\"createdAt\":null,\"updatedAt\":null},{\"versionId\":\"VR_1\",\"pageTemplateId\":\"PGT_2\",\"id\":\"PG_3\",\"title\":\"test title 3\",\"summary\":\"test summary 3\",\"permission\":\"PU\",\"createdAt\":null,\"updatedAt\":null}],\"total\":10,\"nextBatch\":{\"paramKey\":\"nextBatchId\",\"paramValue\":\"PG_4\"}},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			getPagesCalls: []getPagesCall{
				{
					pageParams: pageservice.GetPagesParams{
						NextBatchID: "",
						UserID:      "UR_1",
					},
					returnPages: []page.Page{
						getPage("PG_1", "test title", "test summary", "VR_1", "PGT_1", permission.TypePrivate),
						getPage("PG_2", "test title 2 ", "test summary 2", "VR_1", "PGT_1", permission.TypePrivate),
						getPage("PG_3", "test title 3", "test summary 3", "VR_1", "PGT_2", permission.TypePublic),
					},
					returnTotal:       10,
					returnNextBatchID: "PG_4",
				},
			},
		},
		{
			name: "returning no pages",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"result\":{\"batch\":[],\"total\":0},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			getPagesCalls: []getPagesCall{
				{
					pageParams: pageservice.GetPagesParams{
						NextBatchID: "",
						UserID:      "UR_1",
					},
					returnPages:       []page.Page{},
					returnTotal:       0,
					returnNextBatchID: "",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pageService := new(mocks.PageService)
			for index := range tc.getPagesCalls {
				pageService.On("GetPages", mock.Anything, tc.getPagesCalls[index].pageParams).Return(tc.getPagesCalls[index].returnPages, tc.getPagesCalls[index].returnTotal, tc.getPagesCalls[index].returnNextBatchID, tc.getPagesCalls[index].returnErr)
			}
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         http.MethodGet,
				Endpoint:       "pages",
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "GetPages", len(tc.getPagesCalls))
		})
	}
}

type removePageCall struct {
	pageParams pageservice.RemovePageParams
	returnErr  error
}

func TestDeletePage(t *testing.T) {
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
		removePageCalls      []removePageCall
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
			name:   "happy page, local",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:   200,
			removePageCalls: []removePageCall{
				{
					pageParams: pageservice.RemovePageParams{
						Page: page.Page{
							GUID: "PG_1",
						},
						UserID: "UR_1",
					},
				},
			},
		},
		{
			name:   "trying to remove a page that you don't have permission to edit",
			pageID: "PG_1",
			headers: map[string]string{
				"X-USER-ID": "UR_1",
			},
			authN:                handlertestutils.DefaultAuthN("LOCAL"),
			authZ:                handlertestutils.DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authorized\"}}\n",
			expectedStatusCode:   401,
			removePageCalls: []removePageCall{
				{
					pageParams: pageservice.RemovePageParams{
						Page: page.Page{
							GUID: "PG_1",
						},
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
		t.Run(tc.name, func(t *testing.T) {
			pageService := new(mocks.PageService)
			for index := range tc.removePageCalls {
				pageService.On("RemovePage", mock.Anything, tc.removePageCalls[index].pageParams).Return(tc.removePageCalls[index].returnErr)
			}
			routerHandlers := PageRouterHandlers(tc.authZ.APIPath, pageService)
			resp, respBody := handlertestutils.HandleTestRequest(handlertestutils.HandleTestRequestParams{
				Method:         http.MethodDelete,
				Endpoint:       fmt.Sprintf("pages/%v", tc.pageID),
				Headers:        tc.headers,
				Body:           strings.NewReader(tc.requestBody),
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			pageService.AssertNumberOfCalls(t, "RemovePage", len(tc.removePageCalls))
		})
	}
}
