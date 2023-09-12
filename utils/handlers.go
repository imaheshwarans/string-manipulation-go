package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stringinator-go/constants"
	"stringinator-go/models"
)

type endpointHandler func(w http.ResponseWriter, r *http.Request) error

type HandledError models.HandledError

func ErrorHandler(eh endpointHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Unknown Error", http.StatusInternalServerError)
			}
		}()
		if err := eh(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func JsonResponseHandler(h func(http.ResponseWriter, *http.Request) (interface{}, int, error)) endpointHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Header.Get("Accept") != constants.HTTPMediaTypeJson {
			return models.HandledError{
				StatusCode: http.StatusUnsupportedMediaType,
				Message:    "Invalid Accept type",
			}
		}

		data, status, err := h(w, r) // execute application handler
		if err != nil {
			return models.HandledError{
				StatusCode: status,
				Message:    err.Error(),
			}
		}
		w.Header().Set("Content-Type", constants.HTTPMediaTypeJson)
		w.WriteHeader(status)
		if data != nil {
			// Send JSON response back to the client application
			err = json.NewEncoder(w).Encode(data)
			if err != nil {
				fmt.Println("Failed to encode data in json")
			}
		}
		return nil
	}
}
