package api

// SAMLAssertionDataResult describe the typical Data result
type SAMLAssertionDataResult struct {
	StateToken string `json:"state_token"`
	User       struct {
		LastName  string `json:"lastname"`
		UserName  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"firstname"`
		ID        int    `json:"id"`
	}
	Devices     []SAMLAssertionDevice
	CallbackURL string `json:"callback_url"`
}

// SAMLAssertionDevice provides in the result, thelist of devices.
type SAMLAssertionDevice struct {
	DeviceID   int    `json:"device_id"`
	DeviceType string `json:"device_type"`
}
