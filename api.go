package onelogin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/clarsonneur/onelogin/common"
	"github.com/op/go-logging"
)

// OneLoginURL is the default OneLogin API endpoint
const (
	OneLoginURL                = "https://api.%s.onelogin.com"
	TimeSleepOnResponsePending = 15
	MaxIterGetSAMLResponse     = 6
)

// API is the core API object
type API struct {
	// Allow custom base URL to override the generated URL
	CustomURL string

	// OneLogin service shard (eu, us, etc)
	Shard string

	// Token struct for managing the OAuth token
	Token *OAuthTokenResult

	SubDomain    string
	ClientID     string
	ClientSecret string
}

// NewAPI create the main API object
func NewAPI(shard string, clientID string, clientSecret string, subdomain string, loglevel logging.Level) *API {
	ol := API{Shard: shard, ClientID: clientID, ClientSecret: clientSecret, SubDomain: subdomain}
	ol.SetLogLevel(loglevel)
	return &ol
}

// GetURL creates a URL given the URI and any given args.
// Returns a URL
func (o *API) GetURL(uri string, args ...string) string {
	// Handle cases where the uri requires variable replacements (ie.  /api/1/user/%d/roles)
	fulluri := uri

	if len(args) > 0 {
		// Convert to slice of interface so that the slice can be sent in as a variadic argument
		argint := make([]interface{}, len(args))
		for index, value := range args {
			argint[index] = value
		}
		fulluri = fmt.Sprintf(uri, argint...)
	}

	if o.CustomURL != "" {
		return fmt.Sprintf("%s/%s", o.CustomURL, fulluri)
	}

	return fmt.Sprintf("%s/%s", fmt.Sprintf(OneLoginURL, o.Shard), fulluri)
}

// SetLogLevel define the log level
func (o *API) SetLogLevel(loglevel logging.Level) {
	SetLogLevel(loglevel)
}

// ObtainAPIAccess initialize the access to the API.
func (o *API) ObtainAPIAccess() (err error) {
	o.Token = NewOAuthTokenResult()

	return o.Token.Obtain(o)
}

// SAMLAuthenticate used to authenticate a user thanks to SAML
func (o *API) SAMLAuthenticate(user, pass, appID, ip string, mfa, deviceIndex int) (result *AwsSAMLAssertion, err error) {
	result = NewAwsSAMLAssertion(user, pass)
	assertion := NewSAMLAssertionResult()
	_, err = assertion.Post(o, user, pass, appID, o.SubDomain, ip)

	if err != nil {
		return
	}
	if assertion.Status.Error {
		err = fmt.Errorf("%d: %s", assertion.Status.Code, assertion.Status.Message)
		return
	}
	if assertion.Status.Type == "success" && assertion.Status.Message == "success" {
		result.SetDecoded(assertion.Data)
		return
	}
	err = json.Unmarshal(assertion.Data, &assertion.data)
	if err != nil {
		return
	}

	// Display list of MFA Devices to use
	fmt.Printf("\nMFA Required\n")
	var device SAMLAssertionDevice
	if deviceIndex == -1 {
		fmt.Print("Authenticate using one of these devices:\n")
		fmt.Print("-----------------------------------------------------------------------\n")
		for index, device := range assertion.data[0].Devices {
			fmt.Printf(" %d | %s\n", index, device.DeviceType)
		}
		fmt.Print("-----------------------------------------------------------------------\n")

		device = assertion.data[0].Devices[common.Select(0, len(assertion.data[0].Devices)-1)]
	} else {
		if deviceIndex >= len(assertion.data[0].Devices) {
			err = fmt.Errorf("Invalid index %d. It must be between 0 and %d", deviceIndex, len(assertion.data[0].Devices)-1)
			return
		}
		device = assertion.data[0].Devices[deviceIndex]
		fmt.Printf("Using the MFA device index %d (%s)\n", deviceIndex, device.DeviceType)
	}
	result.MfaVerifyInfo.DeviceID = device.DeviceID
	result.MfaVerifyInfo.DeviceType = device.DeviceType

	verifyFactor := NewVerifyFactorResult()

	var MFACode int

	defer fmt.Printf("\n")

	switch device.DeviceType {
	case "OneLogin SMS":
		fmt.Printf("SMS with OTP token sent to device %d\n", device.DeviceID)
		verifyFactor.Post(o, appID, device.DeviceID, assertion.data[0].StateToken, "", true)
		MFACode = common.GetNumber("Enter the SMS OTP code received:")
	case "OneLogin Protect":
		fmt.Printf("PUSH with OTP token sent to device %d\n", device.DeviceID)
		_, err = verifyFactor.Post(o, appID, device.DeviceID, assertion.data[0].StateToken, "", false)
		// Push. Need to wait for OneLogin to confirm.
		time.Sleep(time.Second * TimeSleepOnResponsePending)
		for i := 0; i < MaxIterGetSAMLResponse; i++ {
			fmt.Printf(".")
			_, err = verifyFactor.Post(o, appID, device.DeviceID, assertion.data[0].StateToken, "", true)
			if err != nil {
				return
			}
			if verifyFactor.Status.Error {
				err = fmt.Errorf("%d: %s", verifyFactor.Status.Code, verifyFactor.Status.Message)
			} else if verifyFactor.Status.Type == "success" {
				result.SetDecoded([]byte(verifyFactor.Data))
				return
			}

			// recheck in couple of seconds
			time.Sleep(time.Second * TimeSleepOnResponsePending)
		}
		fmt.Printf("\nUnable to get your device (%d) authentication.\n", device.DeviceID)
		verifyFactor.Post(o, appID, device.DeviceID, assertion.data[0].StateToken, "", true)
		MFACode = common.GetNumber("Enter the OneProtect OTP code from your mobile application:")
		return

	default:
		fmt.Printf("Retrieve the OTP token from your device %d\n", device.DeviceID)
		if mfa != -1 {
			MFACode = mfa
		} else {
			MFACode = common.GetNumber(fmt.Sprintf("Enter the %s OTP code:", device.DeviceType))
		}
	}
	result.MfaVerifyInfo.OTPToken = MFACode
	_, err = verifyFactor.Post(o, appID, device.DeviceID, assertion.data[0].StateToken, fmt.Sprintf("%d", MFACode), true)
	if verifyFactor.Status.Type == "success" {
		result.SetDecoded([]byte(verifyFactor.Data))
		return
	}
	if verifyFactor.Status.Error {
		err = fmt.Errorf("%s: %s", verifyFactor.Status.Type, verifyFactor.Status.Message)
	}
	return
}

func (o *API) getBearerHeaders() (ret common.Headers) {
	ret = GetHeaders("bearer:" + o.Token.AccessToken)
	return
}
