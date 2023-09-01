package handler

import (
	"monapi/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/version", h.getVersion)
		api.GET("/data", h.getData)
		api.GET("/data/debug", h.getDataDebug)
		api.GET("/last/:id", h.getLastValue)
	}

	return router
}

// Qwer-1234
