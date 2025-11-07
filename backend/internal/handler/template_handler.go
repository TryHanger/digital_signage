package handler

import (
	"net/http"
	"strconv"

	"github.com/TryHanger/digital_signage/backend/internal/model"
	"github.com/TryHanger/digital_signage/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	service *service.TemplateService
}

func NewTemplateHandler(service *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

func (h *TemplateHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/templates")
	{
		// Register routes without trailing slash to avoid Gin's automatic 301 redirect
		// (redirect responses may not include CORS headers and break browser requests).
		group.GET("", h.GetAll)
		group.GET("/:id", h.GetByID)
		group.POST("", h.CreateTemplate)
		group.PUT("/:id", h.UpdateTemplate)
		group.DELETE("/:id", h.Delete)

	}
}

func (h *TemplateHandler) CreateTemplate(c *gin.Context) {
	var template model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "template created"})
}

func (h *TemplateHandler) GetAll(c *gin.Context) {
	templates, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *TemplateHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}

func (h *TemplateHandler) UpdateTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var template model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	template.ID = uint(id)

	if err := h.service.UpdateTemplate(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template updated"})
}
