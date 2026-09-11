package controllers

import (
	"errors"
	"example/go_backoffice/dto/note"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/policies"
	"example/go_backoffice/services"
	"net/http"
	"strconv"

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
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, policies.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *NoteController) GetAllNotes(ctx *gin.Context) {
	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var filter note.NoteFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Filtri non validi"})
		return
	}

	response, err := c.service.GetAllNotes(actor, filter)
	if err != nil {
		if errors.Is(err, policies.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *NoteController) GetNoteByID(ctx *gin.Context) {
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
		return
	}

	response, err := c.service.GetNoteByID(uid, actor)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, policies.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *NoteController) UpdateNote(ctx *gin.Context) {
	id := ctx.Param("id")
	uid64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}
	uid := uint(uid64)

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

	response, err := c.service.UpdateNote(uid, &req, actor)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, policies.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NoteController) DeleteNote(ctx *gin.Context) {
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
		return
	}

	err = c.service.DeleteNote(uid, actor)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, policies.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Nota eliminata con successo"})
}
