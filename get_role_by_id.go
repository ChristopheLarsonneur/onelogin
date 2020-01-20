package onelogin

import (
	"errors"
	"net/http"

	"github.com/clarsonneur/onelogin/common"
)

// https://developers.onelogin.com/api-docs/1/roles/get-role-by-id

const (
	// GetRoleByIDURIPath defined the API Path for such request.
	// As defined by https://developers.onelogin.com/api-docs/1/roles/get-role-by-id
	GetRoleByIDURIPath = "api/1/roles/%d"
)

// GetRoleByIDResult match the result of the end point requested
type GetRoleByIDResult struct {
	Status ResultStatus
	Data   []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

// GetRoleByIDRequest is the input request structure for this API call.
type GetRoleByIDRequest struct {
}

// NewGetUserByID return a new object VerifyFactorResult
func NewGetRoleByID() (ret *GetRoleByIDResult) {
	ret = new(GetRoleByIDResult)
	return
}

// Get the request as defined by the API
func (r *GetRoleByIDResult) Get(a *API, id int64) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("GetRoleByIDResult is nil")
	}

	input := GetRoleByIDResult{}

	response, err = common.Request("GET", a.getBearerHeaders(), a.GetURL(common.BuildURI(GetUserByIDURIPath, id)), input, r)
	return
}
