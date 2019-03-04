// Package handlertestutils are tools for testing handlers
package handlertestutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
)

// HandleTestRequestParams are the params for the HandleTestRequest function.
type HandleTestRequestParams struct {
	Method         string
	Endpoint       string
	Params         url.Values
	Headers        map[string]string
	Body           io.Reader
	RouterHandlers []api.RouterHandler
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
	uri := fmt.Sprintf("http://test.com/%v/%v", p.AuthZ.APIPath, p.Endpoint)
	params := p.Params.Encode()
	if params != "" {
		uri = uri + "?" + params
	}
	r := httptest.NewRequest(p.Method, uri, p.Body)
	for key, value := range p.Headers {
		r.Header.Set(key, value)
	}
	w := httptest.NewRecorder()
	testHandler.ServeHTTP(w, r)
	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return resp, string(respBody)
}

// DefaultAuthN is a quick way to pass in the AuthN struct to a test
func DefaultAuthN(datacenter string) api.AuthN {
	return api.AuthN{
		Datacenter:      datacenter,
		AdminAuthSecret: "SECRET",
	}
}

// DefaultAuthZ is a quick way to pass in the AuthZ struct to a test
func DefaultAuthZ() api.AuthZ {
	return api.AuthZ{
		APIPath: "api/test",
	}
}
