package onelogin

import "github.com/clarsonneur/onelogin/common"

// GetHeaders Compile headers for the API call
func GetHeaders(authorization string) common.Headers {
	return common.Headers{
		"Authorization": authorization,
		"Content-Type":  "application/json",
	}
}
