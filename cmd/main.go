package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/mrofisr/go-api/internal/handler"
	repository "github.com/mrofisr/go-api/internal/repository/postgres"
	"github.com/mrofisr/go-api/internal/router"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()
	pgxConn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%v:%v@localhost:5432/%v", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB")))
	if err != nil {
		log.Fatal(err)
	}
	// Close the connection when the main function exits
	defer pgxConn.Close(ctx)
	personRepo := repository.NewPersonRepository(pgxConn)
	// Create a new PersonHandler
	personHandler := handler.NewPersonHandler(personRepo)
	// Run the server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/person", router.PersonRouter(personHandler))
	// Create health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Run the server
	fmt.Println("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
