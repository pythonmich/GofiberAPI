package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	logs "github.com/sirupsen/logrus"
)
// errorStructResponse is our response struct
type errorStructResponse struct {
	failedField string
	tag string
	value string
}
// ValidateStruct will validate all requests as the field tags require
func ValidateStruct(arg interface{}) []*errorStructResponse {
	logs.WithField("func", "req_validator.go -> ValidateStruct()").Info()
	var errors []*errorStructResponse
	validate := validator.New()
	err := validate.Struct(arg);if err != nil{
		for _, err := range err.(validator.ValidationErrors) {
			var elements errorStructResponse
			elements.failedField = err.StructNamespace()
			elements.tag = err.Tag()
			elements.value = err.Param()
			errors = append(errors, &elements)
		}
	}
	return errors
}

func validateResponse(errs []*errorStructResponse) []fiber.Map {
	logs.WithField("func", "req_validator.go -> validateResponse()").Info()
	var resps []fiber.Map
	var resp fiber.Map
	for _, err := range errs {
		switch err.tag {
		case "required":
			resp = fiber.Map{
				"failed_field": err.failedField,
				"tag":          err.tag,
			}
		case "email":
			resp = fiber.Map{
				"failed_field": err.failedField,
				"tag":          err.tag,
				"value":        "invalid email",
			}
		default:
			resp = fiber.Map{
				"failed_field": err.failedField,
				"tag":          err.tag,
				"value":        err.value,
			}
		}

		resps = append(resps, resp)
	}
	return resps
}

