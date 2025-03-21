package person

import (
	"context"
	"errors"
	"log"
	"testValidate/internal/database"
	"testValidate/internal/erro"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	mapa := Validate(ps, newperk, true)
	if mapa == nil {
		bytes, err := bcrypt.GenerateFromPassword([]byte(newperk.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			errors := make(map[string]string)
			errors["Hash"] = erro.ErrorHash.Error()
			return errors
		}
		newperk.Password = string(bytes)
		userID := uuid.New()
		newperk.Id = userID
		err = ps.Repo.Add(userID, newperk.Name, newperk.Email, newperk.Password, ctx)
		if err != nil && errors.Is(err, erro.ErrorUniqueEmail) {
			log.Println(err)
			errors := make(map[string]string)
			errors["DbAdd"] = erro.ErrorUniqueEmail.Error()
			return errors
		} else if err != nil {
			log.Println(err)
			errors := make(map[string]string)
			errors["DbAdd"] = err.Error()
			return errors
		}
		return nil
	}
	return mapa
}
func (ps *PersonService) Authentication(newperk *Person, ctx context.Context) map[string]string {
	mapa := Validate(ps, newperk, false)

	if mapa == nil {
		flag, err := ps.Repo.AuthenticateUser(newperk.Email, newperk.Password, ctx)
		if !flag && errors.Is(err, erro.ErrorEmailNotRegister) {
			log.Println(err)
			errors := make(map[string]string)
			errors["Authority"] = erro.ErrorEmailNotRegister.Error()

			return errors
		} else if !flag && errors.Is(err, erro.ErrorInvalidPerson) {
			log.Println(err)
			errors := make(map[string]string)
			errors["Authority"] = erro.ErrorInvalidPerson.Error()

			return errors
		} else if !flag {
			log.Println(err)
			errors := make(map[string]string)
			errors["Authority"] = err.Error()

			return errors
		}
		return nil
	}
	return mapa
}
