package onelogin

import (
	"encoding/base64"
	"fmt"
	"github.com/clarsonneur/onelogin/common"
	"time"
)

// https://developers.onelogin.com/api-docs/1/oauth20-tokens/generate-tokens-2
const (
	TokenURIPath = "auth/oauth2/v2/token"
)

// OAuthTokenResult map auth/oauth2/v2/token result response
type OAuthTokenResult struct {
	ResultStatus
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	AccountID   int    `json:"account_id"`
}

// TokenRequest map auth/oauth2/v2/token request
type TokenRequest struct {
	GrantType string `json:"grant_type"`
}

// NewOAuthTokenResult return an OAth token result object
func NewOAuthTokenResult() (ret *OAuthTokenResult) {
	ret = new(OAuthTokenResult)
	return
}

// Obtain get or refresh a token to manage the OneLogin API
func (t *OAuthTokenResult) Obtain(a *API) (err error) {
	if t.isExpired() {
		return t.getToken(a)
	}
	// No need to refresh it
	return nil
}

func (t *OAuthTokenResult) getToken(a *API) (err error) {
	url := a.GetURL(TokenURIPath)

	input := TokenRequest{
		GrantType: "client_credentials",
	}

	// base64 encode "clientId:clientSecret"
	authorization := base64.StdEncoding.EncodeToString([]byte(a.ClientID + ":" + a.ClientSecret))
	headers := GetHeaders("Basic " + authorization)

	_, err = common.Request("POST", headers, url, input, t)

	if t.ResultStatus.Error {
		err = fmt.Errorf("APIToken error: %s", t.ResultStatus.Message)
	}

	if t.AccessToken == "" {
		err = fmt.Errorf("APIToken error: Enable to obtain Access Token thanks to API keys")
	}
	return
}

// isExpired return false if the current token is valid. false otherwise.
// An empty token is considered as invalid and will return expired = true.
func (t *OAuthTokenResult) isExpired() (expired bool) {
	expired = true
	if t.AccessToken == "" {
		return
	}
	expiredDate, err := time.Parse(time.RFC3339, t.CreatedAt)
	if err != nil {
		fmt.Printf("Error to interpret RFC3339 created date '%s'. Considering API token as expired\n", t.CreatedAt)
		return
	}
	expiredDate.Add(time.Second * time.Duration(t.ExpiresIn))
	return time.Now().After(expiredDate)
}
