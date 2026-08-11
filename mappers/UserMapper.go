// mappers/UserMapper.go
package mappers

import (
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
)

// ToUserModel converte una UserRequest nel model da persistere
// ToUserModel converte una UserRequest nel model da persistere
func ToUserModel(req *user.UserRequest) *models.User {
	model := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Status:    enums.StatusActive,
		Email:     req.Email,
		Password:  req.Password,
	}

	if req.ForeignId != nil {
		model.ForeignId = req.ForeignId
	}

	return model
}

// ToUserResponse converte un model nella response da esporre
func ToUserResponse(u *models.User) user.UserResponse {
	return user.UserResponse{
		Id:        int(u.ID),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Role:      enums.Role(u.Role),
		Status:    enums.Status(u.Status),
		Email:     u.Email,
	}
}

// ToUserResponses converte uno slice di model, utile per GetAllUsers
func ToUserResponses(users []models.User) []user.UserResponse {
	res := make([]user.UserResponse, len(users))
	for i, u := range users {
		res[i] = ToUserResponse(&u)
	}
	return res
}
