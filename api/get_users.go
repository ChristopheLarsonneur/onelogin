package api

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/clarsonneur/onelogin/common"
)

// https://developers.onelogin.com/api-docs/1/users/get-users

const (
	// GetUsersURIPath defined the API Path for such request.
	// As defined by https://developers.onelogin.com/api-docs/1/users/get-users
	GetUsersURIPath = "api/1/users"
)

// GetUsersResult match the result of the end point requested
type GetUsersResult struct {
	Status     ResultStatus
	Pagination ResultPagination
	Data       Users `json:"data"`
	url        *url.URL
}

// GetUsersRequest is the input request structure for this API call.
type GetUsersRequest struct {
}

// NewGetUsers return a new object GetUsersResult
func NewGetUsers() (ret *GetUsersResult) {
	ret = new(GetUsersResult)
	return
}

// Get the request as defined by the API
func (r *GetUsersResult) Get(a *Core) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("GetUsersResult is nil")
	}

	r.url, err = url.Parse(a.GetURL(GetUsersURIPath))
	response, err = common.Request("GET", a.getBearerHeaders(), r.url.String(), nil, r)
	return
}

// Next return the next pagination result
// if response and err is nil, then there is no more next page to get.
func (r *GetUsersResult) Next(a *Core) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("GetUsersResult is nil")
	}

	if q := r.url.Query(); q.Get("after_cursor") == r.Pagination.AfterCursor {
		return
	}

	common.SetQuery(r.url, map[string]string{
		"after_cursor": r.Pagination.AfterCursor},
	)

	response, err = common.Request("GET", a.getBearerHeaders(), r.url.String(), nil, r)
	return
}
