package controllers

import (
	"errors"
	"example/go_backoffice/dto/note"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/policies"
	"example/go_backoffice/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoteController struct {
	service services.NoteService
}

func NewNoteController(service services.NoteService) *NoteController {
	return &NoteController{service: service}
}

func (c *NoteController) AssignNote(ctx *gin.Context) {
	var req note.NoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati nota non validi"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	req.ActorID = actor.UserID

	response, err := c.service.AssignNote(&req, actor)
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
	ctx.JSON(http.StatusCreated, response)
}
