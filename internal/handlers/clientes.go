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

type ClienteHandler struct {
	service *services.ClienteService
}

func NewClienteHandler() *ClienteHandler {
	return &ClienteHandler{
		service: services.NewClienteService(),
	}
}

// APIResponse estrutura para respostas da API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse estrutura para respostas paginadas
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// === API ENDPOINTS ===

// ListClientesAPI lista clientes via API
func (h *ClienteHandler) ListClientesAPI(w http.ResponseWriter, r *http.Request) {
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

	// Buscar clientes
	clientes, total, err := h.service.GetAllClientes(page, limit, search)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: "Erro ao buscar clientes",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	totalPages := (total + limit - 1) / limit

	paginatedResponse := PaginatedResponse{
		Data:       clientes,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	response := APIResponse{
		Success: true,
		Message: "Clientes encontrados com sucesso",
		Data:    paginatedResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetClienteAPI busca cliente por ID via API
func (h *ClienteHandler) GetClienteAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/api/clientes/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		response := APIResponse{
			Success: false,
			Message: "ID do cliente é obrigatório",
			Error:   "ID não fornecido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	cliente, err := h.service.GetClienteByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "não encontrado") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "inválido") {
			status = http.StatusBadRequest
		}

		response := APIResponse{
			Success: false,
			Message: "Erro ao buscar cliente",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Cliente encontrado com sucesso",
		Data:    cliente,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateClienteAPI cria cliente via API
func (h *ClienteHandler) CreateClienteAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var input models.ClienteInput
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

	cliente, err := h.service.CreateCliente(input)
	if err != nil {
		response := APIResponse{
			Success: false,
			Message: "Erro ao criar cliente",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Cliente criado com sucesso",
		Data:    cliente,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateClienteAPI atualiza cliente via API
func (h *ClienteHandler) UpdateClienteAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/api/clientes/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		response := APIResponse{
			Success: false,
			Message: "ID do cliente é obrigatório",
			Error:   "ID não fornecido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var input models.ClienteInput
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

	cliente, err := h.service.UpdateCliente(id, input)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "não encontrado") {
			status = http.StatusNotFound
		}

		response := APIResponse{
			Success: false,
			Message: "Erro ao atualizar cliente",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Cliente atualizado com sucesso",
		Data:    cliente,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteClienteAPI deleta cliente via API
func (h *ClienteHandler) DeleteClienteAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/api/clientes/")
	id := strings.Split(path, "/")[0]

	if id == "" {
		response := APIResponse{
			Success: false,
			Message: "ID do cliente é obrigatório",
			Error:   "ID não fornecido",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := h.service.DeleteCliente(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "não encontrado") {
			status = http.StatusNotFound
		}

		response := APIResponse{
			Success: false,
			Message: "Erro ao deletar cliente",
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := APIResponse{
		Success: true,
		Message: "Cliente deletado com sucesso",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// === WEB INTERFACE ===

// ListClientesWeb exibe página de listagem de clientes
func (h *ClienteHandler) ListClientesWeb(w http.ResponseWriter, r *http.Request) {
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

	// Buscar clientes
	clientes, total, err := h.service.GetAllClientes(page, 20, search)
	if err != nil {
		http.Error(w, "Erro ao buscar clientes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (total + 20 - 1) / 20

	data := struct {
		Clientes   []models.Cliente
		Total      int
		Page       int
		TotalPages int
		Search     string
		HasPrev    bool
		HasNext    bool
	}{
		Clientes:   clientes,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
		Search:     search,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
	}

	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Clientes - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
        .btn { padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-danger { background: #dc3545; }
        .btn-danger:hover { background: #c82333; }
        .btn-sm { padding: 5px 10px; font-size: 12px; }
        .search-box { margin-bottom: 20px; }
        .search-box input { padding: 8px; width: 300px; border: 1px solid #ddd; border-radius: 4px; }
        table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #f8f9fa; font-weight: bold; }
        tr:hover { background-color: #f5f5f5; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 12px; }
        .status.ativo { background: #d4edda; color: #155724; }
        .status.inativo { background: #f8d7da; color: #721c24; }
        .pagination { display: flex; justify-content: center; gap: 10px; }
        .pagination a { padding: 8px 12px; text-decoration: none; border: 1px solid #ddd; border-radius: 4px; }
        .pagination a.current { background: #007bff; color: white; }
        .actions { display: flex; gap: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Gestão de Clientes</h1>
            <a href="/clientes/novo" class="btn">+ Novo Cliente</a>
        </div>

        <div class="search-box">
            <form method="GET">
                <input type="text" name="search" placeholder="Buscar por nome ou email..." value="{{.Search}}">
                <button type="submit" class="btn">Buscar</button>
                {{if .Search}}<a href="/clientes" class="btn" style="background: #6c757d;">Limpar</a>{{end}}
            </form>
        </div>

        <p>Total: {{.Total}} cliente(s)</p>

        <table>
            <thead>
                <tr>
                    <th>Nome</th>
                    <th>Email</th>
                    <th>Telefone</th>
                    <th>Tipo</th>
                    <th>Status</th>
                    <th>Ações</th>
                </tr>
            </thead>
            <tbody>
                {{range .Clientes}}
                <tr>
                    <td>{{.Nome}}</td>
                    <td>{{if .Email}}{{.Email}}{{else}}-{{end}}</td>
                    <td>{{if .Telefone}}{{.Telefone}}{{else}}-{{end}}</td>
                    <td>{{if eq .TipoPessoa "fisica"}}Física{{else}}Jurídica{{end}}</td>
                    <td>
                        <span class="status {{if .Ativo}}ativo{{else}}inativo{{end}}">
                            {{if .Ativo}}Ativo{{else}}Inativo{{end}}
                        </span>
                    </td>
                    <td class="actions">
                        <a href="/clientes/{{.ID}}" class="btn btn-sm">Ver</a>
                        <a href="/clientes/{{.ID}}/editar" class="btn btn-sm">Editar</a>
                        <button onclick="deleteCliente('{{.ID}}', '{{.Nome}}')" class="btn btn-sm btn-danger">Excluir</button>
                    </td>
                </tr>
                {{else}}
                <tr>
                    <td colspan="6" style="text-align: center; padding: 40px;">
                        {{if .Search}}
                            Nenhum cliente encontrado para "{{.Search}}"
                        {{else}}
                            Nenhum cliente cadastrado
                        {{end}}
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>

        {{if gt .TotalPages 1}}
        <div class="pagination">
            {{if .HasPrev}}
                <a href="?page={{sub .Page 1}}{{if .Search}}&search={{.Search}}{{end}}">&laquo; Anterior</a>
            {{end}}
            
            {{range $i := .TotalPages}}
                {{$page := add $i 1}}
                <a href="?page={{$page}}{{if $.Search}}&search={{$.Search}}{{end}}" 
                   {{if eq $page $.Page}}class="current"{{end}}>{{$page}}</a>
            {{end}}
            
            {{if .HasNext}}
                <a href="?page={{add .Page 1}}{{if .Search}}&search={{.Search}}{{end}}">Próximo &raquo;</a>
            {{end}}
        </div>
        {{end}}
    </div>

    <script>
        function deleteCliente(id, nome) {
            if (confirm('Tem certeza que deseja excluir o cliente "' + nome + '"?')) {
                fetch('/api/clientes/' + id, {
                    method: 'DELETE'
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        alert('Cliente excluído com sucesso!');
                        location.reload();
                    } else {
                        alert('Erro ao excluir cliente: ' + data.error);
                    }
                })
                .catch(error => {
                    alert('Erro ao excluir cliente: ' + error);
                });
            }
        }
    </script>
</body>
</html>`

	// Funções auxiliares para template
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}

	t, err := template.New("clientes").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}

// ShowClienteWeb exibe detalhes de um cliente
func (h *ClienteHandler) ShowClienteWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/clientes/")
	id := strings.Split(path, "/")[0]

	if id == "" || id == "novo" {
		http.Error(w, "ID do cliente é obrigatório", http.StatusBadRequest)
		return
	}

	cliente, err := h.service.GetClienteByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			http.Error(w, "Cliente não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar cliente: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Nome}} - Cliente</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
        .btn { padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; }
        .btn:hover { background: #0056b3; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #545b62; }
        .info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
        .info-item { margin-bottom: 15px; }
        .info-label { font-weight: bold; color: #666; margin-bottom: 5px; }
        .info-value { font-size: 16px; }
        .status { padding: 4px 8px; border-radius: 4px; font-size: 14px; }
        .status.ativo { background: #d4edda; color: #155724; }
        .status.inativo { background: #f8d7da; color: #721c24; }
        .actions { display: flex; gap: 10px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Nome}}</h1>
            <span class="status {{if .Ativo}}ativo{{else}}inativo{{end}}">
                {{if .Ativo}}Ativo{{else}}Inativo{{end}}
            </span>
        </div>

        <div class="info-grid">
            <div>
                <div class="info-item">
                    <div class="info-label">Email:</div>
                    <div class="info-value">{{if .Email}}{{.Email}}{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Telefone:</div>
                    <div class="info-value">{{if .Telefone}}{{.Telefone}}{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">CPF/CNPJ:</div>
                    <div class="info-value">{{if .CPFCNPJ}}{{.CPFCNPJ}}{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Tipo de Pessoa:</div>
                    <div class="info-value">{{if eq .TipoPessoa "fisica"}}Física{{else}}Jurídica{{end}}</div>
                </div>
            </div>
            <div>
                <div class="info-item">
                    <div class="info-label">Endereço:</div>
                    <div class="info-value">{{if .Endereco}}{{.Endereco}}{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Cidade:</div>
                    <div class="info-value">{{if .Cidade}}{{.Cidade}}{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">Estado:</div>
                    <div class="info-value">{{if .Estado}}{{.Estado}}{{else}}-{{end}}</div>
                </div>
                <div class="info-item">
                    <div class="info-label">CEP:</div>
                    <div class="info-value">{{if .CEP}}{{.CEP}}{{else}}-{{end}}</div>
                </div>
            </div>
        </div>

        {{if .Observacoes}}
        <div class="info-item">
            <div class="info-label">Observações:</div>
            <div class="info-value">{{.Observacoes}}</div>
        </div>
        {{end}}

        <div class="info-item">
            <div class="info-label">Cadastrado em:</div>
            <div class="info-value">{{.CreatedAt.Format "02/01/2006 15:04"}}</div>
        </div>

        <div class="actions">
            <a href="/clientes" class="btn btn-secondary">← Voltar</a>
            <a href="/clientes/{{.ID}}/editar" class="btn">Editar</a>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("cliente").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, cliente); err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}

// NewClienteWeb exibe formulário para novo cliente
func (h *ClienteHandler) NewClienteWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.showClienteForm(w, nil, "")
	} else if r.Method == http.MethodPost {
		h.createClienteWeb(w, r)
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// EditClienteWeb exibe formulário para editar cliente
func (h *ClienteHandler) EditClienteWeb(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	path := strings.TrimPrefix(r.URL.Path, "/clientes/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "editar" {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	id := parts[0]

	if r.Method == http.MethodGet {
		cliente, err := h.service.GetClienteByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "não encontrado") {
				http.Error(w, "Cliente não encontrado", http.StatusNotFound)
			} else {
				http.Error(w, "Erro ao buscar cliente: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
		h.showClienteForm(w, cliente, "")
	} else if r.Method == http.MethodPost {
		h.updateClienteWeb(w, r, id)
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// createClienteWeb processa criação via web
func (h *ClienteHandler) createClienteWeb(w http.ResponseWriter, r *http.Request) {
	input := h.parseClienteForm(r)

	cliente, err := h.service.CreateCliente(input)
	if err != nil {
		h.showClienteForm(w, nil, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/clientes/%s", cliente.ID), http.StatusSeeOther)
}

// updateClienteWeb processa atualização via web
func (h *ClienteHandler) updateClienteWeb(w http.ResponseWriter, r *http.Request, id string) {
	input := h.parseClienteForm(r)

	cliente, err := h.service.UpdateCliente(id, input)
	if err != nil {
		// Buscar cliente atual para mostrar no form
		currentCliente, _ := h.service.GetClienteByID(id)
		h.showClienteForm(w, currentCliente, err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/clientes/%s", cliente.ID), http.StatusSeeOther)
}

// parseClienteForm extrai dados do formulário
func (h *ClienteHandler) parseClienteForm(r *http.Request) models.ClienteInput {
	r.ParseForm()

	input := models.ClienteInput{
		Nome:       strings.TrimSpace(r.FormValue("nome")),
		TipoPessoa: r.FormValue("tipo_pessoa"),
	}

	// Campos opcionais
	if email := strings.TrimSpace(r.FormValue("email")); email != "" {
		input.Email = &email
	}
	if telefone := strings.TrimSpace(r.FormValue("telefone")); telefone != "" {
		input.Telefone = &telefone
	}
	if endereco := strings.TrimSpace(r.FormValue("endereco")); endereco != "" {
		input.Endereco = &endereco
	}
	if cidade := strings.TrimSpace(r.FormValue("cidade")); cidade != "" {
		input.Cidade = &cidade
	}
	if estado := strings.TrimSpace(r.FormValue("estado")); estado != "" {
		input.Estado = &estado
	}
	if cep := strings.TrimSpace(r.FormValue("cep")); cep != "" {
		input.CEP = &cep
	}
	if cpfcnpj := strings.TrimSpace(r.FormValue("cpf_cnpj")); cpfcnpj != "" {
		input.CPFCNPJ = &cpfcnpj
	}
	if observacoes := strings.TrimSpace(r.FormValue("observacoes")); observacoes != "" {
		input.Observacoes = &observacoes
	}

	// Checkbox ativo
	ativo := r.FormValue("ativo") == "on"
	input.Ativo = &ativo

	return input
}

// showClienteForm exibe formulário de cliente
func (h *ClienteHandler) showClienteForm(w http.ResponseWriter, cliente *models.Cliente, errorMsg string) {
	isEdit := cliente != nil
	title := "Novo Cliente"
	action := "/clientes/novo"
	if isEdit {
		title = "Editar Cliente"
		action = fmt.Sprintf("/clientes/%s/editar", cliente.ID)
	}

	data := struct {
		Title    string
		Action   string
		IsEdit   bool
		Cliente  *models.Cliente
		Error    string
	}{
		Title:   title,
		Action:  action,
		IsEdit:  isEdit,
		Cliente: cliente,
		Error:   errorMsg,
	}

	tmpl := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Sistema de Gestão</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { margin-bottom: 20px; }
        .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        .form-group { margin-bottom: 15px; }
        .form-group.full-width { grid-column: 1 / -1; }
        label { display: block; margin-bottom: 5px; font-weight: bold; color: #333; }
        input, select, textarea { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; font-size: 14px; }
        textarea { height: 80px; resize: vertical; }
        .checkbox-group { display: flex; align-items: center; gap: 10px; }
        .checkbox-group input { width: auto; }
        .btn { padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; border: none; cursor: pointer; }
        .btn:hover { background: #0056b3; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #545b62; }
        .actions { display: flex; gap: 10px; margin-top: 20px; }
        .error { background: #f8d7da; color: #721c24; padding: 10px; border-radius: 4px; margin-bottom: 20px; }
        .required { color: red; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
        </div>

        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}

        <form method="POST" action="{{.Action}}">
            <div class="form-grid">
                <div class="form-group">
                    <label for="nome">Nome <span class="required">*</span></label>
                    <input type="text" id="nome" name="nome" required 
                           value="{{if .Cliente}}{{.Cliente.Nome}}{{end}}">
                </div>

                <div class="form-group">
                    <label for="tipo_pessoa">Tipo de Pessoa <span class="required">*</span></label>
                    <select id="tipo_pessoa" name="tipo_pessoa" required>
                        <option value="fisica" {{if and .Cliente (eq .Cliente.TipoPessoa "fisica")}}selected{{end}}>Física</option>
                        <option value="juridica" {{if and .Cliente (eq .Cliente.TipoPessoa "juridica")}}selected{{end}}>Jurídica</option>
                    </select>
                </div>

                <div class="form-group">
                    <label for="email">Email</label>
                    <input type="email" id="email" name="email" 
                           value="{{if and .Cliente .Cliente.Email}}{{.Cliente.Email}}{{end}}">
                </div>

                <div class="form-group">
                    <label for="telefone">Telefone</label>
                    <input type="tel" id="telefone" name="telefone" 
                           value="{{if and .Cliente .Cliente.Telefone}}{{.Cliente.Telefone}}{{end}}">
                </div>

                <div class="form-group">
                    <label for="cpf_cnpj">CPF/CNPJ</label>
                    <input type="text" id="cpf_cnpj" name="cpf_cnpj" 
                           value="{{if and .Cliente .Cliente.CPFCNPJ}}{{.Cliente.CPFCNPJ}}{{end}}">
                </div>

                <div class="form-group">
                    <label for="cep">CEP</label>
                    <input type="text" id="cep" name="cep" 
                           value="{{if and .Cliente .Cliente.CEP}}{{.Cliente.CEP}}{{end}}">
                </div>

                <div class="form-group full-width">
                    <label for="endereco">Endereço</label>
                    <input type="text" id="endereco" name="endereco" 
                           value="{{if and .Cliente .Cliente.Endereco}}{{.Cliente.Endereco}}{{end}}">
                </div>

                <div class="form-group">
                    <label for="cidade">Cidade</label>
                    <input type="text" id="cidade" name="cidade" 
                           value="{{if and .Cliente .Cliente.Cidade}}{{.Cliente.Cidade}}{{end}}">
                </div>

                <div class="form-group">
                    <label for="estado">Estado</label>
                    <input type="text" id="estado" name="estado" maxlength="2" 
                           value="{{if and .Cliente .Cliente.Estado}}{{.Cliente.Estado}}{{end}}">
                </div>

                <div class="form-group full-width">
                    <label for="observacoes">Observações</label>
                    <textarea id="observacoes" name="observacoes">{{if and .Cliente .Cliente.Observacoes}}{{.Cliente.Observacoes}}{{end}}</textarea>
                </div>

                <div class="form-group">
                    <div class="checkbox-group">
                        <input type="checkbox" id="ativo" name="ativo" 
                               {{if not .Cliente}}checked{{else}}{{if .Cliente.Ativo}}checked{{end}}{{end}}>
                        <label for="ativo">Cliente Ativo</label>
                    </div>
                </div>
            </div>

            <div class="actions">
                <a href="/clientes" class="btn btn-secondary">Cancelar</a>
                <button type="submit" class="btn">{{if .IsEdit}}Atualizar{{else}}Criar{{end}} Cliente</button>
            </div>
        </form>
    </div>
</body>
</html>`

	t, err := template.New("form").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erro no template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Erro ao renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}