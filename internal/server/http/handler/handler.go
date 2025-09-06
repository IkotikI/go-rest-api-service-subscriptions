package handler

import (
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	g *gin.RouterGroup

	service *service.Service

	subscription *SubscriptionHandler
	swagger      *SwaggerController
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}

}

func (h *Handler) InitRoutes(g *gin.RouterGroup) {
	h.subscription = NewSubscriptionHandler(g, h.service.Subscriptions)
	h.swagger = NewSwaggerController(g)
}
