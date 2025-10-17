package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	// Inicializar handlers
	clienteHandler := handlers.NewClienteHandler()
	produtoHandler := handlers.NewProdutoHandler()

	// Configurar rotas
	http.HandleFunc("/", middleware.LoggingMiddleware(handlers.HomeHandler))
	http.HandleFunc("/health", middleware.LoggingMiddleware(handlers.HealthHandler))
	http.HandleFunc("/static/", middleware.LoggingMiddleware(handlers.StaticHandler))

	// Rotas da API de clientes
	http.HandleFunc("/api/clientes", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			clienteHandler.ListClientesAPI(w, r)
		case http.MethodPost:
			clienteHandler.CreateClienteAPI(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/api/clientes/", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			clienteHandler.GetClienteAPI(w, r)
		case http.MethodPut:
			clienteHandler.UpdateClienteAPI(w, r)
		case http.MethodDelete:
			clienteHandler.DeleteClienteAPI(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}))

	// Rotas da interface web de clientes
	http.HandleFunc("/clientes", middleware.LoggingMiddleware(clienteHandler.ListClientesWeb))
	http.HandleFunc("/clientes/", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/clientes/novo" {
			clienteHandler.NewClienteWeb(w, r)
		} else if strings.HasSuffix(path, "/editar") {
			clienteHandler.EditClienteWeb(w, r)
		} else {
			clienteHandler.ShowClienteWeb(w, r)
		}
	}))

	// Rotas da API de produtos
	http.HandleFunc("/api/produtos", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			produtoHandler.ListProdutosAPI(w, r)
		case http.MethodPost:
			produtoHandler.CreateProdutoAPI(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/api/produtos/", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			produtoHandler.GetProdutoAPI(w, r)
		case http.MethodPut:
			produtoHandler.UpdateProdutoAPI(w, r)
		case http.MethodDelete:
			produtoHandler.DeleteProdutoAPI(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}))

	// Rotas da API de categorias
	http.HandleFunc("/api/categorias", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			produtoHandler.GetCategoriasAPI(w, r)
		case http.MethodPost:
			produtoHandler.CreateCategoriaAPI(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	}))

	// Rotas da interface web de produtos
	http.HandleFunc("/produtos", middleware.LoggingMiddleware(produtoHandler.ListProdutosWeb))
	http.HandleFunc("/produtos/", middleware.LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/produtos/novo" {
			produtoHandler.NewProdutoWeb(w, r)
		} else if strings.HasSuffix(path, "/editar") {
			produtoHandler.EditProdutoWeb(w, r)
		} else {
			produtoHandler.ShowProdutoWeb(w, r)
		}
	}))

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