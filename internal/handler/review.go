package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/service"
)

type ReviewHandler interface {
	RouteRegister
}

type reviewHandler struct {
	*Handler
	service service.ReviewService
}

func NewReviewHandler(handler *Handler, service service.ReviewService) ReviewHandler {
	return &reviewHandler{Handler: handler, service: service}
}

func (h *reviewHandler) RegisterRoute(router fiber.Router) {
	// TODO
}
