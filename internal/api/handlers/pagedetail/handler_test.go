package pagedetailhandler

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	pagedetailservice "github.com/Pergamene/project-spiderweb-service/internal/services/pagedetail"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/handlertestutils"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/page/mocks"
)

type updatePageDetailCall struct {
	pageDetailParams pagedetailservice.UpdatePageDetailParams
	returnErr        error
}

func TestUpdatePageDetail(t *testing.T) {
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
					pageParams: pageservice.UpdatePageParams{
						Page:   getPage("PG_1", 0, ""),
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
					pageParams: pageservice.UpdatePageParams{
						Page:   getPage("PG_1", 0, ""),
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
