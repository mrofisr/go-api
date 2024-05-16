package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mrofisr/go-api/internal/model"
	repository "github.com/mrofisr/go-api/internal/repository/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type PersonHandler struct {
	Repo repository.PersonRepository
}

func (ph *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("person-handler").Start(r.Context(), "GetPerson")
	span.SetAttributes(
		attribute.String("http.handler", "GetPerson"),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.path", r.URL.Path),
		attribute.String("http.host", r.Host),
		attribute.Int("http.status_code", http.StatusOK),
	)
	defer span.End()
	// Implementation
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
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("person-handler").Start(r.Context(), "GetPersonByID")
	span.SetAttributes(
		attribute.String("http.handler", "GetPersonByID"),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.path", r.URL.Path),
		attribute.String("http.host", r.Host),
		attribute.Int("http.status_code", http.StatusOK),
	)
	defer span.End()
	// Implementation
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
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("person-handler").Start(r.Context(), "CreatePerson")
	span.SetAttributes(
		attribute.String("http.handler", "CreatePerson"),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.path", r.URL.Path),
		attribute.String("http.host", r.Host),
		attribute.Int("http.status_code", http.StatusOK),
	)
	defer span.End()
	// Implementation
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
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("person-handler").Start(r.Context(), "UpdatePerson")
	span.SetAttributes(
		attribute.String("http.handler", "UpdatePerson"),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.path", r.URL.Path),
		attribute.String("http.host", r.Host),
		attribute.Int("http.status_code", http.StatusOK),
	)
	defer span.End()
	// Implementation
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
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("person-handler").Start(r.Context(), "DeletePerson")
	span.SetAttributes(
		attribute.String("http.handler", "DeletePerson"),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.path", r.URL.Path),
		attribute.String("http.host", r.Host),
		attribute.Int("http.status_code", http.StatusOK),
	)
	defer span.End()
	// Implementation
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

func (ph *PersonHandler) CountPerson(w http.ResponseWriter, r *http.Request) {
	tracer := otel.GetTracerProvider()
	ctx, span := tracer.Tracer("person-handler").Start(r.Context(), "CountPerson")
	span.SetAttributes(
		attribute.String("http.handler", "CountPerson"),
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.path", r.URL.Path),
		attribute.String("http.host", r.Host),
		attribute.Int("http.status_code", http.StatusOK),
	)
	defer span.End()
	// Implementation
	count, err := ph.Repo.Count(ctx)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"user_count": ` + strconv.Itoa(count) + `}`))
}

func NewPersonHandler(repo repository.PersonRepository) *PersonHandler {
	return &PersonHandler{Repo: repo}
}
