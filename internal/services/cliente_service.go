package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"webserver-go/internal/database"
	"webserver-go/internal/models"

	"github.com/google/uuid"
)

type ClienteService struct {
	db *sql.DB
}

func NewClienteService() *ClienteService {
	return &ClienteService{
		db: database.DB,
	}
}

// CreateCliente cria um novo cliente
func (s *ClienteService) CreateCliente(input models.ClienteInput) (*models.Cliente, error) {
	// Validações básicas
	if err := s.validateClienteInput(input); err != nil {
		return nil, err
	}

	// Gerar ID
	id := uuid.New().String()
	now := time.Now()

	// Valores padrão
	ativo := true
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	cliente := &models.Cliente{
		ID:          id,
		Nome:        input.Nome,
		Email:       input.Email,
		Telefone:    input.Telefone,
		Endereco:    input.Endereco,
		Cidade:      input.Cidade,
		Estado:      input.Estado,
		CEP:         input.CEP,
		CPFCNPJ:     input.CPFCNPJ,
		TipoPessoa:  input.TipoPessoa,
		Ativo:       ativo,
		Observacoes: input.Observacoes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	query := `
		INSERT INTO clientes (
			id, nome, email, telefone, endereco, cidade, estado, cep, 
			cpf_cnpj, tipo_pessoa, ativo, observacoes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := s.db.Exec(query,
		cliente.ID, cliente.Nome, cliente.Email, cliente.Telefone,
		cliente.Endereco, cliente.Cidade, cliente.Estado, cliente.CEP,
		cliente.CPFCNPJ, cliente.TipoPessoa, cliente.Ativo,
		cliente.Observacoes, cliente.CreatedAt, cliente.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "email") {
				return nil, fmt.Errorf("email já está em uso")
			}
			if strings.Contains(err.Error(), "cpf_cnpj") {
				return nil, fmt.Errorf("CPF/CNPJ já está em uso")
			}
		}
		return nil, fmt.Errorf("erro ao criar cliente: %w", err)
	}

	return cliente, nil
}

// GetAllClientes retorna todos os clientes com paginação
func (s *ClienteService) GetAllClientes(page, limit int, search string) ([]models.Cliente, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Query base
	var countQuery string
	var total int
	var err error

	if search == "" {
		// Sem filtro de busca
		countQuery = "SELECT COUNT(*) FROM clientes"
		err = s.db.QueryRow(countQuery).Scan(&total)
	} else {
		// Com filtro de busca
		countQuery = "SELECT COUNT(*) FROM clientes WHERE nome ILIKE $1 OR email ILIKE $1"
		searchPattern := "%" + search + "%"
		err = s.db.QueryRow(countQuery, searchPattern).Scan(&total)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar clientes: %w", err)
	}

	// Buscar clientes
	var query string
	var rows *sql.Rows

	if search == "" {
		// Sem filtro de busca
		query = `
			SELECT id, nome, email, telefone, endereco, cidade, estado, cep,
				   cpf_cnpj, tipo_pessoa, ativo, observacoes, created_at, updated_at
			FROM clientes 
			ORDER BY nome ASC
			LIMIT $1 OFFSET $2
		`
		rows, err = s.db.Query(query, limit, offset)
	} else {
		// Com filtro de busca
		query = `
			SELECT id, nome, email, telefone, endereco, cidade, estado, cep,
				   cpf_cnpj, tipo_pessoa, ativo, observacoes, created_at, updated_at
			FROM clientes 
			WHERE nome ILIKE $1 OR email ILIKE $1
			ORDER BY nome ASC
			LIMIT $2 OFFSET $3
		`
		searchPattern := "%" + search + "%"
		rows, err = s.db.Query(query, searchPattern, limit, offset)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar clientes: %w", err)
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var cliente models.Cliente
		err := rows.Scan(
			&cliente.ID, &cliente.Nome, &cliente.Email, &cliente.Telefone,
			&cliente.Endereco, &cliente.Cidade, &cliente.Estado, &cliente.CEP,
			&cliente.CPFCNPJ, &cliente.TipoPessoa, &cliente.Ativo,
			&cliente.Observacoes, &cliente.CreatedAt, &cliente.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("erro ao escanear cliente: %w", err)
		}
		clientes = append(clientes, cliente)
	}

	return clientes, total, nil
}

// GetClienteByID retorna um cliente pelo ID
func (s *ClienteService) GetClienteByID(id string) (*models.Cliente, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("ID inválido")
	}

	query := `
		SELECT id, nome, email, telefone, endereco, cidade, estado, cep,
			   cpf_cnpj, tipo_pessoa, ativo, observacoes, created_at, updated_at
		FROM clientes 
		WHERE id = $1
	`

	var cliente models.Cliente
	err := s.db.QueryRow(query, id).Scan(
		&cliente.ID, &cliente.Nome, &cliente.Email, &cliente.Telefone,
		&cliente.Endereco, &cliente.Cidade, &cliente.Estado, &cliente.CEP,
		&cliente.CPFCNPJ, &cliente.TipoPessoa, &cliente.Ativo,
		&cliente.Observacoes, &cliente.CreatedAt, &cliente.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cliente não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	return &cliente, nil
}

// UpdateCliente atualiza um cliente
func (s *ClienteService) UpdateCliente(id string, input models.ClienteInput) (*models.Cliente, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("ID inválido")
	}

	// Validações
	if err := s.validateClienteInput(input); err != nil {
		return nil, err
	}

	// Verificar se cliente existe
	existingCliente, err := s.GetClienteByID(id)
	if err != nil {
		return nil, err
	}

	// Valores padrão
	ativo := existingCliente.Ativo
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	query := `
		UPDATE clientes SET 
			nome = $2, email = $3, telefone = $4, endereco = $5, 
			cidade = $6, estado = $7, cep = $8, cpf_cnpj = $9, 
			tipo_pessoa = $10, ativo = $11, observacoes = $12, updated_at = $13
		WHERE id = $1
	`

	_, err = s.db.Exec(query,
		id, input.Nome, input.Email, input.Telefone,
		input.Endereco, input.Cidade, input.Estado, input.CEP,
		input.CPFCNPJ, input.TipoPessoa, ativo,
		input.Observacoes, time.Now(),
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "email") {
				return nil, fmt.Errorf("email já está em uso")
			}
			if strings.Contains(err.Error(), "cpf_cnpj") {
				return nil, fmt.Errorf("CPF/CNPJ já está em uso")
			}
		}
		return nil, fmt.Errorf("erro ao atualizar cliente: %w", err)
	}

	// Retornar cliente atualizado
	return s.GetClienteByID(id)
}

// DeleteCliente deleta um cliente
func (s *ClienteService) DeleteCliente(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("ID inválido")
	}

	// Verificar se cliente existe
	_, err := s.GetClienteByID(id)
	if err != nil {
		return err
	}

	query := "DELETE FROM clientes WHERE id = $1"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar cliente: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar deleção: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cliente não encontrado")
	}

	return nil
}

// validateClienteInput valida os dados de entrada
func (s *ClienteService) validateClienteInput(input models.ClienteInput) error {
	if strings.TrimSpace(input.Nome) == "" {
		return fmt.Errorf("nome é obrigatório")
	}

	if len(input.Nome) < 2 || len(input.Nome) > 255 {
		return fmt.Errorf("nome deve ter entre 2 e 255 caracteres")
	}

	if input.TipoPessoa != "fisica" && input.TipoPessoa != "juridica" {
		return fmt.Errorf("tipo de pessoa deve ser 'fisica' ou 'juridica'")
	}

	// Validar email se fornecido
	if input.Email != nil && *input.Email != "" {
		email := strings.TrimSpace(*input.Email)
		if !strings.Contains(email, "@") || len(email) < 5 {
			return fmt.Errorf("email inválido")
		}
	}

	// Validar estado se fornecido
	if input.Estado != nil && *input.Estado != "" {
		if len(*input.Estado) != 2 {
			return fmt.Errorf("estado deve ter 2 caracteres")
		}
	}

	// Validar CEP se fornecido
	if input.CEP != nil && *input.CEP != "" {
		cep := strings.ReplaceAll(*input.CEP, "-", "")
		if len(cep) != 8 {
			return fmt.Errorf("CEP deve ter 8 dígitos")
		}
	}

	return nil
}