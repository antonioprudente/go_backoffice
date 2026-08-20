// mappers/UserMapper.go
package mappers

import (
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/models"
)

// ToUserModel converte una UserRequest nel model da persistere
func ToUserModel(req *user.UserRequest) *models.User {
	model := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Role:      enums.Role(req.Role),
		Status:    enums.StatusActive,
		Email:     req.Email,
		Password:  req.Password,
	}

	if req.Id != nil {
		model.ID = *req.Id
	}

	if req.ForeignId != nil {
		model.ForeignId = req.ForeignId
	}

	return model
}

// ToUserResponse converte un model nella response da esporre
func ToUserResponse(u *models.User) user.UserResponse {
	// Aggiungi questa protezione per bloccare il crash da puntatore nil
	if u == nil {
		return user.UserResponse{}
	}

	resp := user.UserResponse{
		Id:        int(u.ID), // Riga 31 che prima andava in crash
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Role:      u.Role,
		Status:    u.Status,
		Email:     u.Email,
		ForeignId: u.ForeignId,
	}

	if u.Foreign != nil {
		foreignResp := ToUserResponse(u.Foreign)
		resp.Foreign = &foreignResp
	}

	return resp
}

// ToUserResponses converte uno slice di model, utile per GetAllUsers
func ToUserResponses(users []models.User) []user.UserResponse {
	res := make([]user.UserResponse, len(users))
	for i, u := range users {
		res[i] = ToUserResponse(&u)
	}
	return res
}
