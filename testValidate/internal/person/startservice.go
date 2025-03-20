package person

import (
	"context"
	"log"
	"testValidate/internal/database"
	"testValidate/internal/erro"

	"github.com/go-playground/validator/v10"
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

	mapa := Validate(ps, newperk)
	if mapa == nil {
		bytes, err := bcrypt.GenerateFromPassword([]byte(newperk.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			errors := make(map[string]string)
			errors["Hash"] = erro.ErrorHash.Error()
			return errors
		}
		newperk.Password = string(bytes)
		flag, err := ps.Repo.CheckUniqueEmail(newperk.Email, ctx)

		if flag && err == nil {
			err = ps.Repo.Add(newperk.Name, newperk.Email, newperk.Password, ctx)
			if err != nil {
				log.Println(err)
				errors := make(map[string]string)
				errors["DBAdd"] = erro.ErrorDBAdd.Error()
				return errors
			}
			return nil
		}
		errors := make(map[string]string)
		errors["UniqueEmail"] = erro.ErrorUniqueEmail.Error()
		return errors
	}
	return mapa
}
func (ps *PersonService) Authentication(newperk *Person, ctx context.Context) map[string]string {

	flag, err := ps.Repo.CompareEmail(newperk.Email, ctx)
	if !flag || err != nil {
		errors := make(map[string]string)
		errors["UnregistedEmail"] = erro.ErrorEmailNotRegister.Error()
		return errors
	}
	flag, err = ps.Repo.ComparePassword(newperk.Email, newperk.Password, ctx)
	if !flag || err != nil {
		errors := make(map[string]string)
		errors["Password"] = erro.ErrorInvalidPerson.Error()
		return errors
	}
	return nil
}
