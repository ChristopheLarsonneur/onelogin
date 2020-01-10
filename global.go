package onelogin

// ResultStatus is the common result status
type ResultStatus struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
}
