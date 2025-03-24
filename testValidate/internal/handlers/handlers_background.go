package handlers

import (
	"net/http"
	"text/template"
)

func (handler *Handler) MainPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.Execute(w, nil)
}
func (handler *Handler) PersonPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {

	tmpl.Execute(w, nil)
}
func (handler *Handler) GreetingPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.Execute(w, nil)
}
