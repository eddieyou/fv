package handler

import (
	"context"
	"encoding/json"
	"finverse/homework/internal/store"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	Store Store
}

type Store interface {
	Create(ctx context.Context, specId string, record store.Record) (*store.Data, error)
	Update(ctx context.Context, id int64, record store.Record) (*store.Data, error)
	Get(ctx context.Context, specId string, key string, value string) ([]*store.Data, error)
	Validate(ctx context.Context, specId string, record store.Record) error
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	specId := chi.URLParam(r, "specId")

	if r.Body == nil {
		log.Printf("empty body\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := store.Record{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error when reading body: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &record); err != nil {
		log.Printf("fail to unmarshal to json: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	log.Printf("got create request for spec %s :%+v\n", specId, record)

	ctx := context.Background()
	if err := h.Store.Validate(ctx, specId, record); err != nil {
		log.Printf("validation error: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.Store.Create(ctx, specId, record)
	if err != nil {
		log.Printf("fail to create record: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("created record for spec %s: %+v\n", specId, result)

	resultJson, err := json.Marshal(result)
	if err != nil {
		log.Printf("fail to marshal result to json: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resultJson)
	return
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	specId := chi.URLParam(r, "specId")
	column := chi.URLParam(r, "column")
	value := chi.URLParam(r, "value")

	if specId == "" {
		log.Printf("no specId defined\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if column == "" {
		log.Printf("no column defined\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if value == "" {
		log.Printf("no value defined\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("got query request for spec %s %s=%s\n", specId, column, value)

	ctx := context.Background()
	result, err := h.Store.Get(ctx, specId, column, value)
	if err != nil {
		log.Printf("fail to get record: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("got records %+v\n", result)

	resultJson, err := json.Marshal(result)
	if err != nil {
		log.Printf("fail to marshal result to json: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resultJson)
	return
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	specId := chi.URLParam(r, "specId")
	rowIdStr := chi.URLParam(r, "rowId")

	rowId, err := strconv.ParseInt(rowIdStr, 10, 64)
	if err != nil {
		log.Printf("error parsing rowId %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		log.Printf("empty body\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := store.Record{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error when reading body: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &record); err != nil {
		log.Printf("fail to unmarshal to json: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	log.Printf("got update request for spec %s rowId: %d :%+v\n", specId, rowId, record)

	ctx := context.Background()
	if err := h.Store.Validate(ctx, specId, record); err != nil {
		log.Printf("validation error: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.Store.Update(ctx, rowId, record)
	if err != nil {
		log.Printf("fail to update record: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("updated record: %+v\n", result)

	resultJson, err := json.Marshal(result)
	if err != nil {
		log.Printf("fail to marshal result to json: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resultJson)
	return
}
