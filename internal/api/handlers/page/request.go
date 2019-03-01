package pagehandler

import (
	"encoding/json"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/api"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// CreatePageRequest parameters from the CreatePage call
type CreatePageRequest struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// NewCreatePageRequest extracts the CreatePageRequest
func NewCreatePageRequest(r *http.Request, p httprouter.Params) (CreatePageRequest, error) {
	var request CreatePageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, errors.New("invalid request")
	}
	return request.validate()
}

func (request CreatePageRequest) validate() (CreatePageRequest, error) {
	if request.Title == "" {
		return request, errors.New("must provide title")
	}
	return request, nil
}

// UpdatePageRequest parameters from the UpdatePage call
type UpdatePageRequest struct {
	GUID    string
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// NewUpdatePageRequest extracts the UpdatePageRequest
func NewUpdatePageRequest(r *http.Request, p httprouter.Params) (UpdatePageRequest, error) {
	var request UpdatePageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, errors.New("invalid request")
	}
	request.GUID = p.ByName(api.PageIDRouteKey)
	return request.validate()
}

func (request UpdatePageRequest) validate() (UpdatePageRequest, error) {
	if request.GUID == "" {
		return request, errors.New("must provide a page id")
	}
	if request.Title == "" {
		return request, errors.New("a page must retain a title")
	}
	return request, nil
}
