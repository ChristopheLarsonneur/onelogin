package api

import (
	"fmt"

	"github.com/clarsonneur/onelogin/common"
)

// OneLoginURL is the default OneLogin API endpoint
const (
	OneLoginURL                = "https://api.%s.onelogin.com"
	TimeSleepOnResponsePending = 15
	MaxIterGetSAMLResponse     = 6
)

// Core is the core API object
type Core struct {
	// Allow custom base URL to override the generated URL
	CustomURL string

	// OneLogin service shard (eu, us, etc)
	Shard string

	SubDomain    string
	ClientID     string
	ClientSecret string

	// Token struct for managing the OAuth token
	Token *OAuthTokenResult
}

// NewAPI create the main API object
func NewAPI(shard string, clientID string, clientSecret string, subdomain string) *Core {
	ol := Core{Shard: shard, ClientID: clientID, ClientSecret: clientSecret, SubDomain: subdomain}
	return &ol
}

// GetURL creates a URL given the URI and any given args.
// Returns a URL
// You can complete the URL with the query part, thanks to common.
func (o *Core) GetURL(uri string, args ...interface{}) string {
	// Handle cases where the uri requires variable replacements (ie.  /api/1/user/%d/roles)
	fulluri := uri

	if len(args) > 0 {
		fulluri = fmt.Sprintf(uri, args...)
	}

	if o.CustomURL != "" {
		return fmt.Sprintf("%s/%s", o.CustomURL, fulluri)
	}

	return fmt.Sprintf("%s/%s", fmt.Sprintf(OneLoginURL, o.Shard), fulluri)
}

// ObtainAPIAccess initialize the access to the API.
func (o *Core) ObtainAPIAccess() (err error) {
	o.Token = NewOAuthTokenResult()

	return o.Token.Obtain(o)
}

func (o *Core) getBearerHeaders() (ret common.Headers) {
	ret = GetHeaders("bearer:" + o.Token.AccessToken)
	return
}
