package handler

import (
	"net/http"
	"strconv"

	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/service"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/schedules")
	{
		group.POST("/", h.CreateSchedule)
		group.GET("/", h.GetSchedules)
		group.GET("/:id", h.GetScheduleByID)
		group.DELETE("/:id", h.DeleteSchedule)
	}
}

func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var schedule model.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateSchedule(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	schedules, err := h.service.GetAllSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) GetScheduleByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	schedule, err := h.service.GetScheduleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteSchedule(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
