package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/service"
)

type CourseGroupHandler interface {
	RouteRegister
}

type courseGroupHandler struct {
	*Handler
	courseGroupService service.CourseGroupService
	accountRepository  repository.AccountRepository
}

func NewCourseGroupHandler(
	handler *Handler,
	groupService service.CourseGroupService,
	accountRepository repository.AccountRepository,
) CourseGroupHandler {
	return &courseGroupHandler{
		Handler:            handler,
		courseGroupService: groupService,
		accountRepository:  accountRepository,
	}
}

func (h *courseGroupHandler) RegisterRoute(router fiber.Router) {
	router.Get("/group/:id<int>", h.GetCourseGroupV1)
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
	c.Context().SetUserValue("FiberCtx", c)

	user, err := h.accountRepository.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}

	groupID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	response, err := h.courseGroupService.GetGroupByIDV1(c.Context(), user, groupID)
	if err != nil {
		return err
	}

	return c.JSON(response)
}
