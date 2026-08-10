package dto/user

type UserResponse struct {
	ID 		  int `json:"id"`
	FirstName string`json:"first_name"`
	LastName  string`json:"last_name"`
	Username  string`json:"username"`
	Status 	  Status`json:"status"`
	Email     string`json:"email"`
}
