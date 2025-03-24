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

type PersonServiceInterface interface {
	Registration(newperk *Person, ctx context.Context) *AuthenticationResponse
	Authentication(newperk *Person, ctx context.Context) *AuthenticationResponse
}
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

type AuthenticationResponse struct {
	Success bool                   `json:"success"`
	UserId  uuid.UUID              `json:"userid,omitempty"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func (ps *PersonService) Registration(newperk *Person, ctx context.Context) *AuthenticationResponse {

	mapa := Validate(ps, newperk, true)
	if mapa != nil {
		return &AuthenticationResponse{Success: false, Errors: mapa}
	}

	errorsv := make(map[string]interface{})
	bytes, err := bcrypt.GenerateFromPassword([]byte(newperk.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		errorsv["Hash"] = erro.ErrorHash.Error()
		return &AuthenticationResponse{Success: false, Errors: errorsv}
	}
	newperk.Password = string(bytes)
	userID := uuid.New()
	newperk.Id = userID
	result := ps.Repo.Add(userID, newperk.Name, newperk.Email, newperk.Password, ctx)
	if !result.Success {
		if errors.Is(result.Error, erro.ErrorUniqueEmail) {
			log.Println(result.Error)
			errorsv["DbAdd"] = erro.ErrorUniqueEmail.Error()
		} else {
			log.Println(result.Error)
			errorsv["DbAdd"] = result.Error.Error()
		}
		return &AuthenticationResponse{Success: false, Errors: errorsv}
	}
	return &AuthenticationResponse{Success: true, UserId: userID}
}

func (ps *PersonService) Authentication(newperk *Person, ctx context.Context) *AuthenticationResponse {
	mapa := Validate(ps, newperk, false)

	if mapa != nil {
		return &AuthenticationResponse{Success: false, Errors: mapa}
	}

	result := ps.Repo.AuthenticateUser(newperk.Email, newperk.Password, ctx)

	if !result.Success {
		errorsv := make(map[string]interface{})
		if errors.Is(result.Error, erro.ErrorEmailNotRegister) {
			log.Println(result.Error)
			errorsv["Authority"] = erro.ErrorEmailNotRegister.Error()
		} else if errors.Is(result.Error, erro.ErrorInvalidPerson) {
			log.Println(result.Error)
			errorsv["Authority"] = erro.ErrorInvalidPerson.Error()
		} else {
			log.Println(result.Error)
			errorsv["Authority"] = result.Error.Error()
		}
		return &AuthenticationResponse{Success: false, Errors: errorsv}
	}

	return &AuthenticationResponse{Success: true, UserId: result.UserID}
}
