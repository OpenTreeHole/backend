package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/internal/service"
)

type DivisionHandler interface {
	RouteRegister
}

type divisionHandler struct {
	*Handler
	service service.DivisionService
	account repository.AccountRepository
}

func NewDivisionHandler(handler *Handler, service service.DivisionService, account repository.AccountRepository) DivisionHandler {
	return &divisionHandler{Handler: handler, service: service, account: account}
}

func (h *divisionHandler) RegisterRoute(router fiber.Router) {
	router.Get("/divisions", h.ListDivisions)
	router.Get("/divisions/:id", h.GetDivision)
	router.Post("/divisions", h.CreateDivision)
	router.Put("/divisions/:id", h.ModifyDivision)
	router.Delete("/divisions/:id", h.DeleteDivision)
}

// ListDivisions godoc
// @Summary list all divisions
// @Description list all divisions
// @Tags Division
// @Accept json
// @Produce json
// @Router /divisions [get]
// @Success 200 {array} schema.DivisionResponse
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *divisionHandler) ListDivisions(c *fiber.Ctx) (err error) {
	divisions, err := h.service.ListDivisions(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(divisions)
}

// GetDivision godoc
// @Summary get a division
// @Description get a division
// @Tags Division
// @Accept json
// @Produce json
// @Router /divisions/{id} [get]
// @Param id path string true "division id"
// @Success 200 {object} schema.DivisionResponse
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *divisionHandler) GetDivision(c *fiber.Ctx) (err error) {

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	division, err := h.service.GetDivision(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(division)
}

// CreateDivision godoc
// @Summary create a division
// @Description create a division, only admin can create
// @Tags Division
// @Accept json
// @Produce json
// @Router /divisions [post]
// @Param json body schema.DivisionCreateRequest true "json"
// @Success 201 {object} schema.DivisionResponse
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *divisionHandler) CreateDivision(c *fiber.Ctx) (err error) {
	var body schema.DivisionCreateRequest
	err = h.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	c.Context().SetUserValue("FiberCtx", c)
	user, err := h.account.GetCurrentUser(c.Context())

	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return c.JSON(schema.Forbidden())
	}

	response, err := h.service.CreateDivision(c.Context(), &body)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// ModifyDivision godoc
// @Summary modify a division
// @Description modify a division, only admin can modify
// @Tags Division
// @Accept json
// @Produce json
// @Router /divisions/{id} [put]
// @Param id path string true "division id"
// @Param json body schema.DivisionModifyRequest true "json"
// @Success 200 {object} schema.DivisionResponse
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *divisionHandler) ModifyDivision(c *fiber.Ctx) (err error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var body schema.DivisionModifyRequest
	err = h.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	c.Context().SetUserValue("FiberCtx", c)
	user, err := h.account.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return c.JSON(schema.Forbidden())
	}

	response, err := h.service.ModifyDivision(c.Context(), id, &body)
	if err != nil {
		return err
	}

	return c.JSON(response)
}

// DeleteDivision godoc
// @Summary delete a division
// @Description delete a division, only admin can delete
// @Tags Division
// @Accept json
// @Produce json
// @Router /divisions/{id} [delete]
// @Param id path string true "division id"
// @Param json body schema.DivisionDeleteRequest true "json"
// @Success 204 {object} schema.HttpBaseError
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *divisionHandler) DeleteDivision(c *fiber.Ctx) (err error) {
	var body schema.DivisionDeleteRequest
	err = h.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	c.Context().SetUserValue("FiberCtx", c)
	user, err := h.account.GetCurrentUser(c.Context())
	if err != nil {
		return err
	}
	if !user.IsAdmin {
		return c.JSON(schema.Forbidden())
	}

	_, err = h.service.DeleteDivision(c.Context(), id, &body)
	if err != nil {
		return err
	}

	return c.JSON(nil)
}
