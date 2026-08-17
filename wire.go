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
		policies.NewUserPolicy,
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
		policies.NewUserPolicy,
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

func InitAgentOperatorController(db *gorm.DB) *controllers.AgentOperatorController {
	wire.Build(
		repositories.NewAgentOperatorRepository,
		repositories.NewUserRepository,
		services.NewAgentOperatorService,
		controllers.NewAgentOperatorController,
	)
	return nil
}

func InitAgencyOperatorController(db *gorm.DB) *controllers.AgencyOperatorController {
	wire.Build(
		repositories.NewAgencyOperatorRepository,
		repositories.NewUserRepository,
		services.NewAgencyOperatorService,
		controllers.NewAgencyOperatorController,
	)
	return nil
}
