package controllers

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgentOperatorController struct {
	service services.AgentOperatorService
}

func NewAgentOperatorController(service services.AgentOperatorService) *AgentOperatorController {
	return &AgentOperatorController{service: service}
}

func (c *AgentOperatorController) AssignAgentToOperator(ctx *gin.Context) {
	agentIdStr := ctx.Param("agentId")
	uid64, err := strconv.ParseUint(agentIdStr, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID agenzia non valido"})
		return
	}
	uAgentId := uint(uid64)

	operatorIdStr := ctx.Param("operatorId")
	uid64, err = strconv.ParseUint(operatorIdStr, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID operatore non valido"})
		return
	}
	uOperatorId := uint(uid64)

	newAgentOperatorReq := pivot.AssignToOpRequest{
		AgentId:    &uAgentId,
		OperatorId: uOperatorId,
	}
	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.AssignAgentToOperator(&newAgentOperatorReq, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
