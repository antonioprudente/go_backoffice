//go:build wireinject
// +build wireinject

package main

import (
	"example/go_backoffice/controllers"
	"example/go_backoffice/repositories"
	"example/go_backoffice/services"
	"os"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func provideJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

// InitUserController dichiara la catena di dipendenze
func InitUserController(db *gorm.DB) *controllers.UserController {
	wire.Build(
		repositories.NewUserRepository,
		services.NewUserService,
		controllers.NewUserController,
	)
	return nil
}

// InitAgentNodeController dichiara la catena di dipendenze
func InitAgentController(db *gorm.DB) *controllers.AgentController {
	wire.Build(
		repositories.NewAgentNodeRepository,
		repositories.NewUserRepository,
		services.NewAgentNodeService,
		services.NewUserService,
		controllers.NewAgentController,
	)
	return nil
}

func InitAuthController(db *gorm.DB) *controllers.AuthController {
	wire.Build(
		repositories.NewUserRepository,
		provideJWTSecret,
		services.NewAuthService,
		controllers.NewAuthController,
	)
	return nil
}
