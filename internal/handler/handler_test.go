package handler_test

import (
	"finverse/homework/internal/handler"
	"finverse/homework/internal/store"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type getTestCase struct {
	name    string
	request *http.Request

	// mock
	mockStoreGet []*store.Data
	mockStoreErr error

	// response to verify
	statusCode int
	body       string
}

func TestGetHandler(t *testing.T) {
	tc := []getTestCase{
		{
			name:    "GetSuccess",
			request: httptest.NewRequest("GET", "http://localhost/s/1/name/apple", nil),
			mockStoreGet: []*store.Data{
				{
					Id: 2,
					Content: store.Record{
						"name":  "Apple",
						"valid": true,
						"count": 123,
					},
				},
			},
			statusCode: http.StatusOK,
			body:       `{}`,
		},
	}

	mockStore := &mockStore{}
	handler := &handler.Handler{
		Store: mockStore,
	}

	r := chi.NewRouter()
	r.Get("/s/{specId}/{column}/{value}", handler.HandleGet)

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.HandleGet(w, tt.request)
			resp := w.Result()

			if tt.statusCode != resp.StatusCode {
				t.Error("invalid status code")
				t.FailNow()
			}

			body, _ := ioutil.ReadAll(resp.Body)
			if tt.body != string(body) {
				t.Error("invalid body")
				t.FailNow()
			}
		})
	}
}

type mockStore struct {
	handler.Store
}
