package person

import (
	"context"
	"fmt"
	"log"
	"testValidate/internal/database"
	"testValidate/internal/erro"

	"github.com/go-playground/validator/v10"
)

type PersonService struct {
	Repo      database.PersonRepository
	Validator *validator.Validate
}

func NewPersonService(repo database.PersonRepository, validator *validator.Validate) *PersonService {
	return &PersonService{
		Repo:      repo,
		Validator: validator,
	}
}
func (ps *PersonService) Registration(newperk *Person, ctx context.Context) map[string]string {
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
					errors[err.Field()] = fmt.Sprint(erro.ErrorMinName)

				default:
					errv := fmt.Sprintf("%s is Null", err.Field())
					errors[err.Field()] = fmt.Sprint(errv)

				}
			}
			return errors
		}
	}
	err = ps.Repo.Add(newperk.Name, newperk.Email, ctx)
	if err != nil {
		log.Println(err)
		errors := make(map[string]string)
		errors[""] = erro.ErrorDBAdd.Error()
		return errors
	}
	return nil

}
func (ps *PersonService) GetPerson(ctx context.Context) ([]Person, error) {
	rows, err := ps.Repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	var persondata []Person
	for rows.Next() {
		var current Person
		err := rows.Scan(&current.Id, &current.Name, &current.Email)
		if err != nil {
			return nil, err
		}
		persondata = append(persondata, current)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return persondata, nil
}
