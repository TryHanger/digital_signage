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

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/schedules")
	{
		group.POST("/", h.Create)
		group.GET("/", h.GetAll)
		group.GET("/:id", h.GetByID)
		group.DELETE("/:id", h.Delete)
		group.PUT("/resolve", h.ResolveConflicts)
		group.PUT("/update", h.UpdateSchedules)
	}
}

// Создание расписания
func (h *ScheduleHandler) Create(c *gin.Context) {
	var schedule model.Schedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conflicts, err := h.service.CreateSchedule(&schedule)
	if err != nil {
		if conflicts != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error":     "time conflict detected",
				"conflicts": conflicts,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

// Получение всех
func (h *ScheduleHandler) GetAll(c *gin.Context) {
	schedules, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

// Получение по ID
func (h *ScheduleHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	schedule, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Расписание не найдено"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

// Удаление
func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteSchedule(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ScheduleHandler) ResolveConflicts(c *gin.Context) {
	var req struct {
		Schedules []model.Schedule `json:"schedules"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSchedules(req.Schedules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedules updated"})
}

func (h *ScheduleHandler) UpdateSchedules(c *gin.Context) {
	var schedules []model.Schedule

	if err := c.ShouldBindJSON(&schedules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSchedules(schedules); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "schedules updated successfully"})
}
