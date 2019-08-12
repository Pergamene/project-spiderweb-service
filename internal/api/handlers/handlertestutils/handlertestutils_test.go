// basic tests for requests that don't reach a particular handler.
package handlertestutils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
)

func TestAPI(t *testing.T) {
	cases := []struct {
		name                 string
		method               string
		endpoint             string
		headers              map[string]string
		authN                api.AuthN
		authZ                api.AuthZ
		datacenter           string
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name:                 "bad endpoint call",
			method:               http.MethodGet,
			endpoint:             "doesnotexist",
			authN:                DefaultAuthN("LOCAL"),
			authZ:                DefaultAuthZ(),
			expectedResponseBody: "{\"meta\":{\"httpStatus\":\"404 - Not Found\",\"message\":\"not found\"}}\n",
			expectedStatusCode:   404,
		},
		{
			name:                 "bad docs call",
			method:               http.MethodGet,
			endpoint:             "docs/doesnotexist",
			authN:                DefaultAuthN("LOCAL"),
			authZ:                DefaultAuthZ(),
			expectedResponseBody: "404 page not found\n",
			expectedStatusCode:   404,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, respBody := HandleTestRequest(HandleTestRequestParams{
				Method:         tc.method,
				Endpoint:       tc.endpoint,
				Headers:        tc.headers,
				Body:           nil,
				RouterHandlers: []api.RouterHandler{},
				AuthZ:          tc.authZ,
				AuthN:          tc.authN,
			})
			require.Equal(t, tc.expectedResponseBody, respBody)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
		})
	}
}
