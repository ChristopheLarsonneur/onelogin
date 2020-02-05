package api

// See https://developers.onelogin.com/api-docs/1/users/user-resource
const (
	StateUnapproved = iota
	StateApproved
	StateRejected
	StateUnlicensed
)

const (
	StatusUnactivated = iota
	StatusActive
	StatusSuspended
	StatusLocked
	StatusPasswordExpired
	StatusPasswordReset
	_ // index ignored. See https://developers.onelogin.com/api-docs/1/users/user-resource
	StatusPasswordPending
	StatusSecurityQuestions
)

// User struct
type User struct {
	Username      string            `json:"username"`
	Email         string            `json:"email"`
	ID            int64             `json:"id"`
	Status        int               `json:"status"`
	State         int               `json:"state"`
	RolesID       []int64           `json:"role_id"`
	ManagerUserID int64             `json:"manager_user_id"`
	MemberOf      string            `json:"member_of"`
	Firstname     string            `json:"firstname"`
	Lastname      string            `json:"lastname"`
	Department    string            `json:"department"`
	CustomAttrs   map[string]string `json:"custom_attributes"`
}
