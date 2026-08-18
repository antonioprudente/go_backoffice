package controllers

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgencyOperatorController struct {
	service services.AgencyOperatorService
}

func NewAgencyOperatorController(service services.AgencyOperatorService) *AgencyOperatorController {
	return &AgencyOperatorController{service: service}
}

func (c *AgencyOperatorController) AssignAgencyToOperator(ctx *gin.Context) {
	var request pivot.AssignToOpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.AssignAgencyToOperator(&request, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AgencyOperatorController) AssignAgenciesToOperator(ctx *gin.Context) {
	var request pivot.ArraysToOpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Dati non validi"})
		return
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.AssignAgenciesToOperator(&request, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, response)
}
