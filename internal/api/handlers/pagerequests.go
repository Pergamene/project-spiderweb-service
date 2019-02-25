package handlers

import (
	"encoding/json"
	"net/http"

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
