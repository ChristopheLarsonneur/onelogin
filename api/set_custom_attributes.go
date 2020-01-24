package api

import (
	"errors"
	"net/http"

	"github.com/clarsonneur/onelogin/common"
)

// https://developers.onelogin.com/api-docs/1/users/set-custom-attribute

const (
	// SetCustomAttrsIDURIPath defined the API Path for such request.
	// As defined by https://developers.onelogin.com/api-docs/1/users/set-custom-attribute
	SetCustomAttrsIDURIPath = "api/1/users/%d/set_custom_attributes"
)

// PutUserAttrsResult match the result of the end point requested
type PutUserAttrsResult struct {
	Status ResultStatus
}

// PutUserAttrsRequest is the input request structure for this API call.
type PutUserAttrsRequest struct {
	CustomAttrs map[string]string `json:"custom_attributes,omitempty"`
}

// NewPutCustomAttrs return a new object PutUserByIDResult
func NewPutCustomAttrs() (ret *PutUserAttrsResult) {
	ret = new(PutUserAttrsResult)
	return
}

// Put the request as defined by the API
func (r *PutUserAttrsResult) Put(a *Core, id int64, input PutUserAttrsRequest) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("PutUserAttrsResult is nil")
	}

	response, err = common.Request("PUT", a.getBearerHeaders(), a.GetURL(SetCustomAttrsIDURIPath, id), input, r)
	return checkResponse(response, err, r.Status)
}
