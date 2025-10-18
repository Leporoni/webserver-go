package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"webserver-go/internal/models"
	"webserver-go/internal/services"
)

// ProdutoHandler handles HTTP requests for produtos
type ProdutoHandler struct {
	service *services.ProdutoService
}

// NewProdutoHandler creates a new produto handler
func NewProdutoHandler() *ProdutoHandler {
	return &ProdutoHandler{
		service: services.NewProdutoService(),
	}
}

// formatFloatPointer é uma função auxiliar para o template que formata um *float64.
// Se o ponteiro for nil, retorna "N/A". Caso contrário, formata o número.
func formatFloatPointer(f *float64) string {
	if f == nil {
		return "N/A"
	}
	return fmt.Sprintf("%.2f", *f)
}

func (h *ProdutoHandler) ViewProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	produto, err := h.service.GetProdutoByID(id)
	if err != nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	// Adicionando a função auxiliar ao mapa de funções do template
	funcs := template.FuncMap{"formatFloat": formatFloatPointer}

	tmpl, err := template.New("view.html").Funcs(funcs).ParseFiles("web/template/layout.html", "web/template/produtos/view.html")
	if err != nil {
		log.Printf("Erro ao fazer parse do template: %v", err)
		http.Error(w, "Erro ao renderizar a página", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", produto)
	if err != nil {
		log.Printf("Erro ao executar o template: %v", err)
		http.Error(w, "Erro ao renderizar a página", http.StatusInternalServerError)
	}
}

// API Handlers

// ListProdutosAPI handles GET /api/produtos
func (h *ProdutoHandler) ListProdutosAPI(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	search := r.URL.Query().Get("search")
	categoriaID := r.URL.Query().Get("categoria_id")

	// Get produtos
	result, err := h.service.GetAllProdutos(page, limit, search, categoriaID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar produtos: %v", err), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Produtos encontrados com sucesso",
		"data":    result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProdutoAPI handles GET /api/produtos/{id}
func (h *ProdutoHandler) GetProdutoAPI(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/api/produtos/")
	if id == "" {
		http.Error(w, "ID do produto é obrigatório", http.StatusBadRequest)
		return
	}

	// Get produto
	produto, err := h.service.GetProdutoByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao buscar produto: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Produto encontrado com sucesso",
		"data":    produto,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateProdutoAPI handles POST /api/produtos
func (h *ProdutoHandler) CreateProdutoAPI(w http.ResponseWriter, r *http.Request) {
	var input models.ProdutoInput

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Create produto
	produto, err := h.service.CreateProduto(&input)
	if err != nil {
		if strings.Contains(err.Error(), "já existe") || strings.Contains(err.Error(), "obrigatório") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao criar produto: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Produto criado com sucesso",
		"data":    produto,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateProdutoAPI handles PUT /api/produtos/{id}
func (h *ProdutoHandler) UpdateProdutoAPI(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/api/produtos/")
	if id == "" {
		http.Error(w, "ID do produto é obrigatório", http.StatusBadRequest)
		return
	}

	var input models.ProdutoInput

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Update produto
	produto, err := h.service.UpdateProduto(id, &input)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if strings.Contains(err.Error(), "já existe") || strings.Contains(err.Error(), "obrigatório") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao atualizar produto: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Produto atualizado com sucesso",
		"data":    produto,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteProdutoAPI handles DELETE /api/produtos/{id}
func (h *ProdutoHandler) DeleteProdutoAPI(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id := strings.TrimPrefix(r.URL.Path, "/api/produtos/")
	if id == "" {
		http.Error(w, "ID do produto é obrigatório", http.StatusBadRequest)
		return
	}

	// Delete produto
	err := h.service.DeleteProduto(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao deletar produto: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Produto deletado com sucesso",
		"data":    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCategoriasAPI handles GET /api/categorias
func (h *ProdutoHandler) GetCategoriasAPI(w http.ResponseWriter, r *http.Request) {
	// Get categorias
	categorias, err := h.service.GetCategorias()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar categorias: %v", err), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Categorias encontradas com sucesso",
		"data":    categorias,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateCategoriaAPI handles POST /api/categorias
func (h *ProdutoHandler) CreateCategoriaAPI(w http.ResponseWriter, r *http.Request) {
	var input models.CategoriaInput

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Create categoria
	categoria, err := h.service.CreateCategoria(&input)
	if err != nil {
		if strings.Contains(err.Error(), "já existe") || strings.Contains(err.Error(), "obrigatório") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao criar categoria: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"success": true,
		"message": "Categoria criada com sucesso",
		"data":    categoria,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Web Handlers

// ListProdutosWeb handles GET /produtos
func (h *ProdutoHandler) ListProdutosWeb(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	search := r.URL.Query().Get("search")
	categoriaID := r.URL.Query().Get("categoria_id")

	// Get produtos
	result, err := h.service.GetAllProdutos(page, 20, search, categoriaID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar produtos: %v", err), http.StatusInternalServerError)
		return
	}

	// Get categorias for filter
	categorias, err := h.service.GetCategorias()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar categorias: %v", err), http.StatusInternalServerError)
		return
	}

	// Render template
	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Produtos - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
        .header h1 { color: #333; margin: 0; }
        .btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-success { background: #28a745; }
        .btn-success:hover { background: #1e7e34; }
        .btn-danger { background: #dc3545; }
        .btn-danger:hover { background: #c82333; }
        .btn-sm { padding: 5px 10px; font-size: 12px; }
        .search-form { display: flex; gap: 10px; margin-bottom: 20px; align-items: center; flex-wrap: wrap; }
        .search-form input, .search-form select { padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
        .search-form input[type="text"] { flex: 1; min-width: 200px; }
        .table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        .table th, .table td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        .table th { background-color: #f8f9fa; font-weight: bold; }
        .table tr:hover { background-color: #f8f9fa; }
        .pagination { display: flex; justify-content: center; gap: 5px; margin-top: 20px; }
        .pagination a, .pagination span { padding: 8px 12px; border: 1px solid #ddd; text-decoration: none; color: #007bff; }
        .pagination .current { background: #007bff; color: white; }
        .pagination a:hover { background: #e9ecef; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .status.ativo { background: #d4edda; color: #155724; }
        .status.inativo { background: #f8d7da; color: #721c24; }
        .price { font-weight: bold; color: #28a745; }
        .actions { display: flex; gap: 5px; }
        .breadcrumb { margin-bottom: 20px; }
        .breadcrumb a { color: #007bff; text-decoration: none; }
        .breadcrumb a:hover { text-decoration: underline; }
        .no-results { text-align: center; padding: 40px; color: #666; }
        .categoria-tag { background: #e9ecef; padding: 2px 6px; border-radius: 3px; font-size: 11px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="breadcrumb">
            <a href="/">Home</a> > <strong>Produtos</strong>
        </div>

        <div class="header">
            <h1>📦 Gestão de Produtos</h1>
            <a href="/produtos/novo" class="btn btn-success">+ Novo Produto</a>
        </div>

        <form class="search-form" method="GET">
            <input type="text" name="search" placeholder="Buscar por nome, descrição ou código..." value="{{.Search}}">
            <select name="categoria_id">
                <option value="">Todas as categorias</option>
                {{range .Categorias}}
                <option value="{{.ID}}" {{if eq $.CategoriaID .ID}}selected{{end}}>{{.Nome}}</option>
                {{end}}
            </select>
            <button type="submit" class="btn">🔍 Buscar</button>
            {{if or .Search .CategoriaID}}
            <a href="/produtos" class="btn" style="background: #6c757d;">Limpar</a>
            {{end}}
        </form>

        {{if .Produtos}}
        <table class="table">
            <thead>
                <tr>
                    <th>Nome</th>
                    <th>Categoria</th>
                    <th>Código</th>
                    <th>Preço Venda</th>
                    <th>Unidade</th>
                    <th>Status</th>
                    <th>Ações</th>
                </tr>
            </thead>
            <tbody>
                {{range .Produtos}}
                <tr>
                    <td>
                        <strong>{{.Nome}}</strong>
                        {{if .Descricao}}<br><small style="color: #666;">{{.Descricao}}</small>{{end}}
                    </td>
                    <td>
                        {{if .Categoria}}
                        <span class="categoria-tag">{{.Categoria.Nome}}</span>
                        {{else}}
                        <span style="color: #999;">Sem categoria</span>
                        {{end}}
                    </td>
                    <td>
                        {{if .CodigoBarras}}{{.CodigoBarras}}{{else}}<span style="color: #999;">-</span>{{end}}
                    </td>
                    <td class="price">R$ {{printf "%.2f" .PrecoVenda}}</td>
                    <td>{{.UnidadeMedida}}</td>
                    <td>
                        {{if .Ativo}}
                        <span class="status ativo">Ativo</span>
                        {{else}}
                        <span class="status inativo">Inativo</span>
                        {{end}}
                    </td>
                    <td class="actions">
                        <a href="/produtos/{{.ID}}" class="btn btn-sm">Ver</a>
                        <a href="/produtos/{{.ID}}/editar" class="btn btn-sm">Editar</a>
                        <button onclick="deleteProduto('{{.ID}}', '{{.Nome}}')" class="btn btn-danger btn-sm">Excluir</button>
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>

        <!-- Pagination -->
        {{if gt .TotalPages 1}}
        <div class="pagination">
            {{if gt .Page 1}}
            <a href="?page={{sub .Page 1}}{{if .Search}}&search={{.Search}}{{end}}{{if .CategoriaID}}&categoria_id={{.CategoriaID}}{{end}}">&laquo; Anterior</a>
            {{end}}
            
            {{range $i := .PageNumbers}}
            {{if eq $i $.Page}}
            <span class="current">{{$i}}</span>
            {{else}}
            <a href="?page={{$i}}{{if $.Search}}&search={{$.Search}}{{end}}{{if $.CategoriaID}}&categoria_id={{$.CategoriaID}}{{end}}">{{$i}}</a>
            {{end}}
            {{end}}
            
            {{if lt .Page .TotalPages}}
            <a href="?page={{add .Page 1}}{{if .Search}}&search={{.Search}}{{end}}{{if .CategoriaID}}&categoria_id={{.CategoriaID}}{{end}}">Próxima &raquo;</a>
            {{end}}
        </div>
        {{end}}

        <div style="margin-top: 30px; text-align: center; color: #666;">
            <p>Total: {{.Total}} produto(s) encontrado(s)</p>
        </div>
        {{else}}
        <div class="no-results">
            <h3>Nenhum produto encontrado</h3>
            {{if or .Search .CategoriaID}}
            <p>Tente ajustar os filtros de busca.</p>
            <a href="/produtos" class="btn">Ver todos os produtos</a>
            {{else}}
            <p>Comece criando seu primeiro produto.</p>
            <a href="/produtos/novo" class="btn btn-success">+ Criar Primeiro Produto</a>
            {{end}}
        </div>
        {{end}}
    </div>

    <script>
        function deleteProduto(id, nome) {
            if (confirm('Tem certeza que deseja excluir o produto "' + nome + '"?')) {
                fetch('/api/produtos/' + id, {
                    method: 'DELETE'
                })
                .then(response => {
                    if (response.ok) {
                        alert('Produto excluído com sucesso!');
                        location.reload();
                    } else {
                        return response.text().then(text => {
                            throw new Error(text);
                        });
                    }
                })
                .catch(error => {
                    alert('Erro ao excluir produto: ' + error.message);
                });
            }
        }
    </script>
</body>
</html>`

	// Template functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
	}

	// Generate page numbers for pagination
	var pageNumbers []int
	start := 1
	end := result.TotalPages
	if end > 10 {
		if result.Page <= 5 {
			end = 10
		} else if result.Page > result.TotalPages-5 {
			start = result.TotalPages - 9
		} else {
			start = result.Page - 4
			end = result.Page + 5
		}
	}
	for i := start; i <= end; i++ {
		pageNumbers = append(pageNumbers, i)
	}

	// Template data
	data := struct {
		Produtos     []models.Produto
		Categorias   []models.Categoria
		Total        int
		Page         int
		TotalPages   int
		PageNumbers  []int
		Search       string
		CategoriaID  string
	}{
		Produtos:    result.Data.([]models.Produto),
		Categorias:  categorias,
		Total:       result.Total,
		Page:        result.Page,
		TotalPages:  result.TotalPages,
		PageNumbers: pageNumbers,
		Search:      search,
		CategoriaID: categoriaID,
	}

	t, err := template.New("produtos").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, data)
}

// NewProdutoWeb handles GET /produtos/novo
func (h *ProdutoHandler) NewProdutoWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.createProdutoWeb(w, r)
		return
	}

	// Get categorias for select
	categorias, err := h.service.GetCategorias()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar categorias: %v", err), http.StatusInternalServerError)
		return
	}

	// Render form template
	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Novo Produto - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { margin-bottom: 30px; }
        .header h1 { color: #333; margin: 0; }
        .form-group { margin-bottom: 20px; }
        .form-group label { display: block; margin-bottom: 5px; font-weight: bold; color: #333; }
        .form-group input, .form-group select, .form-group textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        .form-group textarea { height: 80px; resize: vertical; }
        .form-row { display: flex; gap: 20px; }
        .form-row .form-group { flex: 1; }
        .btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-success { background: #28a745; }
        .btn-success:hover { background: #1e7e34; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #545b62; }
        .actions { display: flex; gap: 10px; margin-top: 30px; }
        .breadcrumb { margin-bottom: 20px; }
        .breadcrumb a { color: #007bff; text-decoration: none; }
        .breadcrumb a:hover { text-decoration: underline; }
        .required { color: #dc3545; }
        .help-text { font-size: 12px; color: #666; margin-top: 5px; }
        .checkbox-group { display: flex; align-items: center; gap: 10px; }
        .checkbox-group input[type="checkbox"] { width: auto; }
    </style>
</head>
<body>
    <div class="container">
        <div class="breadcrumb">
            <a href="/">Home</a> > <a href="/produtos">Produtos</a> > <strong>Novo Produto</strong>
        </div>

        <div class="header">
            <h1>📦 Novo Produto</h1>
        </div>

        <form method="POST">
            <div class="form-row">
                <div class="form-group">
                    <label for="nome">Nome <span class="required">*</span></label>
                    <input type="text" id="nome" name="nome" required maxlength="255">
                    <div class="help-text">Nome do produto (2-255 caracteres)</div>
                </div>
                <div class="form-group">
                    <label for="categoria_id">Categoria</label>
                    <select id="categoria_id" name="categoria_id">
                        <option value="">Selecione uma categoria</option>
                        {{range .Categorias}}
                        <option value="{{.ID}}">{{.Nome}}</option>
                        {{end}}
                    </select>
                </div>
            </div>

            <div class="form-group">
                <label for="descricao">Descrição</label>
                <textarea id="descricao" name="descricao" placeholder="Descrição detalhada do produto"></textarea>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="codigo_barras">Código de Barras</label>
                    <input type="text" id="codigo_barras" name="codigo_barras" maxlength="50">
                    <div class="help-text">Código de barras único (opcional)</div>
                </div>
                <div class="form-group">
                    <label for="unidade_medida">Unidade de Medida <span class="required">*</span></label>
                    <select id="unidade_medida" name="unidade_medida" required>
                        <option value="UN">Unidade (UN)</option>
                        <option value="KG">Quilograma (KG)</option>
                        <option value="G">Grama (G)</option>
                        <option value="L">Litro (L)</option>
                        <option value="ML">Mililitro (ML)</option>
                        <option value="M">Metro (M)</option>
                        <option value="CM">Centímetro (CM)</option>
                        <option value="M2">Metro Quadrado (M²)</option>
                        <option value="M3">Metro Cúbico (M³)</option>
                        <option value="CX">Caixa (CX)</option>
                        <option value="PC">Peça (PC)</option>
                        <option value="PAR">Par (PAR)</option>
                        <option value="DZ">Dúzia (DZ)</option>
                    </select>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="preco_custo">Preço de Custo (R$)</label>
                    <input type="number" id="preco_custo" name="preco_custo" step="0.01" min="0" placeholder="0.00">
                    <div class="help-text">Preço de aquisição/produção</div>
                </div>
                <div class="form-group">
                    <label for="preco_venda">Preço de Venda (R$) <span class="required">*</span></label>
                    <input type="number" id="preco_venda" name="preco_venda" step="0.01" min="0.01" required placeholder="0.00">
                    <div class="help-text">Preço de venda ao cliente</div>
                </div>
                <div class="form-group">
                    <label for="margem_lucro">Margem de Lucro (%)</label>
                    <input type="number" id="margem_lucro" name="margem_lucro" step="0.1" min="0" max="1000" placeholder="25.0">
                    <div class="help-text">Porcentagem de lucro - calculada automaticamente se não informada</div>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="peso">Peso (kg)</label>
                    <input type="number" id="peso" name="peso" step="0.001" min="0" placeholder="2.550">
                    <div class="help-text">Peso em quilos (ex: 2.550 = 2kg e 550g)</div>
                </div>
                <div class="form-group">
                    <label for="dimensoes">Dimensões</label>
                    <input type="text" id="dimensoes" name="dimensoes" placeholder="Ex: 10x20x30 cm" maxlength="50">
                </div>
            </div>

            <div class="form-group">
                <label for="imagem_url">URL da Imagem</label>
                <input type="url" id="imagem_url" name="imagem_url" placeholder="https://exemplo.com/imagem.jpg" maxlength="500">
                <div class="help-text">Link para imagem do produto</div>
            </div>

            <div class="form-group">
                <label for="observacoes">Observações</label>
                <textarea id="observacoes" name="observacoes" placeholder="Observações adicionais sobre o produto"></textarea>
            </div>

            <div class="form-group">
                <div class="checkbox-group">
                    <input type="checkbox" id="ativo" name="ativo" checked>
                    <label for="ativo">Produto ativo</label>
                </div>
                <div class="help-text">Produtos inativos não aparecem nas listagens principais</div>
            </div>

            <div class="actions">
                <button type="submit" class="btn btn-success">💾 Salvar Produto</button>
                <a href="/produtos" class="btn btn-secondary">❌ Cancelar</a>
            </div>
        </form>
    </div>

    <script>
        // Auto-calculate margin when prices change
        document.getElementById('preco_custo').addEventListener('input', calculateMargin);
        document.getElementById('preco_venda').addEventListener('input', calculateMargin);

        function calculateMargin() {
            const precoCusto = parseFloat(document.getElementById('preco_custo').value) || 0;
            const precoVenda = parseFloat(document.getElementById('preco_venda').value) || 0;
            
            if (precoCusto > 0 && precoVenda > precoCusto) {
                const margem = ((precoVenda - precoCusto) / precoCusto) * 100;
                document.getElementById('margem_lucro').value = margem.toFixed(2);
            }
        }
    </script>
</body>
</html>`

	data := struct {
		Categorias []models.Categoria
	}{
		Categorias: categorias,
	}

	t, err := template.New("novo-produto").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, data)
}

// createProdutoWeb handles POST /produtos/novo
func (h *ProdutoHandler) createProdutoWeb(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	// Build input from form
	input := &models.ProdutoInput{
		Nome:          strings.TrimSpace(r.FormValue("nome")),
		UnidadeMedida: strings.TrimSpace(r.FormValue("unidade_medida")),
	}

	// Optional fields
	if desc := strings.TrimSpace(r.FormValue("descricao")); desc != "" {
		input.Descricao = &desc
	}

	if codigo := strings.TrimSpace(r.FormValue("codigo_barras")); codigo != "" {
		input.CodigoBarras = &codigo
	}

	if catID := strings.TrimSpace(r.FormValue("categoria_id")); catID != "" {
		input.CategoriaID = &catID
	}

	if precoCustoStr := r.FormValue("preco_custo"); precoCustoStr != "" {
		if precoCusto, err := strconv.ParseFloat(precoCustoStr, 64); err == nil {
			input.PrecoCusto = &precoCusto
		}
	}

	if precoVendaStr := r.FormValue("preco_venda"); precoVendaStr != "" {
		if precoVenda, err := strconv.ParseFloat(precoVendaStr, 64); err == nil {
			input.PrecoVenda = precoVenda
		}
	}

	if margemStr := r.FormValue("margem_lucro"); margemStr != "" {
		if margem, err := strconv.ParseFloat(margemStr, 64); err == nil {
			input.MargemLucro = &margem
		}
	}

	if pesoStr := r.FormValue("peso"); pesoStr != "" {
		if peso, err := strconv.ParseFloat(pesoStr, 64); err == nil {
			input.Peso = &peso
		}
	}

	if dim := strings.TrimSpace(r.FormValue("dimensoes")); dim != "" {
		input.Dimensoes = &dim
	}

	if img := strings.TrimSpace(r.FormValue("imagem_url")); img != "" {
		input.ImagemURL = &img
	}

	if obs := strings.TrimSpace(r.FormValue("observacoes")); obs != "" {
		input.Observacoes = &obs
	}

	ativo := r.FormValue("ativo") == "on"
	input.Ativo = &ativo

	// Create produto
	produto, err := h.service.CreateProduto(input)
	if err != nil {
		// TODO: Show error in form instead of generic error page
		http.Error(w, fmt.Sprintf("Erro ao criar produto: %v", err), http.StatusBadRequest)
		return
	}

	// Redirect to produto details
	http.Redirect(w, r, fmt.Sprintf("/produtos/%s", produto.ID), http.StatusSeeOther)
}

// ShowProdutoWeb handles GET /produtos/{id}
func (h *ProdutoHandler) ShowProdutoWeb(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/produtos/")
	if strings.Contains(path, "/") {
		// This is an edit request, handle elsewhere
		return
	}

	id := path
	if id == "" {
		http.Error(w, "ID do produto é obrigatório", http.StatusBadRequest)
		return
	}

	// Get produto
	produto, err := h.service.GetProdutoByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao buscar produto: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Render template
	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Nome}} - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30px; }
        .header h1 { color: #333; margin: 0; }
        .btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-danger { background: #dc3545; }
        .btn-danger:hover { background: #c82333; }
        .actions { display: flex; gap: 10px; }
        .breadcrumb { margin-bottom: 20px; }
        .breadcrumb a { color: #007bff; text-decoration: none; }
        .breadcrumb a:hover { text-decoration: underline; }
        .info-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .info-card { background: #f8f9fa; padding: 20px; border-radius: 8px; border-left: 4px solid #007bff; }
        .info-card h3 { margin-top: 0; color: #333; }
        .info-item { margin-bottom: 10px; }
        .info-item strong { color: #333; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .status.ativo { background: #d4edda; color: #155724; }
        .status.inativo { background: #f8d7da; color: #721c24; }
        .price { font-size: 18px; font-weight: bold; color: #28a745; }
        .categoria-tag { background: #e9ecef; padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .image-preview { max-width: 200px; max-height: 200px; border-radius: 8px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="breadcrumb">
            <a href="/">Home</a> > <a href="/produtos">Produtos</a> > <strong>{{.Nome}}</strong>
        </div>

        <div class="header">
            <h1>📦 {{.Nome}}</h1>
            <div class="actions">
                <a href="/produtos/{{.ID}}/editar" class="btn">✏️ Editar</a>
                <button onclick="deleteProduto('{{.ID}}', '{{.Nome}}')" class="btn btn-danger">🗑️ Excluir</button>
            </div>
        </div>

        {{if .ImagemURL}}
        <img src="{{.ImagemURL}}" alt="{{.Nome}}" class="image-preview">
        {{end}}

        <div class="info-grid">
            <div class="info-card">
                <h3>📋 Informações Básicas</h3>
                <div class="info-item">
                    <strong>Nome:</strong> {{.Nome}}
                </div>
                {{if .Descricao}}
                <div class="info-item">
                    <strong>Descrição:</strong> {{.Descricao}}
                </div>
                {{end}}
                {{if .CodigoBarras}}
                <div class="info-item">
                    <strong>Código de Barras:</strong> {{.CodigoBarras}}
                </div>
                {{end}}
                <div class="info-item">
                    <strong>Categoria:</strong> 
                    {{if .Categoria}}
                    <span class="categoria-tag">{{.Categoria.Nome}}</span>
                    {{else}}
                    <span style="color: #999;">Sem categoria</span>
                    {{end}}
                </div>
                <div class="info-item">
                    <strong>Status:</strong> 
                    {{if .Ativo}}
                    <span class="status ativo">Ativo</span>
                    {{else}}
                    <span class="status inativo">Inativo</span>
                    {{end}}
                </div>
            </div>

            <div class="info-card">
                <h3>💰 Preços e Margem</h3>
                <div class="card mb-4">
            <div class="card-header">
                <i class="fas fa-dollar-sign"></i> Preços e Margem
            </div>
            <div class="card-body">
                <p class="card-text"><strong>Preço de Custo:</strong> R$ {{ .PrecoCusto | formatFloat }}</p>
                <p class="card-text"><strong>Preço de Venda:</strong> R$ {{ .PrecoVenda | formatFloat }}</p>
                <p class="card-text"><strong>Margem de Lucro:</strong> {{ .MargemLucro | formatFloat }}%</p>
                <p class="card-text"><strong>Unidade de Medida:</strong> {{ .UnidadeMedida }}</p>
            </div>
        </div>
    </div>
    <div class="col-md-6">
        <div class="card mb-4">
            <div class="card-header">
                <i class="fas fa-ruler-combined"></i> Especificações Físicas
            </div>
            <div class="card-body">
                <p class="card-text"><strong>Peso:</strong> {{ .Peso | formatFloat }} kg</p>
            </div>
        </div>
        <div class="card mb-4">
            <div class="card-header">
                <i class="fas fa-info-circle"></i> Informações do Sistema
            </div>
            <div class="card-body">
                <p class="card-text"><strong>Criado em:</strong> {{.CreatedAt.Format "02/01/2006 15:04"}}</p>
                <p class="card-text"><strong>Atualizado em:</strong> {{.UpdatedAt.Format "02/01/2006 15:04"}}</p>
                <p class="card-text"><strong>ID:</strong> <code>{{.ID}}</code></p>
            </div>
        </div>
    </div>
</div>

{{if .Observacoes}}
<div class="info-card">
    <h3>📝 Observações</h3>
    <p>{{.Observacoes}}</p>
</div>
{{end}}

<div style="margin-top: 30px; text-align: center;">
    <a href="/produtos" class="btn">← Voltar para Produtos</a>
</div>
</div>

<script>
    function deleteProduto(id, nome) {
        if (confirm('Tem certeza que deseja excluir o produto "' + nome + '"?')) {
            fetch('/api/produtos/' + id, {
                method: 'DELETE'
            })
            .then(response => {
                if (response.ok) {
                    alert('Produto excluído com sucesso!');
                    window.location.href = '/produtos';
                } else {
                    return response.text().then(text => {
                        throw new Error(text);
                    });
                }
            })
            .catch(error => {
                alert('Erro ao excluir produto: ' + error.message);
            });
        }
    }
</script>
</body>
</html>`

	// Registrar função auxiliar de formatação para uso nos pipes do template
	funcs := template.FuncMap{"formatFloat": formatFloatPointer}

	t, err := template.New("produto").Funcs(funcs).Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, produto)
}

// EditProdutoWeb handles GET /produtos/{id}/editar
func (h *ProdutoHandler) EditProdutoWeb(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	path := strings.TrimPrefix(r.URL.Path, "/produtos/")
	id := strings.TrimSuffix(path, "/editar")

	if r.Method == "POST" {
		h.updateProdutoWeb(w, r, id)
		return
	}

	// Get produto
	produto, err := h.service.GetProdutoByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Erro ao buscar produto: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Get categorias for select
	categorias, err := h.service.GetCategorias()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar categorias: %v", err), http.StatusInternalServerError)
		return
	}

	// Render edit form template (similar to new form but with values filled)
	// For brevity, I'll create a simplified version
	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Editar {{.Produto.Nome}} - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { margin-bottom: 30px; }
        .header h1 { color: #333; margin: 0; }
        .form-group { margin-bottom: 20px; }
        .form-group label { display: block; margin-bottom: 5px; font-weight: bold; color: #333; }
        .form-group input, .form-group select, .form-group textarea { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        .form-group textarea { height: 80px; resize: vertical; }
        .form-row { display: flex; gap: 20px; }
        .form-row .form-group { flex: 1; }
        .btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-success { background: #28a745; }
        .btn-success:hover { background: #1e7e34; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #545b62; }
        .actions { display: flex; gap: 10px; margin-top: 30px; }
        .breadcrumb { margin-bottom: 20px; }
        .breadcrumb a { color: #007bff; text-decoration: none; }
        .breadcrumb a:hover { text-decoration: underline; }
        .required { color: #dc3545; }
        .help-text { font-size: 12px; color: #666; margin-top: 5px; }
        .checkbox-group { display: flex; align-items: center; gap: 10px; }
        .checkbox-group input[type="checkbox"] { width: auto; }
    </style>
</head>
<body>
    <div class="container">
        <div class="breadcrumb">
            <a href="/">Home</a> > <a href="/produtos">Produtos</a> > <a href="/produtos/{{.Produto.ID}}">{{.Produto.Nome}}</a> > <strong>Editar</strong>
        </div>

        <div class="header">
            <h1>✏️ Editar Produto</h1>
        </div>

        <form method="POST">
            <div class="form-row">
                <div class="form-group">
                    <label for="nome">Nome <span class="required">*</span></label>
                    <input type="text" id="nome" name="nome" value="{{.Produto.Nome}}" required maxlength="255">
                </div>
                <div class="form-group">
                    <label for="categoria_id">Categoria</label>
                    <select id="categoria_id" name="categoria_id">
                        <option value="">Selecione uma categoria</option>
                        {{range .Categorias}}
                        <option value="{{.ID}}" {{if and $.Produto.CategoriaID (eq $.Produto.CategoriaID .ID)}}selected{{end}}>{{.Nome}}</option>
                        {{end}}
                    </select>
                </div>
            </div>

            <div class="form-group">
                <label for="descricao">Descrição</label>
                <textarea id="descricao" name="descricao">{{if .Produto.Descricao}}{{.Produto.Descricao}}{{end}}</textarea>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="codigo_barras">Código de Barras</label>
                    <input type="text" id="codigo_barras" name="codigo_barras" value="{{if .Produto.CodigoBarras}}{{.Produto.CodigoBarras}}{{end}}" maxlength="50">
                </div>
                <div class="form-group">
                    <label for="unidade_medida">Unidade de Medida <span class="required">*</span></label>
                    <select id="unidade_medida" name="unidade_medida" required>
                        <option value="UN" {{if eq .Produto.UnidadeMedida "UN"}}selected{{end}}>Unidade (UN)</option>
                        <option value="KG" {{if eq .Produto.UnidadeMedida "KG"}}selected{{end}}>Quilograma (KG)</option>
                        <option value="G" {{if eq .Produto.UnidadeMedida "G"}}selected{{end}}>Grama (G)</option>
                        <option value="L" {{if eq .Produto.UnidadeMedida "L"}}selected{{end}}>Litro (L)</option>
                        <option value="ML" {{if eq .Produto.UnidadeMedida "ML"}}selected{{end}}>Mililitro (ML)</option>
                        <option value="M" {{if eq .Produto.UnidadeMedida "M"}}selected{{end}}>Metro (M)</option>
                        <option value="CM" {{if eq .Produto.UnidadeMedida "CM"}}selected{{end}}>Centímetro (CM)</option>
                        <option value="M2" {{if eq .Produto.UnidadeMedida "M2"}}selected{{end}}>Metro Quadrado (M²)</option>
                        <option value="M3" {{if eq .Produto.UnidadeMedida "M3"}}selected{{end}}>Metro Cúbico (M³)</option>
                        <option value="CX" {{if eq .Produto.UnidadeMedida "CX"}}selected{{end}}>Caixa (CX)</option>
                        <option value="PC" {{if eq .Produto.UnidadeMedida "PC"}}selected{{end}}>Peça (PC)</option>
                        <option value="PAR" {{if eq .Produto.UnidadeMedida "PAR"}}selected{{end}}>Par (PAR)</option>
                        <option value="DZ" {{if eq .Produto.UnidadeMedida "DZ"}}selected{{end}}>Dúzia (DZ)</option>
                    </select>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="preco_custo">Preço de Custo (R$)</label>
                    <input type="number" id="preco_custo" name="preco_custo" step="0.01" min="0" value="{{if .Produto.PrecoCusto}}{{printf "%.2f" .Produto.PrecoCusto}}{{end}}">
                </div>
                <div class="form-group">
                    <label for="preco_venda">Preço de Venda (R$) <span class="required">*</span></label>
                    <input type="number" id="preco_venda" name="preco_venda" step="0.01" min="0.01" value="{{printf "%.2f" .Produto.PrecoVenda}}" required>
                </div>
                <div class="form-group">
                    <label for="margem_lucro">Margem de Lucro (%)</label>
                    <input type="number" id="margem_lucro" name="margem_lucro" step="0.1" min="0" max="1000" value="{{if .Produto.MargemLucro}}{{printf "%.1f" .Produto.MargemLucro}}{{end}}">
                    <div class="help-text">Porcentagem de lucro sobre o custo</div>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="peso">Peso (kg)</label>
                    <input type="number" id="peso" name="peso" step="0.001" min="0" value="{{if .Produto.Peso}}{{printf "%.3f" .Produto.Peso}}{{end}}">
                    <div class="help-text">Peso em quilos (ex: 2.550 = 2kg e 550g)</div>
                </div>
                <div class="form-group">
                    <label for="dimensoes">Dimensões</label>
                    <input type="text" id="dimensoes" name="dimensoes" value="{{if .Produto.Dimensoes}}{{.Produto.Dimensoes}}{{end}}" maxlength="50">
                </div>
            </div>

            <div class="form-group">
                <label for="imagem_url">URL da Imagem</label>
                <input type="url" id="imagem_url" name="imagem_url" value="{{if .Produto.ImagemURL}}{{.Produto.ImagemURL}}{{end}}" maxlength="500">
            </div>

            <div class="form-group">
                <label for="observacoes">Observações</label>
                <textarea id="observacoes" name="observacoes">{{if .Produto.Observacoes}}{{.Produto.Observacoes}}{{end}}</textarea>
            </div>

            <div class="form-group">
                <div class="checkbox-group">
                    <input type="checkbox" id="ativo" name="ativo" {{if .Produto.Ativo}}checked{{end}}>
                    <label for="ativo">Produto ativo</label>
                </div>
            </div>

            <div class="actions">
                <button type="submit" class="btn btn-success">💾 Salvar Alterações</button>
                <a href="/produtos/{{.Produto.ID}}" class="btn btn-secondary">❌ Cancelar</a>
            </div>
        </form>
    </div>

    <script>
        // Auto-calculate margin when prices change
        document.getElementById('preco_custo').addEventListener('input', calculateMargin);
        document.getElementById('preco_venda').addEventListener('input', calculateMargin);

        function calculateMargin() {
            const precoCusto = parseFloat(document.getElementById('preco_custo').value) || 0;
            const precoVenda = parseFloat(document.getElementById('preco_venda').value) || 0;
            
            if (precoCusto > 0 && precoVenda > precoCusto) {
                const margem = ((precoVenda - precoCusto) / precoCusto) * 100;
                document.getElementById('margem_lucro').value = margem.toFixed(2);
            }
        }
    </script>
</body>
</html>`

	data := struct {
		Produto    *models.Produto
		Categorias []models.Categoria
	}{
		Produto:    produto,
		Categorias: categorias,
	}

	t, err := template.New("editar-produto").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, data)
}

// updateProdutoWeb handles POST /produtos/{id}/editar
func (h *ProdutoHandler) updateProdutoWeb(w http.ResponseWriter, r *http.Request, id string) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	// Build input from form (similar to create)
	input := &models.ProdutoInput{
		Nome:          strings.TrimSpace(r.FormValue("nome")),
		UnidadeMedida: strings.TrimSpace(r.FormValue("unidade_medida")),
	}

	// Optional fields (same logic as create)
	if desc := strings.TrimSpace(r.FormValue("descricao")); desc != "" {
		input.Descricao = &desc
	}

	if codigo := strings.TrimSpace(r.FormValue("codigo_barras")); codigo != "" {
		input.CodigoBarras = &codigo
	}

	if catID := strings.TrimSpace(r.FormValue("categoria_id")); catID != "" {
		input.CategoriaID = &catID
	}

	if precoCustoStr := r.FormValue("preco_custo"); precoCustoStr != "" {
		if precoCusto, err := strconv.ParseFloat(precoCustoStr, 64); err == nil {
			input.PrecoCusto = &precoCusto
		}
	}

	if precoVendaStr := r.FormValue("preco_venda"); precoVendaStr != "" {
		if precoVenda, err := strconv.ParseFloat(precoVendaStr, 64); err == nil {
			input.PrecoVenda = precoVenda
		}
	}

	if margemStr := r.FormValue("margem_lucro"); margemStr != "" {
		if margem, err := strconv.ParseFloat(margemStr, 64); err == nil {
			input.MargemLucro = &margem
		}
	}

	if pesoStr := r.FormValue("peso"); pesoStr != "" {
		if peso, err := strconv.ParseFloat(pesoStr, 64); err == nil {
			input.Peso = &peso
		}
	}

	if dim := strings.TrimSpace(r.FormValue("dimensoes")); dim != "" {
		input.Dimensoes = &dim
	}

	if img := strings.TrimSpace(r.FormValue("imagem_url")); img != "" {
		input.ImagemURL = &img
	}

	if obs := strings.TrimSpace(r.FormValue("observacoes")); obs != "" {
		input.Observacoes = &obs
	}

	ativo := r.FormValue("ativo") == "on"
	input.Ativo = &ativo

	// Update produto
	_, err = h.service.UpdateProduto(id, input)
	if err != nil {
		// TODO: Show error in form instead of generic error page
		http.Error(w, fmt.Sprintf("Erro ao atualizar produto: %v", err), http.StatusBadRequest)
		return
	}

	// Redirect to produto details
	http.Redirect(w, r, fmt.Sprintf("/produtos/%s", id), http.StatusSeeOther)
}