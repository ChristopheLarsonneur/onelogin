package api

import (
	"fmt"
	"net/http"

	"github.com/clarsonneur/onelogin/common"
)

// ResultStatus is the common result status
type ResultStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
}

// GetHeaders Compile headers for the API call
func GetHeaders(authorization string) common.Headers {
	return common.Headers{
		"Authorization": authorization,
		"Content-Type":  "application/json",
	}
}

func checkResponse(response *http.Response, inputErr error, status ResultStatus) (ret *http.Response, err error) {
	ret = response
	err = inputErr
	if err != nil {
		return
	}

	if status.Error {
		err = fmt.Errorf("%d %s: %s", response.StatusCode, status.Type, status.Message)
	}
	return
}
