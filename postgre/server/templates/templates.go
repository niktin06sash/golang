package templates

import (
	"log"

	"text/template"
)

func MainPage() *template.Template {

	tmpl, err := template.ParseFiles("templates/mainp.html")
	if err != nil {
		log.Println(err)
	}
	return tmpl
}
func PersonInterface() *template.Template {
	tmpl, err := template.ParseFiles("templates/personinterface.html")
	if err != nil {
		log.Println(err)
	}
	return tmpl
}
