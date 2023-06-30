package common

import (
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
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

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}

func Validate(model any) error {
	errors := validate.Struct(model)
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

// ValidateQuery parse, set default and validate query
func ValidateQuery[T any](c *fiber.Ctx) (*T, error) {
	model := new(T)
	err := c.QueryParser(model)
	if err != nil {
		return nil, BadRequest(err.Error())
	}
	err = defaults.Set(model)
	if err != nil {
		return nil, err
	}

	// validate
	return model, Validate(model)
}

// ValidateBody parse, set default and validate body
// supports json only, if empty, using defaults
func ValidateBody[T any](c *fiber.Ctx) (*T, error) {
	body := c.Body()
	model := new(T)
	// empty request body, return default value
	if len(body) == 0 {
		return model, defaults.Set(model)
	}

	// unmarshal json
	err := json.Unmarshal(body, model)
	if err != nil {
		return nil, BadRequest(err.Error())
	}

	// set default value
	err = defaults.Set(model)
	if err != nil {
		return nil, err
	}

	// validate
	return model, Validate(model)
}
