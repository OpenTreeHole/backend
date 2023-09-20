package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/internal/service"
)

type AccountHandler interface {
	RouteRegister
	Login(c *fiber.Ctx) (err error)
	Register(c *fiber.Ctx) (err error)
	ResetPassword(c *fiber.Ctx) (err error)
}

type accountHandler struct {
	*Handler
	service service.AccountService
}

func NewAccountHandler(handler *Handler, service service.AccountService) AccountHandler {
	return &accountHandler{Handler: handler, service: service}
}

func (h *accountHandler) RegisterRoute(router fiber.Router) {
	router.Post("/login", h.Login)
	router.Post("/register", h.Register)
	router.Put("/reset-password", h.ResetPassword)
}

// Login godoc
// @Summary login
// @Description login with email and password
// @Tags Account
// @Accept json
// @Produce json
// @Router /login [post]
// @Param json body schema.LoginRequest true "LoginRequest"
// @Success 200 {object} schema.TokenResponse
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *accountHandler) Login(c *fiber.Ctx) (err error) {
	var body schema.LoginRequest
	err = h.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	tokenResponse, err := h.service.Login(
		c.Context(), // TODO: create new context wrapping the request context
		body.Email,
		body.Password,
	)
	if err != nil {
		return err
	}

	return c.JSON(tokenResponse)
}

// Register godoc
// @Summary register
// @Description register with email, password and optional verification code if enabled
// @Tags Account
// @Accept json
// @Produce json
// @Router /register [post]
// @Param json body schema.RegisterRequest true "RegisterRequest"
// @Success 201 {object} schema.TokenResponse
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *accountHandler) Register(c *fiber.Ctx) (err error) {
	var body schema.RegisterRequest
	err = h.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	tokenResponse, err := h.service.Register(
		c.Context(),
		body.Email,
		body.Password,
		string(body.Verification),
		true,
	)
	if err != nil {
		return err
	}

	return c.JSON(tokenResponse)
}

// ResetPassword godoc
// @Summary reset password
// @Description reset password with email, password and optional verification code if enabled
// @Tags Account
// @Accept json
// @Produce json
// @Router /register [put]
// @Param json body schema.ResetPasswordRequest true "ResetPasswordRequest"
// @Success 201 {object} schema.TokenResponse
// @Failure 400 {object} schema.HttpError
// @Failure 500 {object} schema.HttpBaseError
func (h *accountHandler) ResetPassword(c *fiber.Ctx) (err error) {
	var body schema.ResetPasswordRequest
	err = h.ValidateBody(c, &body)
	if err != nil {
		return err
	}

	tokenResponse, err := h.service.ResetPassword(
		c.Context(),
		body.Email,
		body.Password,
		string(body.Verification),
	)
	if err != nil {
		return err
	}

	return c.JSON(tokenResponse)
}
