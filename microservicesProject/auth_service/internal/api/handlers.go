package api

import (
	"microservicesProject/auth_service/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	jsonResponseType = "application/json"
)

type Handler struct {
	services *service.Service
}
type HTTPResponse struct {
	Success bool                   `json:"success"`
	Errors  map[string]interface{} `json:"errors"`
	UserID  uuid.UUID              `json:"data"`
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *mux.Router {
	m := mux.NewRouter()
	m.HandleFunc("/reg", h.Registration).Methods("POST")
	m.HandleFunc("/auth", h.Authorization).Methods("POST")
	m.HandleFunc("/check-session", h.SessionChecker).Methods("GET")
	return m
}
