package controllers

import (
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgentController struct {
	userService services.UserService
	nodeService services.AgentNodeService
}

func NewAgentController(userService services.UserService, nodeService services.AgentNodeService) *AgentController {
	return &AgentController{
		userService: userService,
		nodeService: nodeService,
	}
}

func (c *AgentController) CreateAgentNode(ctx *gin.Context) {
	var newAgentNodeReq agent_node.AgentNodeRequest

	role, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato nella richiesta"})
		return
	}
	newAgentNodeReq.Agent.Role = role.(string)

	if err := ctx.ShouldBindJSON(&newAgentNodeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati utente non validi"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.nodeService.CreateNode(&newAgentNodeReq, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AgentController) GetAgentsTree(ctx *gin.Context) {
	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	agents, err := c.nodeService.GetTrees(actor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero degli agenti"})
		return
	}
	ctx.JSON(http.StatusOK, agents)
}
