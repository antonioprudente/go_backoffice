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

func (c *AgentOperatorController) RemoveAgentFromOperator(ctx *gin.Context) {
	agId := ctx.Param("agentID")
	uid64, err := strconv.ParseUint(agId, 10, 64) // Usiamo 64 specificamente per ParseUint
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "agentID non valido"})
		return
	}
	agentID := uint(uid64)

	opId := ctx.Param("operatorID")
	uid64, err = strconv.ParseUint(opId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "operatorID non valido"})
		return
	}
	operatorID := uint(uid64)

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.RemoveAgentFromOperator(&agentID, &operatorID, actor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Controlla se res è nil oppure se punta a false
	if res == nil || !*res {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Nessuna associazione trovata da eliminare",
			"data":    false,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Agente rimosso con successo",
		"data":    true,
	})
}
