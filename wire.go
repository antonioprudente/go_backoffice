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
		repositories.NewAgencyOperatorRepository,
		repositories.NewActivityLogRepository,
		policies.NewUserPolicy,
		services.NewUserService,
		services.NewActivityLogService,
		controllers.NewUserController,
	)
	return nil
}

func InitAgentController(db *gorm.DB) *controllers.AgentController {
	wire.Build(
		repositories.NewAgentNodeRepository,
		repositories.NewUserRepository,
		repositories.NewScopeRepository,
		repositories.NewAgentOperatorRepository,
		repositories.NewAgencyOperatorRepository,
		repositories.NewActivityLogRepository,
		policies.NewUserPolicy,
		services.NewAgentNodeService,
		services.NewUserService,
		services.NewActivityLogService,
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

func InitScopeController(db *gorm.DB) *controllers.ScopeController {
	wire.Build(
		repositories.NewScopeRepository,
		repositories.NewUserRepository,
		repositories.NewActivityLogRepository,
		services.NewScopeService,
		services.NewActivityLogService,
		controllers.NewScopeController,
	)
	return nil
}

func InitNoteController(db *gorm.DB) *controllers.NoteController {
	wire.Build(
		repositories.NewNoteRepository,
		repositories.NewActivityLogRepository,
		repositories.NewUserRepository,
		repositories.NewScopeRepository,
		policies.NewUserPolicy,
		services.NewNoteService,
		services.NewActivityLogService,
		controllers.NewNoteController,
	)
	return nil
}

func InitActivityLogController(db *gorm.DB) *controllers.ActivityLogController {
	wire.Build(
		repositories.NewActivityLogRepository,
		services.NewActivityLogService,
		controllers.NewActivityLogController,
	)
	return nil
}
