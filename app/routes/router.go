package routes

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /ping", pingHandler)
	return router
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "pong",
	}
	json.NewEncoder(w).Encode(response)
}
