package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"webserver-go/internal/database"
	"webserver-go/internal/handlers"
	"webserver-go/internal/middleware"
)

func main() {
	// Conectar ao banco de dados
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Executar migrations
	if err := database.RunMigrations(database.DB); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// Configurar rotas
	http.HandleFunc("/", middleware.LoggingMiddleware(handlers.HomeHandler))
	http.HandleFunc("/health", middleware.LoggingMiddleware(handlers.HealthHandler))
	http.HandleFunc("/static/", middleware.LoggingMiddleware(handlers.StaticHandler))

	// Configurar servidor
	port := ":8080"
	fmt.Printf("Sistema de Gestão iniciando na porta %s\n", port)
	fmt.Println("Acesse: http://localhost:8080")
	fmt.Println("Health check: http://localhost:8080/health")
	fmt.Println("Database: Connected and migrations applied")

	// Configurar graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down gracefully...")
		database.Close()
		os.Exit(0)
	}()

	// Iniciar servidor
	log.Fatal(http.ListenAndServe(port, nil))
}