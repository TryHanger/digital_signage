package handler

import (
	"github.com/TryHanger/digital_signage/internal/model"
	"github.com/TryHanger/digital_signage/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(s *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: s}
}

func (h *ScheduleHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/schedules")
	{
		group.POST("/", h.Create)
		group.GET("/", h.GetAll)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
	}

	// отдельный маршрут: активный контент для монитора
	r.GET("/monitors/:id/current", h.GetActiveContent)
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var schedule model.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, schedule)
}

func (h *ScheduleHandler) GetAll(c *gin.Context) {
	schedules, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	schedule, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var schedule model.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schedule.ID = uint(id)
	if err := h.service.Update(&schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ScheduleHandler) GetActiveContent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	schedule, err := h.service.GetActiveContent(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Нет активного контента"})
		return
	}
	c.JSON(http.StatusOK, schedule.Content)
}
