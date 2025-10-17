package handler

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	service *service.MonitorService
}

func NewMonitorHandler(service *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{service: service}
}

func (h *MonitorHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/monitors")
	group.GET("/", h.GetAll)
	group.POST("/", h.Create)

}

func (h *MonitorHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAllMonitors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *MonitorHandler) Create(c *gin.Context) {
	var monitor model.Monitor
	if err := c.ShouldBindJSON(&monitor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateMonitor(&monitor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, monitor)
}
