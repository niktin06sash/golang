package handlers

import (
	"net/http"
	"testValidate/internal/person"
	"text/template"
)

type Handler struct {
}
type HandlerInterface interface {
	Registration(w http.ResponseWriter, r *http.Request, ps person.PersonServiceInterface)
	Authentication(w http.ResponseWriter, r *http.Request, ps person.PersonServiceInterface)
	ProtectPersonPage(w http.ResponseWriter, r *http.Request)
	ProtectGreetingPage(w http.ResponseWriter, r *http.Request)
	ProtectMainPage(w http.ResponseWriter, r *http.Request)
	MainPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template)
	PersonPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template)
	GreetingPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
}

func NewHandler() *Handler {
	return &Handler{}
}
