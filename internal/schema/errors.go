package schema

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ValidateFieldError is the error detail for validation errors
//
// see validator.FieldError for more information
type ValidateFieldError struct {
	// Tag is the validation tag that failed.
	// use alias if defined
	//
	// e.g. "required", "min", "max", etc.
	Tag string `json:"tag"`

	// Field is the field name that failed validation
	// use registered tag name if registered
	Field string `json:"field"`

	// Kind is the kind of the field type
	Kind reflect.Kind `json:"-"`

	// Param is the parameter for the validation
	Param string `json:"param"`

	// Value is the actual value that failed validation
	Value any `json:"value"`

	// Message is the error message
	Message string `json:"message"`
}

func (e *ValidateFieldError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	// construct error message
	// if you create a custom validation tag, you may need to switch case here
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
		e.Message = e.Field + "格式不正确"
	}

	return e.Message
}

// ValidationErrors is a list of ValidateFieldError
// for use in custom error messages post validation
//
// see validator.ValidationErrors for more information
type ValidationErrors []ValidateFieldError

func (e ValidationErrors) Error() string {
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

type HttpBaseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HttpBaseError) Error() string {
	return e.Message
}

type HttpError struct {
	HttpBaseError
	ValidationDetail ValidationErrors `json:"validation_detail,omitempty"`
}

func BadRequest(messages ...string) *HttpBaseError {
	message := "Bad Request"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &HttpBaseError{
		Code:    400,
		Message: message,
	}
}

func Unauthorized(messages ...string) *HttpBaseError {
	message := "Invalid JWT Token"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &HttpBaseError{
		Code:    401,
		Message: message,
	}
}

func Forbidden(messages ...string) *HttpBaseError {
	message := "Forbidden"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &HttpBaseError{
		Code:    403,
		Message: message,
	}
}

func NotFound(messages ...string) *HttpBaseError {
	message := "Not Found"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &HttpBaseError{
		Code:    404,
		Message: message,
	}
}

func InternalServerError(messages ...string) *HttpBaseError {
	message := "Internal Server Error"
	if len(messages) > 0 {
		message = messages[0]
	}
	return &HttpBaseError{
		Code:    500,
		Message: message,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var httpError HttpError

	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpError.Code = 404
	} else {
		switch e := err.(type) {
		case *HttpError:
			httpError = *e
		case *fiber.Error:
			httpError.Code = e.Code
		case ValidationErrors:
			httpError.Code = 400
			httpError.ValidationDetail = e
		case fiber.MultiError:
			httpError.Code = 400

			var stringBuilder strings.Builder
			for _, err := range e {
				stringBuilder.WriteString(err.Error())
				stringBuilder.WriteString("\n")
			}
			httpError.Message = stringBuilder.String()
		default:
			httpError.Code = 500
			httpError.Message = err.Error()
		}
	}

	// parse status code
	// when status code is 400xxx to 599xxx, use leading 3 numbers instead
	// else use 500
	statusCode := httpError.Code
	statusCodeString := strconv.Itoa(statusCode)
	if len(statusCodeString) > 3 {
		statusCodeString = statusCodeString[:3]
		newStatusCode, err := strconv.Atoi(statusCodeString)
		if err == nil && newStatusCode >= 400 && newStatusCode < 600 {
			statusCode = newStatusCode
		} else {
			statusCode = 500
		}
	}

	return c.Status(statusCode).JSON(&httpError)
}
