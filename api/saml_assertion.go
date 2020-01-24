package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/clarsonneur/onelogin/common"
)

const (
	// SAMLAssertionURIPath API Path
	// As defined by https://developers.onelogin.com/api-docs/1/saml-assertions/generate-saml-assertion
	SAMLAssertionURIPath = "api/1/saml_assertion"
)

// SAMLAssertionResult match the result of the end point requested
type SAMLAssertionResult struct {
	Status ResultStatus
	Data   json.RawMessage `json:"data"`
	data   []SAMLAssertionDataResult
}

// SAMLAssertionRequest is the SAML assertion request structure
type SAMLAssertionRequest struct {
	User      string `json:"username_or_email"`
	Password  string `json:"password"`
	AppID     string `json:"app_id"`
	SubDomain string `json:"subdomain"`
	IPAddress string `json:"ip_address,omitempty"`
}

// NewSAMLAssertionResult creates tge SAMl Assertion result
func NewSAMLAssertionResult() (ret *SAMLAssertionResult) {
	ret = new(SAMLAssertionResult)
	return
}

// Post the SAMLAssertion request and saved it to the SAMLAssertionResult
func (r *SAMLAssertionResult) Post(a *Core, user, pass, appID, subDomain, IP string) (response *http.Response, err error) {
	if r == nil {
		return nil, errors.New("SAMLAssertionResult is nil")
	}

	input := SAMLAssertionRequest{
		User:      user,
		Password:  pass,
		AppID:     appID,
		SubDomain: subDomain,
		IPAddress: IP,
	}

	response, err = common.Request("POST", a.getBearerHeaders(), a.GetURL(SAMLAssertionURIPath), input, r)
	return checkResponse(response, err, r.Status)
}
