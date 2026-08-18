package controllers

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgentOperatorController struct {
	service services.AgentOperatorService
}

func NewAgentOperatorController(service services.AgentOperatorService) *AgentOperatorController {
	return &AgentOperatorController{service: service}
}

func (c *AgentOperatorController) AssignAgentToOperator(ctx *gin.Context) {
	var request pivot.AssignToOpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.AssignAgentToOperator(&request, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AgentOperatorController) AssignAgentsToOperator(ctx *gin.Context) {
	var request pivot.ArraysToOpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.AssignAgentsToOperator(&request, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response)
}
