package onelogin

import (
	"errors"
	"net/http"

	"github.com/clarsonneur/onelogin/common"
)

// https://developers.onelogin.com/api-docs/1/users/get-user-by-id

const (
	// GetUserByIDURIPath defined the API Path for such request.
	// As defined by https://developers.onelogin.com/api-docs/1/users/get-user-by-id
	GetUserByIDURIPath = "api/1/users/%d"
)

// GetUserByIDResult match the result of the end point requested
type GetUserByIDResult struct {
	Status ResultStatus
	Data   struct {
		Email         string            `json:"email"`
		ID            int16             `json:"id"`
		Status        int               `json:"status"`
		State         int               `json:"state"`
		RolesID       []int16           `json:"role_id"`
		ManagerUserID int16             `json:"manager_user_id"`
		MemberOf      string            `json:"member_of"`
		CustomAttrs   map[string]string `json:"custom_attributes"`
	} `json:"data"`
}

// GetUserByIDRequest is the input request structure for this API call.
type GetUserByIDRequest struct {
}

// NewGetUserByIDResult return a new object VerifyFactorResult
func NewGetUserByIDResult() (ret *GetUserByIDResult) {
	ret = new(GetUserByIDResult)
	return
}

// Get the request as defined by the API
func (r *GetUserByIDResult) Get(a *API, id int16) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("GetUserByIDResult is nil")
	}

	input := GetUserByIDResult{}

	response, err = common.Request("GET", a.getBearerHeaders(), a.GetURL(common.BuildURI(GetUserByIDURIPath, id)), input, r)
	return
}
