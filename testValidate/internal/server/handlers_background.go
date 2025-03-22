package server

import (
	"net/http"
)

func (server *Server) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := server.MapaHtml["startpage"]
	tmpl.Execute(w, nil)
}
func (server *Server) PersonPage(w http.ResponseWriter, r *http.Request) {
	tmpl := server.MapaHtml["personpage"]
	tmpl.Execute(w, nil)
}
func (server *Server) GreetingPage(w http.ResponseWriter, r *http.Request) {
	tmpl := server.MapaHtml["greetingpage"]
	tmpl.Execute(w, nil)
}
