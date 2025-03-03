package util

import (
	"encoding/json"
	"net/http"
)

// Структура JSON-ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}

// Функция для отправки JSON-ошибки
func SendErrorResponse(w http.ResponseWriter, code int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
}
