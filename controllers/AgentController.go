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

	// Binding JSON
	if err := ctx.ShouldBindJSON(&newAgentNodeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati utente non validi"})
		return
	}

	// Recupero targetRole dal contesto
	role, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato nella richiesta"})
		return
	}

	// Verifica che agent sia popolato
	if newAgentNodeReq.Agent == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Il campo 'agent' è obbligatorio"})
		return
	}

	// verifica il formato del ruolo e lo inserisce alla request
	roleStr, ok := role.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Formato ruolo non valido"})
		return
	}
	newAgentNodeReq.Agent.Role = roleStr

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// invocazione del service per il salvataggio
	response, err := c.nodeService.CreateNode(&newAgentNodeReq, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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
