package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"example/go_backoffice/models"
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

func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero degli utenti"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}

	user, err := c.service.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Utente non trovato"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero dell'utente"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var newUser models.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati utente non validi"})
		return
	}

	if err := c.service.CreateUser(&newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Impossibile creare l'utente"})
		return
	}

	ctx.JSON(http.StatusCreated, newUser)
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
