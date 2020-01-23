package api

// User struct
type User struct {
	Email         string            `json:"email"`
	ID            int64             `json:"id"`
	Status        int               `json:"status"`
	State         int               `json:"state"`
	RolesID       []int64           `json:"role_id"`
	ManagerUserID int64             `json:"manager_user_id"`
	MemberOf      string            `json:"member_of"`
	CustomAttrs   map[string]string `json:"custom_attributes"`
}
