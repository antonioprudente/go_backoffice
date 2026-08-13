package pivot

import (
	"example/go_backoffice/dto/user"
)

type AssignToOpResponse struct {
	Operator user.UserResponse `json:"operator"`
	Agent    user.UserResponse `json:"agent,omitempty"`
	Agency   user.UserResponse `json:"agency,omitempty"`
}
