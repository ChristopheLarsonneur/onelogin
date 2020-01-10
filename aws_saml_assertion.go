package onelogin

import (
	"encoding/base64"
	"fmt"
)

// AwsSAMLAssertion provide information back to the caller.
type AwsSAMLAssertion struct {
	SamlResponse        []byte
	EncodedSamlResponse []byte
	MfaVerifyInfo       struct {
		DeviceID   int
		DeviceType string
		OTPToken   int
	}
	User        string
	Password    string
	OLSubdomain string
}

// NewAwsSAMLAssertion creates the AwsSAMLAssertion object.
func NewAwsSAMLAssertion(user, pass string) (ret *AwsSAMLAssertion) {
	ret = new(AwsSAMLAssertion)
	ret.User = user
	ret.Password = pass
	return
}

// SetDecoded save the response data field, base64 decoded
func (a *AwsSAMLAssertion) SetDecoded(data []byte) (err error) {
	if a == nil {
		return fmt.Errorf("SetDecode: AwsSAMLAssertion object is nil")
	}
	a.EncodedSamlResponse = data
	a.SamlResponse, err = base64.RawStdEncoding.DecodeString(string(data))
	return
}
