package api

import (
	"errors"
	"net/http"

	"github.com/clarsonneur/onelogin/common"
)

// https://developers.onelogin.com/api-docs/1/roles/get-roles

const (
	// GetRolesURIPath defined the API Path for such request.
	// As defined by https://developers.onelogin.com/api-docs/1/roles/get-roles
	GetRolesURIPath = "api/1/roles"
)

// GetRolesResult match the result of the end point requested
type GetRolesResult struct {
	Status ResultStatus
	Data   Roles `json:"data"`
}

// GetRolesRequest is the input request structure for this API call.
type GetRolesRequest struct {
}

// NewGetRoles return a new object VerifyFactorResult
func NewGetRoles() (ret *GetRolesResult) {
	ret = new(GetRolesResult)
	return
}

// Get the request as defined by the API
func (r *GetRolesResult) Get(a *Core) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("GetRolesResult is nil")
	}

	input := GetRolesResult{}

	response, err = common.Request("GET", a.getBearerHeaders(), a.GetURL(GetRolesURIPath), input, r)
	return checkResponse(response, err, r.Status)
}
