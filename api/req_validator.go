package api

// Validate user requests
import (
	"FiberFinanceAPI/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

var (
	invalidEmail  = "invalid email"
	requiredField = "field is empty"
)

// errorStructResponse is our response struct
type errorStructResponse struct {
	failedField string
	tag         string
	value       string
}

// validateStruct will validate all requests as the field tags require
func validateStruct(arg interface{}, logs *utils.StandardLogger) []*errorStructResponse {
	logs.WithField("func", "req_validator.go -> validateStruct()").Debug("validating struct")
	var errors []*errorStructResponse
	validate := validator.New()
	err := validate.Struct(arg)
	if err != nil {
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

func validateResponse(errs []*errorStructResponse, logs *utils.StandardLogger) []fiber.Map {
	logs.WithField("func", "req_validator.go -> validateResponse()").Debug("validating request")

	var resps []fiber.Map
	var resp fiber.Map
	for _, err := range errs {
		// f is fields
		f := strings.Split(err.failedField, ".")
		// last element in the slice
		failedFieldName := f[len(f)-1]

		switch err.tag {
		case "required":
			logs.WithField(fmt.Sprintf("case %s", err.tag), fmt.Sprintf("failed field %s", err.failedField)).Warn()
			resp = validationErrorResponse(failedFieldName, err.tag, fmt.Sprint(requiredField))
		case "email":
			logs.WithField(fmt.Sprintf("case %s", err.tag), fmt.Sprintf("failed field %s", err.failedField)).Warn()
			resp = validationErrorResponse(failedFieldName, err.tag, invalidEmail)
		default:
			logs.WithField(fmt.Sprintf("case %s", err.tag), fmt.Sprintf("failed field %s", err.failedField)).Warn()
			resp = validationErrorResponse(failedFieldName, err.tag, err.value)
		}
		logs.WithField(fmt.Sprintf("case %s", err.tag), fmt.Sprintf("failed field %s", err.failedField)).Warn()
		resps = append(resps, resp)
	}
	return resps
}

func validationErrorResponse(failedField, tag, value string) fiber.Map {
	return fiber.Map{
		"failed_field": failedField,
		"tag":          tag,
		"value":        value,
	}
}
