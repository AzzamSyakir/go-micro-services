package response

import "net/http"

type Response[T any] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
	Code    int    `json:"code,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func NewResponse[T any](w http.ResponseWriter, result *Response[T]) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(result.Code)
	encodeErr := json.NewEncoder(w).Encode(result)
	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
