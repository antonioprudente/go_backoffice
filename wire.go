//go:build wireinject
// +build wireinject

package main

import (
	"example/go_backoffice/controllers"
	"example/go_backoffice/repositories"
	"example/go_backoffice/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// InitUserController dichiara la catena di dipendenze
func InitUserController(db *gorm.DB) *controllers.UserController {
	wire.Build(
		repositories.NewUserRepository,
		services.NewUserService,
		controllers.NewUserController,
	)
	return nil
}
