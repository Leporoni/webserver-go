package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// Response estrutura para resposta JSON
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// HomeHandler handler para a rota principal
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Bem-vindo ao Sistema de Gestão!",
		Timestamp: time.Now(),
		Status:    "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HealthHandler handler para a rota de saúde
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Sistema funcionando perfeitamente",
		Timestamp: time.Now(),
		Status:    "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StaticHandler handler para servir arquivos estáticos
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/"+r.URL.Path[8:]) // Remove "/static/" do path
}