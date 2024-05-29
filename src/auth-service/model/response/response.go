package response

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func NewResponse[T any](w http.ResponseWriter, result *Response[T]) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(result.Code)
	encodeErr := json.NewEncoder(w).Encode(result)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
