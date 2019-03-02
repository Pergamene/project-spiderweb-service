package handlerutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/healthcheck"
	"github.com/Pergamene/project-spiderweb-service/internal/api/handlers/page"
)

// HandleTestRequestParams are the params for the HandleTestRequest function.
type HandleTestRequestParams struct {
	Method         string
	Endpoint       string
	Headers        map[string]string
	Body           *io.Reader
	RouterHandlers api.RouterHandlers
	AuthZ          api.AuthZ
	AuthN          api.AuthN
}

// HandleTestRequest handles making the request for a given test and returning the response and response body.
func HandleTestRequest(p HandleTestRequestParams) (*http.Response, string) {
	router := api.NewRouter(p.AuthZ.APIPath, "static/test", p.RouterHandlers)
	testHandler := api.Handler{
		AuthN:      p.AuthN,
		AuthZ:      p.AuthZ,
		Router:     router,
		Datacenter: p.AuthN.Datacenter,
		APIPath:    p.AuthZ.APIPath,
	}
	url := fmt.Sprintf("http://test.com/%v/%v", p.AuthZ.APIPath, p.Endpoint)
	r := httptest.NewRequest(p.Method, url, nil)
	for key, value := range p.Headers {
		r.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return resp, string(respBody)
}

// GetBaseRouterHandlers returns a RouterHandlers with default handler instantiation.
// This is to ensure no nil pointers are referenced.
func GetBaseRouterHandlers() api.RouterHandlers {
	return api.RouterHandlers{
		HealthcheckHandler: healthcheckhandler.HealthcheckHandler{},
		PageHandler:        pagehandler.PageHandler{},
	}
}
