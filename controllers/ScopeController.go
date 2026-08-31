package controllers

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScopeController struct {
	service services.ScopeService
}

func NewScopeControllerController(service services.ScopeService) *ScopeController {
	return &ScopeController{service: service}
}

func (c *ScopeController) AssignToOperator(ctx *gin.Context) {
	var request pivot.ArraysToOpRequest
	response, err := c.service.AssignToOperator(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}
