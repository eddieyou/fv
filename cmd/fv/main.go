package main

import (
	"context"
	"finverse/homework/internal/handler"
	"finverse/homework/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
)

func main() {
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("FV_DB"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	handler := &handler.Handler{
		Store: &store.Store{
			Pool: dbpool,
		},
	}

	// hardcode the creds for simplicity, can be loaded from somewhare
	creds := map[string]string{"fv": "fvpass"}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.BasicAuth("fv", creds))

	r.Get("/s/{specId}/{column}/{value}", handler.HandleGet)
	r.Post("/s/{specId}", handler.HandleCreate)
	r.Post("/s/{specId}/r/{rowId}", handler.HandleUpdate)
	http.ListenAndServe(":3000", r)
}
