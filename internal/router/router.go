package router

import (
	"ex_produce/internal/customer"
	"github.com/go-chi/chi/v5"
)

func InitRouter(service customer.Service) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/customer/get", service.GetUser)
	router.Post("/customer/save", service.SaveUser)
	return router
}
