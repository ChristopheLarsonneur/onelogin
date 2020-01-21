package onelogin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/clarsonneur/onelogin/api"
	"github.com/clarsonneur/onelogin/common"
	"github.com/op/go-logging"
)

// OneLoginURL is the default OneLogin API endpoint
const (
	OneLoginURL                = "https://api.%s.onelogin.com"
	TimeSleepOnResponsePending = 15
	MaxIterGetSAMLResponse     = 6
)

// Service is the core OneLogin service object, connected to the OneLogin Service through the API (api.Core).
type Service struct {
	core *api.Core

}

// NewService create the main API object
func NewService(shard string, clientID string, clientSecret string, subdomain string, loglevel logging.Level) (ret *Service) {
	ret = new(Service)
	ret.core = api.NewAPI(shard, clientID, clientSecret, subdomain)
	ret.SetLogLevel(loglevel)

	return
}

// SetLogLevel define the log level
func (o *Service) SetLogLevel(loglevel logging.Level) {
	SetLogLevel(loglevel)
}

// SAMLAuthenticate used to authenticate a user thanks to SAML
func (o *Service) SAMLAuthenticate(user, pass, appID, ip string, mfa, deviceIndex int) (result *AwsSAMLAssertion, err error) {
	result = NewAwsSAMLAssertion(user, pass)
	assertion := api.NewSAMLAssertionResult()
	_, err = assertion.Post(o.core, user, pass, appID, o.core.SubDomain, ip)

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

	var data []api.SAMLAssertionDataResult
	err = json.Unmarshal(assertion.Data, &data)
	if err != nil {
		return
	}

	// Display list of MFA Devices to use
	fmt.Printf("\nMFA Required\n")
	var device api.SAMLAssertionDevice
	if deviceIndex == -1 {
		fmt.Print("Authenticate using one of these devices:\n")
		fmt.Print("-----------------------------------------------------------------------\n")
		for index, device := range data[0].Devices {
			fmt.Printf(" %d | %s\n", index, device.DeviceType)
		}
		fmt.Print("-----------------------------------------------------------------------\n")

		device = data[0].Devices[common.Select(0, len(data[0].Devices)-1)]
	} else {
		if deviceIndex >= len(data[0].Devices) {
			err = fmt.Errorf("Invalid index %d. It must be between 0 and %d", deviceIndex, len(data[0].Devices)-1)
			return
		}
		device = data[0].Devices[deviceIndex]
		fmt.Printf("Using the MFA device index %d (%s)\n", deviceIndex, device.DeviceType)
	}
	result.MfaVerifyInfo.DeviceID = device.DeviceID
	result.MfaVerifyInfo.DeviceType = device.DeviceType

	verifyFactor := api.NewVerifyFactorResult()

	var MFACode int

	defer fmt.Printf("\n")

	switch device.DeviceType {
	case "OneLogin SMS":
		fmt.Printf("SMS with OTP token sent to device %d\n", device.DeviceID)
		verifyFactor.Post(o.core, appID, device.DeviceID, data[0].StateToken, "", true)
		MFACode = common.GetNumber("Enter the SMS OTP code received:")
	case "OneLogin Protect":
		fmt.Printf("PUSH with OTP token sent to device %d\n", device.DeviceID)
		_, err = verifyFactor.Post(o.core, appID, device.DeviceID, data[0].StateToken, "", false)
		// Push. Need to wait for OneLogin to confirm.
		time.Sleep(time.Second * TimeSleepOnResponsePending)
		for i := 0; i < MaxIterGetSAMLResponse; i++ {
			fmt.Printf(".")
			_, err = verifyFactor.Post(o.core, appID, device.DeviceID, data[0].StateToken, "", true)
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
		verifyFactor.Post(o.core, appID, device.DeviceID, data[0].StateToken, "", true)
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
	_, err = verifyFactor.Post(o.core, appID, device.DeviceID, data[0].StateToken, fmt.Sprintf("%d", MFACode), true)
	if verifyFactor.Status.Type == "success" {
		result.SetDecoded([]byte(verifyFactor.Data))
		return
	}
	if verifyFactor.Status.Error {
		err = fmt.Errorf("%s: %s", verifyFactor.Status.Type, verifyFactor.Status.Message)
	}
	return
}
