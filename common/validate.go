package common

import (
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorDetailElement struct {
	Tag         string       `json:"tag"`
	Field       string       `json:"field"`
	Kind        reflect.Kind `json:"-"`
	Value       any          `json:"value"`
	Param       string       `json:"param"`
	StructField string       `json:"struct_field"`
	Message     string       `json:"message"`
}

func (e *ErrorDetailElement) Error() string {
	if e.Message != "" {
		return e.Message
	}

	switch e.Tag {
	case "min":
		if e.Kind == reflect.String {
			e.Message = e.Field + "至少" + e.Param + "字符"
		} else {
			e.Message = e.Field + "至少为" + e.Param
		}
	case "max":
		if e.Kind == reflect.String {
			e.Message = e.Field + "限长" + e.Param + "字符"
		} else {
			e.Message = e.Field + "至多为" + e.Param
		}
	case "required":
		e.Message = e.Field + "不能为空"
	case "email":
		e.Message = "邮箱格式不正确"
	default:
		e.Message = e.StructField + "格式不正确"
	}

	return e.Message
}

type ErrorDetail []*ErrorDetailElement

func (e ErrorDetail) Error() string {
	if len(e) == 0 {
		return "Validation Error"
	}

	if len(e) == 1 {
		return e[0].Error()
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString(e[0].Error())
	for _, err := range e[1:] {
		stringBuilder.WriteString(", ")
		stringBuilder.WriteString(err.Error())
	}
	return stringBuilder.String()
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
				Field:       err.Field(),
				Tag:         err.Tag(),
				Param:       err.Param(),
				Kind:        err.Kind(),
				Value:       err.Value(),
				StructField: err.StructField(),
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
