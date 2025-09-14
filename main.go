package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Estrutura para resposta JSON
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// Handler para a rota principal
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Bem-vindo ao servidor web Go!",
		Timestamp: time.Now(),
		Status:    "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler para a rota de saúde
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Servidor funcionando perfeitamente",
		Timestamp: time.Now(),
		Status:    "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler para servir arquivos estáticos
func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/"+r.URL.Path[8:]) // Remove "/static/" do path
}

// Middleware para logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func main() {
	// Configurar rotas
	http.HandleFunc("/", loggingMiddleware(homeHandler))
	http.HandleFunc("/health", loggingMiddleware(healthHandler))
	http.HandleFunc("/static/", loggingMiddleware(staticHandler))

	// Configurar servidor
	port := ":8080"
	fmt.Printf("Servidor Go iniciando na porta %s\n", port)
	fmt.Println("Acesse: http://localhost:8080")
	fmt.Println("Health check: http://localhost:8080/health")

	// Iniciar servidor
	log.Fatal(http.ListenAndServe(port, nil))
}