package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"webserver-go/internal/models"
	"webserver-go/internal/services"
)

type ProdutoHandler struct {
	produtoService   *services.ProdutoService
	categoriaService *services.CategoriaService
}

func NewProdutoHandler() *ProdutoHandler {
	return &ProdutoHandler{
		produtoService:   services.NewProdutoService(),
		categoriaService: services.NewCategoriaService(),
	}
}

// === API ENDPOINTS ===

// ListProdutosAPI lista produtos via API
func (h *ProdutoHandler) ListProdutosAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parâmetros de paginação
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

	// Buscar produtos
	produtos, total, err := h.produtoService.GetAllProdutos(page, limit, search, categoriaID)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: "Erro ao buscar produtos",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	paginatedResponse := PaginatedResponse{
		Data:       produtos,
		Total:      int(total),
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	response := APIResponse{
		Success: true,
		Message: "Produtos encontrados com sucesso",
		Data:    paginatedResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProdutoAPI busca produto por ID via API
func (h *ProdutoHandler) GetProdutoAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/api/produtos/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		response := APIResponse{
			Success: false,
			Message: "ID do produto é obrigatório",
			Error:   "ID não fornecido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	produto, err := h.produtoService.GetProdutoByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "não encontrado") {
			status = http.StatusNotFound
		}

		response := APIResponse{
			Success: false,
			Message: "Erro ao buscar produto",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Produto encontrado com sucesso",
		Data:    produto,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateProdutoAPI cria produto via API
func (h *ProdutoHandler) CreateProdutoAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var input models.ProdutoGormInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response := APIResponse{
			Success: false,
			Message: "Dados inválidos",
			Error:   "Erro ao decodificar JSON: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	produto, err := h.produtoService.CreateProduto(input)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: "Erro ao criar produto",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Produto criado com sucesso",
		Data:    produto,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateProdutoAPI atualiza produto via API
func (h *ProdutoHandler) UpdateProdutoAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/api/produtos/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		response := APIResponse{
			Success: false,
			Message: "ID do produto é obrigatório",
			Error:   "ID não fornecido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var input models.ProdutoGormInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response := APIResponse{
			Success: false,
			Message: "Dados inválidos",
			Error:   "Erro ao decodificar JSON: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	produto, err := h.produtoService.UpdateProduto(id, input)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "não encontrado") {
			status = http.StatusNotFound
		}

		response := APIResponse{
			Success: false,
			Message: "Erro ao atualizar produto",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Produto atualizado com sucesso",
		Data:    produto,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteProdutoAPI deleta produto via API
func (h *ProdutoHandler) DeleteProdutoAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/api/produtos/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		response := APIResponse{
			Success: false,
			Message: "ID do produto é obrigatório",
			Error:   "ID não fornecido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := h.produtoService.DeleteProduto(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "não encontrado") {
			status = http.StatusNotFound
		}

		response := APIResponse{
			Success: false,
			Message: "Erro ao deletar produto",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Produto deletado com sucesso",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// === WEB INTERFACE ===

// ListProdutosWeb exibe página de listagem de produtos
func (h *ProdutoHandler) ListProdutosWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parâmetros de paginação
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	search := r.URL.Query().Get("search")
	categoriaID := r.URL.Query().Get("categoria_id")

	// Buscar produtos
	produtos, total, err := h.produtoService.GetAllProdutos(page, 20, search, categoriaID)
	if err != nil {
		http.Error(w, "Erro ao buscar produtos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Buscar categorias para filtro
	categorias, err := h.categoriaService.GetAllCategorias(true)
	if err != nil {
		http.Error(w, "Erro ao buscar categorias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := int((total + 20 - 1) / 20)

	data := struct {
		Produtos     []models.ProdutoGorm
		Categorias   []models.CategoriaGorm
		Total        int64
		Page         int
		TotalPages   int
		Search       string
		CategoriaID  string
		HasPrev      bool
		HasNext      bool
	}{
		Produtos:    produtos,
		Categorias:  categorias,
		Total:       total,
		Page:        page,
		TotalPages:  totalPages,
		Search:      search,
		CategoriaID: categoriaID,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}

	tmpl := h.getProdutoListTemplate()

	// Funções auxiliares para template
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"formatPrice": func(price float64) string {
			return fmt.Sprintf("R$ %.2f", price)
		},
	}

	t, err := template.New("produtos").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}

// ShowProdutoWeb exibe detalhes de um produto
func (h *ProdutoHandler) ShowProdutoWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/produtos/")
	id := strings.Split(path, "/")[0]

	if id == "" || id == "novo" {
		http.Error(w, "ID do produto é obrigatório", http.StatusBadRequest)
		return
	}

	produto, err := h.produtoService.GetProdutoByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, "Produto não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar produto: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	tmpl := h.getProdutoDetailTemplate()

	funcMap := template.FuncMap{
		"formatPrice": func(price float64) string {
			return fmt.Sprintf("R$ %.2f", price)
		},
		"formatPricePtr": func(price *float64) string {
			if price == nil {
				return "-"
			}
			return fmt.Sprintf("R$ %.2f", *price)
		},
	}

	t, err := template.New("produto").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, produto); err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}

// NewProdutoWeb exibe formulário para novo produto
func (h *ProdutoHandler) NewProdutoWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.showProdutoForm(w, nil, "")
	} else if r.Method == http.MethodPost {
		h.createProdutoWeb(w, r)
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// EditProdutoWeb exibe formulário para editar produto
func (h *ProdutoHandler) EditProdutoWeb(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/produtos/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "editar" {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	id := parts[0]

	if r.Method == http.MethodGet {
		produto, err := h.produtoService.GetProdutoByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "não encontrado") {
				http.Error(w, "Produto não encontrado", http.StatusNotFound)
			} else {
				http.Error(w, "Erro ao buscar produto: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		h.showProdutoForm(w, produto, "")
	} else if r.Method == http.MethodPost {
		h.updateProdutoWeb(w, r, id)
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// createProdutoWeb processa criação via web
func (h *ProdutoHandler) createProdutoWeb(w http.ResponseWriter, r *http.Request) {
	input := h.parseProdutoForm(r)

	produto, err := h.produtoService.CreateProduto(input)
	if err != nil {
		h.showProdutoForm(w, nil, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/produtos/%s", produto.ID), http.StatusSeeOther)
}

// updateProdutoWeb processa atualização via web
func (h *ProdutoHandler) updateProdutoWeb(w http.ResponseWriter, r *http.Request, id string) {
	input := h.parseProdutoForm(r)

	produto, err := h.produtoService.UpdateProduto(id, input)
	if err != nil {
		// Buscar produto atual para mostrar no form
		currentProduto, _ := h.produtoService.GetProdutoByID(id)
		h.showProdutoForm(w, currentProduto, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/produtos/%s", produto.ID), http.StatusSeeOther)
}

// parseProdutoForm extrai dados do formulário
func (h *ProdutoHandler) parseProdutoForm(r *http.Request) models.ProdutoGormInput {
	r.ParseForm()

	input := models.ProdutoGormInput{
		Nome:          strings.TrimSpace(r.FormValue("nome")),
		UnidadeMedida: r.FormValue("unidade_medida"),
	}

	// Preço de venda (obrigatório)
	if precoVenda, err := strconv.ParseFloat(r.FormValue("preco_venda"), 64); err == nil {
		input.PrecoVenda = precoVenda
	}

	// Campos opcionais
	if descricao := strings.TrimSpace(r.FormValue("descricao")); descricao != "" {
		input.Descricao = &descricao
	}
	if codigoBarras := strings.TrimSpace(r.FormValue("codigo_barras")); codigoBarras != "" {
		input.CodigoBarras = &codigoBarras
	}
	if categoriaID := strings.TrimSpace(r.FormValue("categoria_id")); categoriaID != "" {
		input.CategoriaID = &categoriaID
	}
	if precoCusto, err := strconv.ParseFloat(r.FormValue("preco_custo"), 64); err == nil && precoCusto > 0 {
		input.PrecoCusto = &precoCusto
	}
	if margemLucro, err := strconv.ParseFloat(r.FormValue("margem_lucro"), 64); err == nil {
		input.MargemLucro = &margemLucro
	}
	if peso, err := strconv.ParseFloat(r.FormValue("peso"), 64); err == nil && peso > 0 {
		input.Peso = &peso
	}
	if dimensoes := strings.TrimSpace(r.FormValue("dimensoes")); dimensoes != "" {
		input.Dimensoes = &dimensoes
	}
	if imagemURL := strings.TrimSpace(r.FormValue("imagem_url")); imagemURL != "" {
		input.ImagemURL = &imagemURL
	}
	if observacoes := strings.TrimSpace(r.FormValue("observacoes")); observacoes != "" {
		input.Observacoes = &observacoes
	}

	// Checkbox ativo
	ativo := r.FormValue("ativo") == "on"
	input.Ativo = &ativo

	return input
}

// showProdutoForm exibe formulário de produto
func (h *ProdutoHandler) showProdutoForm(w http.ResponseWriter, produto *models.ProdutoGorm, errorMsg string) {
	// Buscar categorias para dropdown
	categorias, err := h.categoriaService.GetAllCategorias(true)
	if err != nil {
		http.Error(w, "Erro ao buscar categorias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	isEdit := produto != nil
	title := "Novo Produto"
	action := "/produtos/novo"
	if isEdit {
		title = "Editar Produto"
		action = fmt.Sprintf("/produtos/%s/editar", produto.ID)
	}

	data := struct {
		Title      string
		Action     string
		IsEdit     bool
		Produto    *models.ProdutoGorm
		Categorias []models.CategoriaGorm
		Error      string
	}{
		Title:      title,
		Action:     action,
		IsEdit:     isEdit,
		Produto:    produto,
		Categorias: categorias,
		Error:      errorMsg,
	}

	tmpl := h.getProdutoFormTemplate()

	funcMap := template.FuncMap{
		"formatFloat": func(f *float64) string {
			if f == nil {
				return ""
			}
			return fmt.Sprintf("%.2f", *f)
		},
	}

	t, err := template.New("form").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}