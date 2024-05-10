package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mrofisr/go-api/internal/model"
	repository "github.com/mrofisr/go-api/internal/repository/postgres"
)

type PersonHandler struct {
	Repo repository.PersonRepository
}

func (ph *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	// Implementation
	ctx := r.Context()
	persons, err := ph.Repo.FindAll(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Implement your response here
	jsonPersons, err := json.Marshal(persons)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPersons)
	w.WriteHeader(http.StatusOK)
}

func (ph *PersonHandler) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	// Implementation
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	person, err := ph.Repo.FindById(ctx, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Implement your response here
	jsonPerson, err := json.Marshal(person)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPerson)
	w.WriteHeader(http.StatusOK)
}

func (ph *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	// Implementation
	ctx := r.Context()
	// Data from request body
	person := model.Person{}
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = ph.Repo.Create(ctx, person.Name, person.Age)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Person created"}`))
}

func (ph *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	// Implementation
	ctx := r.Context()
	newPerson := model.Person{}
	err := json.NewDecoder(r.Body).Decode(&newPerson)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = ph.Repo.Update(ctx, id, newPerson.Name, newPerson.Age)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Person updated", "id": ` + strconv.Itoa(id) + `}`))
}

func (ph *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Implementation
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = ph.Repo.Delete(ctx, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Person deleted", "id": ` + strconv.Itoa(id) + `}`))
}

func NewPersonHandler(repo repository.PersonRepository) *PersonHandler {
	return &PersonHandler{Repo: repo}
}
