// dto/user/UserRequest.go
package user // non "dto/user"

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Role      string `json:"role,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ForeignId *uint  `json:"foreign_id,omitempty"`
}
