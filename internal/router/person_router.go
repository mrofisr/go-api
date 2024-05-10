package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrofisr/go-api/internal/handler"
)

func PersonRouter(personHandler *handler.PersonHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", personHandler.GetPerson)
	r.Get("/{id}", personHandler.GetPersonByID)
	r.Post("/", personHandler.CreatePerson)
	r.Put("/{id}", personHandler.UpdatePerson)
	r.Delete("/{id}", personHandler.DeletePerson)
	return r
}
