package handler

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/pkg/log"
)

type RouteRegister interface {
	RegisterRoute(router fiber.Router)
}

type Handler struct {
	logger    *log.Logger
	validator *validator.Validate
}

func NewHandler(logger *log.Logger, validator *validator.Validate) *Handler {
	logger.Info("init handler")
	return &Handler{logger: logger, validator: validator}
}

func NewValidater() *validator.Validate {
	var v = validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	})

	return v
}

func (h *Handler) ValidateStruct(ctx context.Context, model any) error {
	err := h.validator.StructCtx(ctx, model)
	if err == nil {
		return nil
	}

	var rawValidationErrors *validator.ValidationErrors
	if ok := errors.As(err, &rawValidationErrors); ok {
		var validationErrors schema.ValidationErrors
		for _, fe := range *rawValidationErrors {
			validationErrors = append(validationErrors,
				schema.ValidateFieldError{
					Tag:   fe.Error(),
					Field: fe.Field(),
					Kind:  fe.Kind(),
					Param: fe.Param(),
					Value: fe.Value(),
				},
			)
		}
		return &validationErrors
	}

	var invalidValidationError *validator.InvalidValidationError
	if ok := errors.As(err, &invalidValidationError); ok {
		h.logger.Error("invalid validation error", zap.Error(err))
		return err
	}

	h.logger.Error("unknown validation error", zap.Error(err))
	return err
}

// ValidateQuery parse, set default and validate query into model
func (h *Handler) ValidateQuery(c *fiber.Ctx, model any) error {
	// parse query into struct
	// see https://docs.gofiber.io/api/ctx/#queryparser
	err := c.QueryParser(model)
	if err != nil {
		return schema.BadRequest(err.Error())
	}

	// set default value
	err = defaults.Set(model)
	if err != nil {
		return err
	}

	// Validate
	return h.ValidateStruct(c.Context(), model)
}

// ValidateBody parse, set default and validate body based on Content-Type.
// It supports json, xml and form only when struct tag exists; if empty, using defaults.
func (h *Handler) ValidateBody(c *fiber.Ctx, model any) error {
	body := c.Body()

	// empty request body, return default value
	if len(body) == 0 {
		return defaults.Set(model)
	}

	// parse json, xml and form by fiber.BodyParser into struct
	// see https://docs.gofiber.io/api/ctx/#bodyparser
	err := c.BodyParser(model)
	if err != nil {
		return schema.BadRequest(err.Error())
	}

	// set default value
	err = defaults.Set(model)
	if err != nil {
		return err
	}

	// Validate
	return h.ValidateStruct(c.Context(), model)
}
