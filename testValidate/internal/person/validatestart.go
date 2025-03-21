package person

import (
	"fmt"
	"testValidate/internal/erro"

	"github.com/go-playground/validator/v10"
)

func Validate(ps *PersonService, newperk *Person, flag bool) map[string]string {
	if !flag {
		newperk.Name = "qwertyuiopasdfghjklzxcvbn"
	}
	err := ps.Validator.Struct(newperk)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errors := make(map[string]string)
			for _, err := range validationErrors {

				switch err.Tag() {

				case "email":
					errors[err.Field()] = fmt.Sprint(erro.ErrorNotEmail)

				case "min":
					errv := fmt.Sprintf("%s is too short", err.Field())
					errors[err.Field()] = fmt.Sprint(errv)

				default:
					errv := fmt.Sprintf("%s is Null", err.Field())
					errors[err.Field()] = fmt.Sprint(errv)
				}
			}
			return errors
		}
	}

	return nil
}
