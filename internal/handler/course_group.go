package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/internal/service"
)

type CourseGroupHandler interface {
	RouteRegister
}

type courseGroupHandler struct {
	*Handler
	service service.CourseGroupService
}

func NewCourseGroupHandler(handler *Handler, service service.CourseGroupService) CourseGroupHandler {
	return &courseGroupHandler{Handler: handler, service: service}
}

func (h *courseGroupHandler) RegisterRoute(router fiber.Router) {
	router.Get("/group", h.GetCourseGroupV1)
}

// GetCourseGroupV1 godoc
// @Summary /group/{group_id}
// @Description get a course group, old version or v1 version
// @Tags CourseGroup
// @Accept json
// @Produce json
// @Deprecated
// @Router /group/{id} [get]
// @Param id path string true "course group id"
// @Success 200 {object} schema.CourseGroupV1Response
// @Failure 400 {object} schema.HttpError
// @Failure 404 {object} schema.HttpBaseError
// @Failure 500 {object} schema.HttpBaseError
func (h *courseGroupHandler) GetCourseGroupV1(c *fiber.Ctx) (err error) {
	// TODO
	return c.JSON(schema.CourseGroupV1Response{})
}
