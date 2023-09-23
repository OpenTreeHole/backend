package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/schema"
)

type DivisionHandler interface {
	RouteRegister
}

type divisionHandler struct {
	*Handler
}

func NewDivisionHandler(handler *Handler) DivisionHandler {
	return &divisionHandler{Handler: handler}
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
	// TODO:
	return c.JSON([]schema.DivisionResponse{})
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
	// TODO:
	return c.JSON(schema.DivisionResponse{})
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

	// TODO: add service.CreateDivision

	return c.JSON(schema.DivisionResponse{})
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
	var body schema.DivisionModifyRequest
	err = h.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	// TODO: add service.ModifyDivision

	return c.JSON(schema.DivisionResponse{})
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

	// TODO: add service.DeleteDivision

	return c.JSON(schema.DivisionResponse{})
}
