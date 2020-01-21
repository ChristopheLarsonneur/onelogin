package onelogin

import (
	"encoding/json"
	"errors"
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
	lastError error

	allRolesLoaded bool
	roles map[int64]string
}

// NewService create the main API object
func NewService(shard string, clientID string, clientSecret string, subdomain string, loglevel logging.Level) (ret *Service) {
	ret = new(Service)
	ret.core = api.NewAPI(shard, clientID, clientSecret, subdomain)
	ret.SetLogLevel(loglevel)

	ret.roles = make(map[int64]string)
	return
}

// SetLogLevel define the log level
func (o *Service) SetLogLevel(loglevel logging.Level) {
	SetLogLevel(loglevel)
}

// SAMLAuthenticate used to authenticate a user thanks to SAML
func (o *Service) SAMLAuthenticate(user, pass, appID, ip string, mfa, deviceIndex int) (result *AwsSAMLAssertion, err error) {
	if err = o.initCheck() ; err != nil {
		return
	}

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

// CleanError cleanup the last error reported by onelogin API call.
func (o *Service) CleanError() {
	if o == nil {
		return
	}
	o.lastError = nil
}

func (o *Service) setError(err error) error {
	if o == nil {
		return errors.New("onelogin.Service is nil")
	}
	o.lastError = err
	return err
}

// initCheck basically check initial onelogin object status and obtain API access.
func (o *Service) initCheck() (_ error) {
	if o == nil {
		return errors.New("onelogin.Service is nil")
	}
	if  o.core == nil {
		return errors.New("onelogin.api.Core is nil")
	}

	if o.lastError != nil {
		return fmt.Errorf("onelogin.Core is always in error: %s", o.lastError)
	}

	if o.core.Token == nil || o.core.Token.AccessToken == "" {
		if err := o.setError(o.core.ObtainAPIAccess()) ; err != nil {
			return err
		}
	}
	return
}

// GetRoles return the list of all roles from OneLogin
func (o *Service) GetRoles() (ret map[int64]string, err error) {
	if err = o.initCheck() ; err != nil {
		return
	}

	if o.allRolesLoaded {
		ret = o.roles
		return 
	}

	roles := api.NewGetRoles()
	if _, err = roles.Get(o.core) ; err != nil {
		return ret, o.setError(err)
	}

	for _, role := range roles.Data {
		o.roles[role.ID] = role.Name
	}
	ret = o.roles
	return
}

// GetAPI provide the OneLogin api obejct and access to it. (access token)
func (o *Service) GetAPI() (apiCore *api.Core, err error) {
	if err = o.initCheck() ; err != nil {
		return
	}
	apiCore = o.core
	return
}

// GetRoleName return a role name from the role ID
func (o *Service) GetRoleName(id int64) (ret string, err error) {
	if err = o.initCheck() ; err != nil {
		return
	}

	if v, found := o.roles[id] ; found {
		return v, nil
	}

	role := api.NewGetRoleByID()

	if _, err = role.Get(o.core, id) ; err != nil {
		return ret, err
	}
	if len(role.Data) >= 1 {
		ret = role.Data[0].Name
		o.roles[id] = ret
	}
	return
}
