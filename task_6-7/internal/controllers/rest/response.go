package rest

import (
	"encoding/json"
	"net/http"
)

type ResponseBase[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ResponseBaseMulti[T any] struct {
	Message string `json:"message"`
	Data    []T    `json:"data"`
}

type ResponseBaseErr struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func WriteResponse[Data any](w http.ResponseWriter, code int, data Data, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ResponseBase[Data]{
		Message: msg,
		Data:    data,
	}

	json.NewEncoder(w).Encode(resp)
}

func WriteResponseErr(w http.ResponseWriter, code int, err error, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err == nil {
		return
	}

	resp := ResponseBaseErr{
		Message: msg,
		Error:   err.Error(),
	}

	json.NewEncoder(w).Encode(resp)
}
