package api

import (
	"errors"
	"github.com/clarsonneur/onelogin/common"
	"net/http"
	"strconv"
)

// https://developers.onelogin.com/api-docs/1/saml-assertions/verify-factor

const (
	// VerifyFactorURIPath defined the API Path for such request.
	// As defined by https://developers.onelogin.com/api-docs/1/saml-assertions/verify-factor
	VerifyFactorURIPath = "api/1/saml_assertion/verify_factor"
)

// VerifyFactorResult match the result of the end point requested
type VerifyFactorResult struct {
	Status ResultStatus
	Data   string `json:"data"`
}

// VerifyFactorRequest is the input request structure for this API call.
type VerifyFactorRequest struct {
	AppID       string `json:"app_id"`
	DeviceID    string `json:"device_id"`
	StateToken  string `json:"state_token"`
	OTPToken    string `json:"otp_token"`
	DoNotNotify bool   `json:"do_not_notify"`
}

// NewVerifyFactorResult return a new object VerifyFactorResult
func NewVerifyFactorResult() (ret *VerifyFactorResult) {
	ret = new(VerifyFactorResult)
	return
}

// Post the request as defined by the API
func (r *VerifyFactorResult) Post(a *Core, appID string, deviceID int, stateToken, OTPToken string, doNotNotify bool) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("VerifyFactorResult is nil")
	}

	input := VerifyFactorRequest{
		AppID:       appID,
		DeviceID:    strconv.Itoa(deviceID),
		StateToken:  stateToken,
		OTPToken:    OTPToken,
		DoNotNotify: doNotNotify,
	}

	response, err = common.Request("POST", a.getBearerHeaders(), a.GetURL(VerifyFactorURIPath), input, r)
	return
}
