package user

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ForeignId string `json:"foreign_id"`
}
