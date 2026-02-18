package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"spacex-tracker/services"
)

type LaunchHandler struct {
	service services.LaunchService
}

func NewLaunchHandler(service services.LaunchService) *LaunchHandler {
	return &LaunchHandler{
		service: service,
	}
}

func (h *LaunchHandler) GetNext(c *gin.Context) {
	launch, err := h.service.GetNext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch next launch",
		})
		return
	}

	c.JSON(http.StatusOK, launch)
}

func (h *LaunchHandler) GetLatest(c *gin.Context) {
	launch, err := h.service.GetLatest(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch latest launch",
		})
		return
	}

	c.JSON(http.StatusOK, launch)
}

func (h *LaunchHandler) GetUpcoming(c *gin.Context) {
	launches, err := h.service.GetUpcoming(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch upcoming launches",
		})
		return
	}

	c.JSON(http.StatusOK, launches)
}

func (h *LaunchHandler) GetPast(c *gin.Context) {
	sortOrder := c.DefaultQuery("sort", "desc")

	launches, err := h.service.GetPast(c.Request.Context(), sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch past launches",
		})
		return
	}

	c.JSON(http.StatusOK, launches)
}