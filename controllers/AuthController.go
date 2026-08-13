package controllers

import (
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

	resp, err := c.service.Login(&req, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AuthController) Logout(ctx *gin.Context) {

}
