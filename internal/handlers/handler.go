package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nazarkorpal/img-task/internal/services"
)

type Handler struct {
	Image
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Image: NewImageHandler(service.Image),
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r := router.Group("api")
	h.Image.routes(r)

	return router
}

