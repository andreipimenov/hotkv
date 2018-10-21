package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andreipimenov/hotkv/storage"
)

const (
	// Timeout respresents time to live for keys in storage.
	Timeout = 30 * time.Second
)

// APIKeyValue struct contains key and value.
type APIKeyValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value,omitempty"`
}

// APIResponse contains code and optional message for common API responses.
type APIResponse struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// WriteResponse prints API response into http.ResponseWriter.
func WriteResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	j, _ := json.Marshal(data)
	w.Write(j)
}

// NotFound for all unsupported API endpoints.
func NotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, APIResponse{
			Code:    http.StatusText(http.StatusNotFound),
			Message: "Unsupported API endpoint",
		}, http.StatusNotFound)
	})
}

// Ping - health check.
func Ping() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, APIResponse{
			Code:    http.StatusText(http.StatusOK),
			Message: "Pong",
		}, http.StatusOK)
	})
}

// GetKeys - get value by key.
func GetKeys(s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteResponse(w, APIResponse{
				Code:    http.StatusText(http.StatusMethodNotAllowed),
				Message: "HTTP method must be GET in order to get key",
			}, http.StatusMethodNotAllowed)
			return
		}
		uri := strings.Split(r.URL.Path, "/")
		if len(uri) != 4 || uri[3] == "" {
			WriteResponse(w, APIResponse{
				Code:    http.StatusText(http.StatusBadRequest),
				Message: "URI must match the pattern /api/keys/{key}",
			}, http.StatusBadRequest)
			return
		}
		key := uri[3]
		value, err := s.Get(context.Background(), key)
		if err != nil {
			WriteResponse(w, APIResponse{
				Code:    http.StatusText(http.StatusNotFound),
				Message: err.Error(),
			}, http.StatusNotFound)
			return
		}
		WriteResponse(w, APIKeyValue{
			Key:   key,
			Value: value,
		}, http.StatusOK)
	})
}

func SetKeys(s *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteResponse(w, APIResponse{
				Code:    http.StatusText(http.StatusMethodNotAllowed),
				Message: "HTTP method must be POST in order to set key",
			}, http.StatusMethodNotAllowed)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteResponse(w, APIResponse{
				Code:    http.StatusText(http.StatusBadRequest),
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		e := &APIKeyValue{}
		err = json.Unmarshal(body, e)
		if err != nil {
			WriteResponse(w, APIResponse{
				Code:    http.StatusText(http.StatusBadRequest),
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		ctx, _ := context.WithTimeout(context.Background(), Timeout)
		s.Set(ctx, e.Key, e.Value)
		WriteResponse(w, APIResponse{
			Code:    http.StatusText(http.StatusCreated),
			Message: fmt.Sprintf("Key %s is created successfully", e.Key),
		}, http.StatusCreated)
	})
}

func main() {
	s := storage.New()
	http.Handle("/", NotFound())
	http.Handle("/api/ping", Ping())
	http.Handle("/api/keys", SetKeys(s))
	http.Handle("/api/keys/", GetKeys(s))
	log.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
