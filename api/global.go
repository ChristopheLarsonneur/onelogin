package api

import "github.com/clarsonneur/onelogin/common"

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
