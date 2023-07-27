package utils

import (
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorDetailElement struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ErrorDetail []*ErrorDetailElement

func (e *ErrorDetail) Error() string {
	return "Validation Error"
}

var Validate = validator.New()

func init() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func ValidateStruct(model any) error {
	errors := Validate.Struct(model)
	if errors != nil {
		var errorDetail ErrorDetail
		for _, err := range errors.(validator.ValidationErrors) {
			detail := ErrorDetailElement{
				Field: err.Field(),
				Tag:   err.Tag(),
				Value: err.Param(),
			}
			errorDetail = append(errorDetail, &detail)
		}
		return &errorDetail
	}
	return nil
}

// ValidateQuery parse, set default and validate query into model
func ValidateQuery(c *fiber.Ctx, model any) error {
	// parse query into struct
	// see https://docs.gofiber.io/api/ctx/#queryparser
	err := c.QueryParser(model)
	if err != nil {
		return BadRequest(err.Error())
	}

	// set default value
	err = defaults.Set(model)
	if err != nil {
		return err
	}

	// Validate
	return ValidateStruct(model)
}

// ValidateBody parse, set default and validate body based on Content-Type.
// It supports json, xml and form only when struct tag exists; if empty, using defaults.
func ValidateBody(c *fiber.Ctx, model any) error {
	body := c.Body()

	// empty request body, return default value
	if len(body) == 0 {
		return defaults.Set(model)
	}

	// parse json, xml and form by fiber.BodyParser into struct
	// see https://docs.gofiber.io/api/ctx/#bodyparser
	err := c.BodyParser(model)
	if err != nil {
		return BadRequest(err.Error())
	}

	// set default value
	err = defaults.Set(model)
	if err != nil {
		return err
	}

	// Validate
	return ValidateStruct(model)
}
