package controllers

import (
	"example/go_backoffice/services"
)

type AgencyOperatorController struct {
	service services.AgentOperatorService
}

func NewAgencyOperatorController(service services.AgentOperatorService) *AgentOperatorController {
	return &AgentOperatorController{service: service}
}
