package pagehandler

import (
	"encoding/json"
	"net/http"

	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// CreatePageRequest parameters from the CreatePage call
type CreatePageRequest struct {
	Title                string `json:"title"`
	Summary              string `json:"summary"`
	VersionID            string `json:"versionId"`
	PermissionTypeString string `json:"permission"`
	PermissionType       permission.Type
	PageTemplateID       string `json:"pageTemplateId"`
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
	if request.VersionID == "" {
		return request, errors.New("must provide versionId")
	}
	if request.PageTemplateID == "" {
		return request, errors.New("must provide pageTemplateId")
	}
	permissionType, err := permission.GetPermissionType(request.PermissionTypeString)
	if err != nil {
		return request, errors.New("permission is not a valid value")
	}
	request.PermissionType = permissionType
	return request, nil
}

// UpdatePageRequest parameters from the UpdatePage call
type UpdatePageRequest struct {
	GUID                 string
	Title                string `json:"title"`
	Summary              string `json:"summary"`
	VersionID            string `json:"versionId"`
	PermissionTypeString string `json:"permission"`
	PermissionType       permission.Type
	PageTemplateID       string `json:"pageTemplateId"`
}

// NewUpdatePageRequest extracts the UpdatePageRequest
func NewUpdatePageRequest(r *http.Request, p httprouter.Params) (UpdatePageRequest, error) {
	var request UpdatePageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, errors.New("invalid request")
	}
	request.GUID = p.ByName(PageIDRouteKey)
	return request.validate()
}

func (request UpdatePageRequest) validate() (UpdatePageRequest, error) {
	if request.GUID == "" {
		return request, errors.New("must provide a page id")
	}
	if request.PermissionTypeString != "" {
		permissionType, err := permission.GetPermissionType(request.PermissionTypeString)
		if err != nil {
			return request, errors.New("permission is not a valid value")
		}
		request.PermissionType = permissionType
	}
	return request, nil
}

// GetEntirePageRequest parameters from the GetEntirePage call
type GetEntirePageRequest struct {
	GUID string
}

// NewGetEntirePageRequest extracts the GetEntirePageRequest
func NewGetEntirePageRequest(r *http.Request, p httprouter.Params) (GetEntirePageRequest, error) {
	request, err := NewGetPageRequest(r, p)
	return GetEntirePageRequest{
		GUID: request.GUID,
	}, err
}

// GetPageRequest parameters from the GetPage call
type GetPageRequest struct {
	GUID string
}

// NewGetPageRequest extracts the GetPageRequest
func NewGetPageRequest(r *http.Request, p httprouter.Params) (GetPageRequest, error) {
	var request GetPageRequest
	request.GUID = p.ByName(PageIDRouteKey)
	return request.validate()
}

func (request GetPageRequest) validate() (GetPageRequest, error) {
	if request.GUID == "" {
		return request, errors.New("must provide a page id")
	}
	return request, nil
}

// DeletePageRequest parameters from the DeletePage call
type DeletePageRequest struct {
	GUID string
}

// NewDeletePageRequest extracts the DeletePageRequest
func NewDeletePageRequest(r *http.Request, p httprouter.Params) (DeletePageRequest, error) {
	var request DeletePageRequest
	request.GUID = p.ByName(PageIDRouteKey)
	return request.validate()
}

func (request DeletePageRequest) validate() (DeletePageRequest, error) {
	if request.GUID == "" {
		return request, errors.New("must provide a page id")
	}
	return request, nil
}

// GetPagesRequest parameters from the GetPages call
type GetPagesRequest struct {
	NextBatchID string
}

// NewGetPagesRequest extracts the GetPagesRequest
func NewGetPagesRequest(r *http.Request, p httprouter.Params) (GetPagesRequest, error) {
	var request GetPagesRequest
	request.NextBatchID = r.URL.Query().Get("nextBatchId")
	return request.validate()
}

func (request GetPagesRequest) validate() (GetPagesRequest, error) {
	return request, nil
}
