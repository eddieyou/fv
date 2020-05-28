package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"testing"
)

func TestGetSpec(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), "host=localhost port=5432 user=postgres password=fvpgsecret sslmode=disable")
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		t.Fail()
	}
	defer dbpool.Close()

	s := &Store{
		Pool: dbpool,
	}

	spec, err := s.GetSpec(context.Background(), "1")
	log.Printf("Got spec=%+v, err=%+v\n", spec, err)
}
