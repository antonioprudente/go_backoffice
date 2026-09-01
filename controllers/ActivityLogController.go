package controllers

import (
	"example/go_backoffice/models"
	"example/go_backoffice/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActivityLogController struct {
	service services.ActivityLogService
}

func NewActivityLogController(service services.ActivityLogService) *ActivityLogController {
	return &ActivityLogController{service: service}
}

func (c *ActivityLogController) ActivityLog(ctx *gin.Context) {
	var logs []*models.ActivityLog

	logs, err := c.service.GetActivity()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Errore durante il recupero dei log"})
		return
	}
	ctx.JSON(http.StatusOK, logs)
}
