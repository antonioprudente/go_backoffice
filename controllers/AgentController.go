package controllers

import (
	"errors"
	"example/go_backoffice/dto/user"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/policies"
	"example/go_backoffice/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	var newUserRequest user.UserRequest

	if err := ctx.ShouldBindJSON(&newUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati utente non validi"})
		return
	}

	role, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato nella richiesta"})
		return
	}
	roleStr, ok := role.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Formato ruolo non valido"})
		return
	}
	newUserRequest.Role = roleStr

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.nodeService.CreateNode(&newUserRequest, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AgentController) GetFilteredTree(ctx *gin.Context) {
	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	agents, err := c.nodeService.GetTree(actor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero degli agenti"})
		return
	}
	ctx.JSON(http.StatusOK, agents)
}

func (c *AgentController) DeleteAgent(ctx *gin.Context) {
	id := ctx.Param("id")
	uid64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}
	uid := uint(uid64)

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	err = c.nodeService.DeleteNode(uid, actor)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, policies.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Agente eliminato con successo"})
}
