package controllers

import (
	"errors"
	"net/http"

	"example/go_backoffice/dto/auth"
	"example/go_backoffice/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req auth.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	token, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Credenziali non valide"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il login"})
		return
	}

	ctx.JSON(http.StatusOK, auth.LoginResponse{Token: token})
}
