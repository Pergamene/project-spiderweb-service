package healthcheckhandler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	handlerutils "github.com/Pergamene/project-spiderweb-service/internal/api/handlers"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/healthcheck/mocks"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/mock"
)

func TestIsHealthy(t *testing.T) {
	cases := []struct {
		name                            string
		headers                         map[string]string
		authN                           api.AuthN
		authZ                           api.AuthZ
		datacenter                      string
		expectedResponseBody            string
		expectedStatusCode              int
		serviceIsHealthyCalled          bool
		serviceIsHealthyReturnIsHealthy bool
		serviceIsHealthyReturnErr       error
	}{
		{
			name: "standard healthy healthcheck, local",
			authN: api.AuthN{
				Datacenter:      "LOCAL",
				AdminAuthSecret: "SECRET",
			},
			authZ: api.AuthZ{
				APIPath: "api/test",
			},
			expectedResponseBody:            "{\"result\":{\"status\":\"ok\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:              200,
			serviceIsHealthyCalled:          true,
			serviceIsHealthyReturnIsHealthy: true,
		},
		{
			name: "bad healthcheck, local",
			authN: api.AuthN{
				Datacenter:      "LOCAL",
				AdminAuthSecret: "SECRET",
			},
			authZ: api.AuthZ{
				APIPath: "api/test",
			},
			expectedResponseBody:            "{\"result\":{\"status\":\"error\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:              200,
			serviceIsHealthyCalled:          true,
			serviceIsHealthyReturnIsHealthy: false,
		},
		{
			name: "not authenticated",
			authN: api.AuthN{
				Datacenter:      "PROD",
				AdminAuthSecret: "SECRET",
			},
			authZ: api.AuthZ{
				APIPath: "api/test",
			},
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
		{
			name: "standard healthy healthcheck, prod",
			headers: map[string]string{
				"X-ADMIN-AUTH-SECRET": "SECRET",
			},
			authN: api.AuthN{
				Datacenter:      "PROD",
				AdminAuthSecret: "SECRET",
			},
			authZ: api.AuthZ{
				APIPath: "api/test",
			},
			expectedResponseBody:            "{\"result\":{\"status\":\"ok\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:              200,
			serviceIsHealthyCalled:          true,
			serviceIsHealthyReturnIsHealthy: true,
		},
		{
			name: "standard healthy healthcheck, prod, but bad secret",
			headers: map[string]string{
				"X-ADMIN-AUTH-SECRET": "BAD_SECRET",
			},
			authN: api.AuthN{
				Datacenter:      "PROD",
				AdminAuthSecret: "SECRET",
			},
			authZ: api.AuthZ{
				APIPath: "api/test",
			},
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"401 - Unauthorized\",\"message\":\"not authenticated\"}}\n",
			expectedStatusCode:   401,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			healthcheckService := new(mocks.HealthcheckService)
			healthcheckService.On("IsHealthy", mock.Anything).Return(tc.serviceIsHealthyReturnIsHealthy, tc.serviceIsHealthyReturnErr)
			routerHandlers := HealthcheckRouterHandlers(tc.authZ.APIPath, healthcheckService)
			resp, respBody := handlerutils.HandleTestRequest(handlerutils.HandleTestRequestParams{
				Method:         http.MethodGet,
				Endpoint:       "healthcheck",
				Headers:        tc.headers,
				Body:           nil,
				RouterHandlers: routerHandlers,
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			healthcheckService.AssertNumberOfCalls(t, "IsHealthy", testutils.GetExpectedNumberOfCalls(tc.serviceIsHealthyCalled))
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
		})
	}
}
