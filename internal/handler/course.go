package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/service"
)

type CourseHandler interface {
	RouteRegister
}

type courseHandler struct {
	*Handler
	service service.CourseService
}

func NewCourseHandler(handler *Handler, service service.CourseService) CourseHandler {
	return &courseHandler{Handler: handler, service: service}
}

func (h *courseHandler) RegisterRoute(router fiber.Router) {
	// TODO
}
