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
	// Se for uma requisição para API, retornar JSON
	if r.Header.Get("Accept") == "application/json" {
		response := Response{
			Message:   "Bem-vindo ao Sistema de Gestão!",
			Timestamp: time.Now(),
			Status:    "success",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Caso contrário, mostrar página HTML
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 40px; }
        .header h1 { color: #333; margin-bottom: 10px; }
        .header p { color: #666; font-size: 18px; }
        .modules { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-bottom: 40px; }
        .module { background: #f8f9fa; padding: 20px; border-radius: 8px; border-left: 4px solid #007bff; }
        .module h3 { margin-top: 0; color: #333; }
        .module p { color: #666; margin-bottom: 15px; }
        .btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; }
        .btn:hover { background: #0056b3; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 40px; }
        .stat { background: #e9ecef; padding: 20px; border-radius: 8px; text-align: center; }
        .stat h4 { margin: 0; color: #333; font-size: 24px; }
        .stat p { margin: 5px 0 0 0; color: #666; }
        .footer { text-align: center; color: #666; border-top: 1px solid #ddd; padding-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏢 Sistema de Gestão</h1>
            <p>Gerencie clientes, produtos e estoque de forma integrada</p>
        </div>

        <div class="stats">
            <div class="stat">
                <h4>✅</h4>
                <p>Sistema Online</p>
            </div>
            <div class="stat">
                <h4>🗄️</h4>
                <p>PostgreSQL</p>
            </div>
            <div class="stat">
                <h4>🐳</h4>
                <p>Docker Ready</p>
            </div>
            <div class="stat">
                <h4>🚀</h4>
                <p>Go + Nginx</p>
            </div>
        </div>

        <div class="modules">
            <div class="module">
                <h3>👥 Gestão de Clientes</h3>
                <p>Cadastre e gerencie informações completas dos seus clientes, incluindo dados pessoais, endereços e histórico.</p>
                <div style="display: flex; gap: 10px;">
                    <a href="/clientes" class="btn">Ver Clientes</a>
                    <a href="/clientes/novo" class="btn" style="background: #28a745;">+ Novo Cliente</a>
                </div>
            </div>

            <div class="module">
                <h3>📦 Gestão de Produtos</h3>
                <p>Controle seu catálogo de produtos com categorias, preços, especificações e imagens.</p>
                <div style="display: flex; gap: 10px;">
                    <a href="/produtos" class="btn">Ver Produtos</a>
                    <a href="/produtos/novo" class="btn" style="background: #28a745;">+ Novo Produto</a>
                </div>
            </div>

            <div class="module">
                <h3>📊 Controle de Estoque</h3>
                <p>Monitore entradas, saídas e movimentações de estoque com relatórios detalhados.</p>
                <a href="#" class="btn" style="background: #6c757d;">Em Desenvolvimento</a>
            </div>

            <div class="module">
                <h3>📈 Relatórios</h3>
                <p>Gere relatórios personalizados e acompanhe métricas importantes do seu negócio.</p>
                <a href="#" class="btn" style="background: #6c757d;">Em Desenvolvimento</a>
            </div>
        </div>

        <div class="footer">
            <p>Sistema de Gestão v1.0 - Desenvolvido com Go, PostgreSQL e Docker</p>
            <p><a href="/health" style="color: #007bff;">Status do Sistema</a> | <a href="/api/clientes" style="color: #007bff;">API Clientes</a> | <a href="/api/produtos" style="color: #007bff;">API Produtos</a></p>
        </div>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
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