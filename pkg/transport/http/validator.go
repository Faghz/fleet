package inthttp

import (
	"strings"

	"github.com/elzestia/fleet/pkg/transport/http/response"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

var customValidator *CustomValidator

type CustomValidator struct {
	Validator  *validator.Validate
	translator ut.Translator
}

func CreateCustomValidator(validatorLang string) *CustomValidator {
	validator := validator.New()

	RegisterCustomValigator(validator)

	translator := setupTranslator(validatorLang, validator)

	customValidator = &CustomValidator{
		Validator:  validator,
		translator: translator,
	}

	return customValidator
}

func GetValidator() *CustomValidator {
	return customValidator
}

func setupTranslator(validatorLang string, validator *validator.Validate) ut.Translator {
	en := en.New()
	universalTranslator := ut.New(en, en)
	translator, _ := universalTranslator.GetTranslator(validatorLang)
	enTranslations.RegisterDefaultTranslations(validator, translator)

	return translator
}

func (v *CustomValidator) Validate(i interface{}) error {
	if err := v.Validator.Struct(i); err != nil {
		return v.ErrorValidator(err, i)
	}

	return nil
}

func (v *CustomValidator) ErrorValidator(err error, model interface{}) (newError error) {
	validationErrs, isOK := err.(validator.ValidationErrors)
	if !isOK {
		return err
	}
	fieldErrors := []response.FailureError{}

	for _, field := range validationErrs {
		key := strings.Join(strings.Split(field.Namespace(), ".")[1:], "/")

		message := strings.ToLower(field.Translate(v.translator))
		if customMessage, ok := CustomErrorMessages[field.Tag()]; ok {
			message = customMessage
		}

		fieldErrors = append(fieldErrors, response.FailureError{
			Pointer: strings.ToLower(key),
			Message: message,
		})
	}

	return response.GenerateBadRequest("Bad Request", "One or more field has invalid data", fieldErrors...)
}
