package healthcheckhandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/healthcheck/mocks"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/page"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/mock"
)

// var testHandler api.Handler
// var ctx context.Context

func TestMain(m *testing.M) {
	result := m.Run()
	os.Exit(result)
}

type headerSet struct {
	key   string
	value string
}

func TestIsHealthy(t *testing.T) {
	cases := []struct {
		name                            string
		headers                         []headerSet
		authN                           api.AuthN
		authZ                           api.AuthZ
		datacenter                      string
		serviceIsHealthyCalled          bool
		serviceIsHealthyReturnIsHealthy bool
		serviceIsHealthyReturnErr       error
		expectedResponseBody            string
		expectedStatusCode              int
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
			serviceIsHealthyCalled:          true,
			serviceIsHealthyReturnIsHealthy: true,
			expectedResponseBody:            "{\"result\":{\"status\":\"ok\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:              200,
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
			serviceIsHealthyCalled:          true,
			serviceIsHealthyReturnIsHealthy: false,
			expectedResponseBody:            "{\"result\":{\"status\":\"error\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:              200,
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
			headers: []headerSet{
				{
					key:   api.AdminAuthSecretHeaderKey,
					value: "SECRET",
				},
			},
			authN: api.AuthN{
				Datacenter:      "PROD",
				AdminAuthSecret: "SECRET",
			},
			authZ: api.AuthZ{
				APIPath: "api/test",
			},
			serviceIsHealthyCalled:          true,
			serviceIsHealthyReturnIsHealthy: true,
			expectedResponseBody:            "{\"result\":{\"status\":\"ok\"},\"meta\":{\"httpStatus\":\"200 - OK\"}}\n",
			expectedStatusCode:              200,
		},
		{
			name: "standard healthy healthcheck, prod, but bad secret",
			headers: []headerSet{
				{
					key:   api.AdminAuthSecretHeaderKey,
					value: "BAD_SECRET",
				},
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
			healthcheckHandler := HealthcheckHandler{
				HealthcheckService: healthcheckService,
			}
			routerHandlers := api.RouterHandlers{
				HealthcheckHandler: healthcheckHandler,
				PageHandler:        pagehandler.PageHandler{},
			}
			router := api.NewRouter(tc.authZ.APIPath, "static/test", routerHandlers)
			testHandler := api.Handler{
				AuthN:      tc.authN,
				AuthZ:      tc.authZ,
				Router:     router,
				Datacenter: tc.authN.Datacenter,
				APIPath:    tc.authZ.APIPath,
			}
			url := fmt.Sprintf("http://test.com/%v/healthcheck", tc.authZ.APIPath)
			r := httptest.NewRequest(http.MethodGet, url, nil)
			for _, header := range tc.headers {
				r.Header.Set(header.key, header.value)
			}
			w := httptest.NewRecorder()
			testHandler.ServeHTTP(w, r)
			healthcheckService.AssertNumberOfCalls(t, "IsHealthy", testutils.GetExpectedNumberOfCalls(tc.serviceIsHealthyCalled))

			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			require.Equal(t, tc.expectedResponseBody, string(body))
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
		})
	}
}
