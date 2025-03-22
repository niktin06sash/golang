package handlers

import (
	"net/http"
	"text/template"
)

func MainPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.Execute(w, nil)
}
func PersonPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {

	tmpl.Execute(w, nil)
}
func GreetingPage(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	tmpl.Execute(w, nil)
}
