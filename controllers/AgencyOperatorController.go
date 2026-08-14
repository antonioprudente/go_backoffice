package controllers

import (
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/middlewares"
	"example/go_backoffice/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AgencyOperatorController struct {
	service services.AgencyOperatorService
}

func NewAgencyOperatorController(service services.AgencyOperatorService) *AgencyOperatorController {
	return &AgencyOperatorController{service: service}
}

func (c *AgencyOperatorController) AssignAgencyToOperator(ctx *gin.Context) {
	agencyIdStr := ctx.Param("agencyId")
	uid64, err := strconv.ParseUint(agencyIdStr, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID agenzia non valido"})
		return
	}
	uAgencyId := uint(uid64)

	operatorIdStr := ctx.Param("operatorId")
	uid64, err = strconv.ParseUint(operatorIdStr, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID operatore non valido"})
		return
	}
	uOperatorId := uint(uid64)

	newAgencyOperatorReq := pivot.AssignToOpRequest{
		AgencyId:   &uAgencyId,
		OperatorId: uOperatorId,
	}

	actor, err := middlewares.ActorFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.AssignAgencyToOperator(&newAgencyOperatorReq, actor)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
