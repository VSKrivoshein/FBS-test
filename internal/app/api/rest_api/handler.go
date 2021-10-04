package rest_api

import (
	"github.com/VSKrivoshein/FBS-test/internal/app/services/fiboncci"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service fiboncci.Service
}

func NewHandler(service fiboncci.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(JSONLogMiddleware())
	r.Use(ErrorMiddleware())

	r.POST("/fibonacci", h.fibonacci)

	return r
}