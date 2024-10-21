package handlers

import (
	"github.com/DoktorGhost/golibrary-clients/internal/providers"
	"github.com/go-chi/chi"
)

func SetupRoutes(provider *providers.UseCaseProvider) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/login", handlerLogin(provider))
	r.Post("/register", handlerAddUser(provider))

	return r
}
