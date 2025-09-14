package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	// Verificar status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler retornou status code errado: got %v want %v",
			status, http.StatusOK)
	}

	// Verificar Content-Type
	expected := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("handler retornou content-type errado: got %v want %v",
			contentType, expected)
	}

	// Verificar se é JSON válido
	var response Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Resposta não é JSON válido: %v", err)
	}

	// Verificar campos da resposta
	if response.Status != "success" {
		t.Errorf("Status esperado 'success', got %v", response.Status)
	}

	if response.Message == "" {
		t.Error("Message não pode estar vazio")
	}
}

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	// Verificar status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler retornou status code errado: got %v want %v",
			status, http.StatusOK)
	}

	// Verificar se é JSON válido
	var response Response
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Resposta não é JSON válido: %v", err)
	}

	// Verificar campos específicos do health check
	if response.Status != "healthy" {
		t.Errorf("Status esperado 'healthy', got %v", response.Status)
	}

	// Verificar se timestamp é recente (últimos 5 segundos)
	now := time.Now()
	diff := now.Sub(response.Timestamp)
	if diff > 5*time.Second {
		t.Errorf("Timestamp muito antigo: %v", diff)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Handler de teste
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := loggingMiddleware(testHandler)

	handler.ServeHTTP(rr, req)

	// Verificar se o handler original foi chamado
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("middleware não chamou handler original: got %v want %v",
			status, http.StatusOK)
	}

	if body := rr.Body.String(); body != "test" {
		t.Errorf("middleware alterou resposta: got %v want %v",
			body, "test")
	}
}

func TestResponseStruct(t *testing.T) {
	response := Response{
		Message:   "Test message",
		Timestamp: time.Now(),
		Status:    "test",
	}

	// Testar serialização JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Erro ao serializar Response: %v", err)
	}

	// Testar deserialização JSON
	var decoded Response
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("Erro ao deserializar Response: %v", err)
	}

	// Verificar se os dados foram preservados
	if decoded.Message != response.Message {
		t.Errorf("Message não preservado: got %v want %v",
			decoded.Message, response.Message)
	}

	if decoded.Status != response.Status {
		t.Errorf("Status não preservado: got %v want %v",
			decoded.Status, response.Status)
	}
}

// Benchmark para o homeHandler
func BenchmarkHomeHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/", nil)
	
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(homeHandler)
		handler.ServeHTTP(rr, req)
	}
}

// Benchmark para o healthHandler
func BenchmarkHealthHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/health", nil)
	
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(healthHandler)
		handler.ServeHTTP(rr, req)
	}
}