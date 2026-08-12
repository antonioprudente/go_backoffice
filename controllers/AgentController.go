package controllers

import (
	"example/go_backoffice/dto/agent_node"
	"example/go_backoffice/dto/user"
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

	// 1. Bind del JSON inviato nella request
	if err := ctx.ShouldBindJSON(&newAgentNodeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	// 2. Recupero del ruolo impostato dal middleware
	role, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato"})
		return
	}

	roleStr, ok := role.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Formato ruolo non valido"})
		return
	}

	// 3. Verifica ed inizializzazione del puntatore Agent
	if newAgentNodeReq.Agent == nil {
		newAgentNodeReq.Agent = &user.UserRequest{} // Usa user.UserRequest
	}

	// 4. Assegnazione del ruolo
	newAgentNodeReq.Agent.Role = roleStr

	// 5. Chiamata al Service
	response, err := c.nodeService.CreateNode(&newAgentNodeReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante la creazione"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

/*func (c *AgentController) CreateAgentNode(ctx *gin.Context) {
	var newAgentNodeReq agent_node.AgentNodeRequest

	// recupera il ruolo
	role, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato nella richiesta"})
		return
	}
	newAgentNodeReq.Agent.Role = role.(string)

	// validazione campi
	if err := ctx.ShouldBindJSON(&newAgentNodeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati utente non validi"})
		return
	}

	// creazione a db
	response, err := c.nodeService.CreateNode(&newAgentNodeReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Impossibile creare l'agente"})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}*/

func (c *AgentController) GetAgentsTree(ctx *gin.Context) {
	agents, err := c.nodeService.GetTrees()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero degli agenti"})
		return
	}
	ctx.JSON(http.StatusOK, agents)
}
