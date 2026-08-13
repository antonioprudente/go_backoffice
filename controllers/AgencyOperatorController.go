package controllers

import (
	"example/go_backoffice/services"
)

type AgentOperatorController struct {
	service services.AgentOperatorService
}

func NewAgentOperatorController(service services.AgentOperatorService) *AgentOperatorController {
	return &AgentOperatorController{service: service}
}
