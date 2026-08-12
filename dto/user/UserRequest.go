// dto/user/UserRequest.go
package user // non "dto/user"

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	ForeignId *uint  `json:"foreign_id,omitempty"`
}
