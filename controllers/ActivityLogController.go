package controllers

import (
	"errors"
	"example/go_backoffice/policies"
	"example/go_backoffice/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActivityLogController struct {
	service services.ActivityLogService
}

func NewActivityLogController(service services.ActivityLogService) *ActivityLogController {
	return &ActivityLogController{service: service}
}

func (c *ActivityLogController) ActivityLog(ctx *gin.Context) {
	res, err := c.service.GetActivity()

	if err != nil {
		if errors.Is(err, policies.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *ActivityLogController) ActivityLogDetails(ctx *gin.Context) {
	id := ctx.Param("id")
	uid64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID non valido"})
		return
	}
	uid := uint(uid64)

	res, err := c.service.GetLogDetailByID(uid)
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
	ctx.JSON(http.StatusOK, res)

}
