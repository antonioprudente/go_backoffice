package controllers

import (
	"example/go_backoffice/dto/agent_node"
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
	if err := ctx.ShouldBindJSON(&newAgentNodeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati utente non validi"})
		return
	}

	response, err := c.nodeService.CreateNode(&newAgentNodeReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Impossibile creare l'agente"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AgentController) GetAgentsTree(ctx *gin.Context) {
	agents, err := c.nodeService.GetTrees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero degli agenti"})
		return
	}
	ctx.JSON(http.StatusOK, agents)
}
