//go:build wireinject
// +build wireinject

package main

import (
	"example/go_backoffice/controllers"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"
	"example/go_backoffice/services"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitUserController(db *gorm.DB) *controllers.UserController {
	wire.Build(
		repositories.NewUserRepository,
		repositories.NewScopeRepository,
		repositories.NewAgentNodeRepository,
		policies.NewAgencyPolicy,
		policies.NewScopePolicy,
		services.NewUserService,
		controllers.NewUserController,
	)
	return nil
}

func InitAgentController(db *gorm.DB) *controllers.AgentController {
	wire.Build(
		repositories.NewAgentNodeRepository,
		repositories.NewUserRepository,
		repositories.NewScopeRepository,
		policies.NewAgentPolicy,
		policies.NewScopePolicy,
		policies.NewAgencyPolicy,
		services.NewAgentNodeService,
		services.NewUserService,
		controllers.NewAgentController,
	)
	return nil
}

func InitAuthController(db *gorm.DB) *controllers.AuthController {
	wire.Build(
		repositories.NewUserRepository,
		services.NewAuthService,
		controllers.NewAuthController,
	)
	return nil
}
