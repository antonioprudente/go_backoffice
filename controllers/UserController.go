package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/policies"
	"example/go_backoffice/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var newUserReq user.UserRequest
	// ottiene il targetRole dal middleware (su quale tipo di entità si sta operando)

	role, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato nella richiesta"})
		return
	}
	newUserReq.Role = role.(string)

	// validazione request
	if err := ctx.ShouldBindJSON(&newUserReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati utente non validi"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// creazione a db (le policy di autorizzazione fine sono applicate dentro il service)
	response, err := c.service.CreateUser(&newUserReq, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	// ottiene il targetRole dal middleware
	roleVal, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato nella richiesta"})
		return
	}
	targetRole := roleVal.(string)

	var users []user.UserResponse

	users, err := c.service.GetAllByRole(targetRole)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero degli utenti"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	uid64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}
	uid := uint(uid64)

	targetRole, exists := ctx.Get("targetRole")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ruolo non specificato nella richiesta"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.GetUserByIDAndRole(uid, targetRole.(string), actor)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Utente non trovato"})
		case errors.Is(err, policies.ErrForbidden):
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			// ErrUnknownRole, ErrMissingRelation, ErrNotImplemented o errori
			// tecnici veri (query fallita ecc.): questi restano 500
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero dell'utente"})
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) ActiveUserById(ctx *gin.Context) {
	idStr := ctx.Param("id")

	parsedID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}
	userID := uint(parsedID)
	response, err := c.service.ChangeStatus(userID, enums.StatusActive)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Utente non trovato"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante l'aggiornamento dello stato dell'utente"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) SuspendUserById(ctx *gin.Context) {
	idStr := ctx.Param("id")

	parsedID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}
	userID := uint(parsedID)
	response, err := c.service.ChangeStatus(userID, enums.StatusSuspended)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Utente non trovato"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante l'aggiornamento dello stato dell'utente"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) BlockUserById(ctx *gin.Context) {
	idStr := ctx.Param("id")

	parsedID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}
	userID := uint(parsedID)
	response, err := c.service.ChangeStatus(userID, enums.StatusBlocked)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Utente non trovato"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante l'aggiornamento dello stato dell'utente"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}

	if err := c.service.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Impossibile eliminare l'utente"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Utente eliminato con successo"})
}
