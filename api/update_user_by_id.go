package api

import (
	"errors"
	"net/http"

	"github.com/clarsonneur/onelogin/common"
)

// https://developers.onelogin.com/api-docs/1/users/update-user

const (
	// UpdateUserByIDURIPath defined the API Path for such request.
	// As defined by https://developers.onelogin.com/api-docs/1/users/update-user
	UpdateUserByIDURIPath = "api/1/users/%d"
)

// PutUserByIDResult match the result of the end point requested
type PutUserByIDResult struct {
	Status ResultStatus
	Data   Users `json:"data"`
}

// PutUserRequest is the input request structure for this API call.
type PutUserRequest struct {
	CustomAttrs map[string]string `json:"custom_attributes,omitempty"`
}

// NewPutUserByID return a new object PutUserByIDResult
func NewPutUserByID() (ret *PutUserByIDResult) {
	ret = new(PutUserByIDResult)
	return
}

// Put the request as defined by the API
func (r *PutUserByIDResult) Put(a *Core, id int64, input PutUserRequest) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("PutUserByIDResult is nil")
	}

	response, err = common.Request("PUT", a.getBearerHeaders(), a.GetURL(UpdateUserByIDURIPath, id), input, r)
	return
}
