package api

import (
	"FiberFinanceAPI/utils"
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const (
	//{0} should be the field, and {1} should be the param if provided
	invalidCurrencyMSG = "{0} provided is currently not supported"
)

// validates is our request validate interface
type validates interface {
	validateRequests(req interface{}) []fiber.Map
	addTranslation(tag, errMsg string)
}

type validateRequest struct {
	validate   *validator.Validate
	translator ut.Translator
	logs       *utils.StandardLogger
}

func newValidator(logs *utils.StandardLogger) validates {
	v := &validateRequest{
		logs: logs,
	}
	v.validate = validator.New()
	// We register our custom validation before we validate our struct
	err := v.validate.RegisterValidation("currency", validCurrency)
	if err != nil {
		v.logs.WithError(err).Warn("could not register validation")
		return nil
	}
	return v
}

// validCurrency Register our validator for currency supported
var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(utils.CurrencyCode); ok {
		return utils.IsSupportedCurrency(currency)
	}
	return false
}

// validateRequests validates our struct requests
func (v *validateRequest) validateRequests(req interface{}) (errs []fiber.Map) {
	v.logs.WithField("func", "validate_req.go -> validateRequests()").Debug()
	err := v.validate.Struct(req)
	if err != nil {
		var found bool
		english := en.New()
		uni := ut.New(english, english)
		v.translator, found = uni.GetTranslator("en")
		if !found {
			err = errors.New("could not translate error")
			v.logs.Warn(err)
			return nil
		}
		v.addTranslation("currency", invalidCurrencyMSG)
		_ = enTranslation.RegisterDefaultTranslations(v.validate, v.translator)
		errs = v.translateError(err)
	}
	return errs
}

// translateError translates our errors into human readable forms
func (v *validateRequest) translateError(err error) (errs []fiber.Map) {
	v.logs.WithField("func", "validate_req.go -> translateError()").Debug()
	if err == nil {
		return nil
	}
	for _, e := range err.(validator.ValidationErrors) {
		v.logs.WithError(e).Warn(e.Translate(v.translator))
		mErr := fiber.Map{
			"failed": e.Translate(v.translator),
		}
		errs = append(errs, mErr)
	}
	return errs
}

// addTranslation adding a custom translation for our custom errors
func (v *validateRequest) addTranslation(tag, errMsg string) {
	v.logs.WithField("func", "validate_req.go -> addTranslation()").Debug()

	// register our message to our universal-translator
	registerFn := func(ut ut.Translator) error {
		v.logs.WithField("func", "validate_req.go -> translateError() -> registerFn()").Debug()
		return ut.Add(tag, errMsg, false)
	}
	// we are defining the tag we want to be translated
	translationFunc := func(ut ut.Translator, fe validator.FieldError) string {
		v.logs.WithField("func", "validate_req.go -> translateError() -> translationFunc()").Debug()
		param := fe.Param()
		tag = fe.Tag()
		v.logs.WithFields(logrus.Fields{"field": fe.Field(), "param": param, "tag": tag}).Debug("failed field")
		t, err := ut.T(tag, fe.Field(), param) // tag, {0}, {1}, ..
		if err != nil {
			v.logs.WithError(err).Warn(fe.(error).Error())
			return fe.(error).Error()
		}
		v.logs.Info("error message returned successfully")
		return t
	}
	// we use all of them to the built-in function of validator v10.
	_ = v.validate.RegisterTranslation(tag, v.translator, registerFn, translationFunc)

}
