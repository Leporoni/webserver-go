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

// ProdutoService handles business logic for produtos
type ProdutoService struct {
	db *sql.DB
}

// NewProdutoService creates a new produto service
func NewProdutoService() *ProdutoService {
	return &ProdutoService{
		db: database.DB,
	}
}

// CreateProduto creates a new produto
func (s *ProdutoService) CreateProduto(input *models.ProdutoInput) (*models.Produto, error) {
	// Validações básicas
	if err := s.validateProdutoInput(input); err != nil {
		return nil, err
	}

	// Verificar se código de barras já existe (se fornecido)
	if input.CodigoBarras != nil && *input.CodigoBarras != "" {
		exists, err := s.codigoBarrasExists(*input.CodigoBarras, "")
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar código de barras: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("código de barras já existe")
		}
	}

	// Verificar se categoria existe (se fornecida)
	if input.CategoriaID != nil {
		exists, err := s.categoriaExists(*input.CategoriaID)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar categoria: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("categoria não encontrada")
		}
	}

	// Calcular margem de lucro se não fornecida
	if input.MargemLucro == nil && input.PrecoCusto != nil {
		margem := ((input.PrecoVenda - *input.PrecoCusto) / *input.PrecoCusto) * 100
		input.MargemLucro = &margem
	}

	// Definir valores padrão
	ativo := true
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	// Criar produto
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO produtos (
			id, nome, descricao, codigo_barras, categoria_id, 
			preco_custo, preco_venda, margem_lucro, unidade_medida, 
			peso, dimensoes, ativo, imagem_url, observacoes, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)`

	_, err := s.db.Exec(query,
		id, input.Nome, input.Descricao, input.CodigoBarras, input.CategoriaID,
		input.PrecoCusto, input.PrecoVenda, input.MargemLucro, input.UnidadeMedida,
		input.Peso, input.Dimensoes, ativo, input.ImagemURL, input.Observacoes,
		now, now,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar produto: %w", err)
	}

	// Buscar produto criado
	return s.GetProdutoByID(id)
}

// GetAllProdutos returns paginated list of produtos
func (s *ProdutoService) GetAllProdutos(page, limit int, search, categoriaID string) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Base query
	baseQuery := `
		FROM produtos p
		LEFT JOIN categorias c ON p.categoria_id = c.id
		WHERE 1=1
	`
	
	var args []interface{}
	argIndex := 1

	// Filtros
	if search != "" {
		baseQuery += fmt.Sprintf(" AND (LOWER(p.nome) LIKE LOWER($%d) OR LOWER(p.descricao) LIKE LOWER($%d) OR p.codigo_barras LIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if categoriaID != "" {
		baseQuery += fmt.Sprintf(" AND p.categoria_id = $%d", argIndex)
		args = append(args, categoriaID)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("erro ao contar produtos: %w", err)
	}

	// Get produtos
	selectQuery := `
		SELECT 
			p.id, p.nome, p.descricao, p.codigo_barras, p.categoria_id,
			p.preco_custo, p.preco_venda, p.margem_lucro, p.unidade_medida,
			p.peso, p.dimensoes, p.ativo, p.imagem_url, p.observacoes,
			p.created_at, p.updated_at,
			c.nome as categoria_nome
		` + baseQuery + `
		ORDER BY p.nome ASC
		LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)

	args = append(args, limit, offset)

	rows, err := s.db.Query(selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %w", err)
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		var categoriaNome sql.NullString

		err := rows.Scan(
			&p.ID, &p.Nome, &p.Descricao, &p.CodigoBarras, &p.CategoriaID,
			&p.PrecoCusto, &p.PrecoVenda, &p.MargemLucro, &p.UnidadeMedida,
			&p.Peso, &p.Dimensoes, &p.Ativo, &p.ImagemURL, &p.Observacoes,
			&p.CreatedAt, &p.UpdatedAt, &categoriaNome,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %w", err)
		}

		// Adicionar categoria se existir
		if categoriaNome.Valid && p.CategoriaID != nil {
			p.Categoria = &models.Categoria{
				ID:   *p.CategoriaID,
				Nome: categoriaNome.String,
			}
		}

		produtos = append(produtos, p)
	}

	totalPages := (total + limit - 1) / limit

	return &models.PaginatedResponse{
		Data:       produtos,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetProdutoByID returns a produto by ID
func (s *ProdutoService) GetProdutoByID(id string) (*models.Produto, error) {
	query := `
		SELECT 
			p.id, p.nome, p.descricao, p.codigo_barras, p.categoria_id,
			p.preco_custo, p.preco_venda, p.margem_lucro, p.unidade_medida,
			p.peso, p.dimensoes, p.ativo, p.imagem_url, p.observacoes,
			p.created_at, p.updated_at,
			c.id as categoria_id_full, c.nome as categoria_nome, 
			c.descricao as categoria_descricao, c.ativo as categoria_ativo,
			c.created_at as categoria_created_at, c.updated_at as categoria_updated_at
		FROM produtos p
		LEFT JOIN categorias c ON p.categoria_id = c.id
		WHERE p.id = $1
	`

	var p models.Produto
	var categoria models.Categoria
	var categoriaID, categoriaNome, categoriaDescricao sql.NullString
	var categoriaAtivo sql.NullBool
	var categoriaCreatedAt, categoriaUpdatedAt sql.NullTime

	err := s.db.QueryRow(query, id).Scan(
		&p.ID, &p.Nome, &p.Descricao, &p.CodigoBarras, &p.CategoriaID,
		&p.PrecoCusto, &p.PrecoVenda, &p.MargemLucro, &p.UnidadeMedida,
		&p.Peso, &p.Dimensoes, &p.Ativo, &p.ImagemURL, &p.Observacoes,
		&p.CreatedAt, &p.UpdatedAt,
		&categoriaID, &categoriaNome, &categoriaDescricao, &categoriaAtivo,
		&categoriaCreatedAt, &categoriaUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	// Adicionar categoria se existir
	if categoriaID.Valid {
		categoria.ID = categoriaID.String
		categoria.Nome = categoriaNome.String
		if categoriaDescricao.Valid {
			categoria.Descricao = &categoriaDescricao.String
		}
		categoria.Ativo = categoriaAtivo.Bool
		categoria.CreatedAt = categoriaCreatedAt.Time
		categoria.UpdatedAt = categoriaUpdatedAt.Time
		p.Categoria = &categoria
	}

	return &p, nil
}

// UpdateProduto updates an existing produto
func (s *ProdutoService) UpdateProduto(id string, input *models.ProdutoInput) (*models.Produto, error) {
	// Verificar se produto existe
	existing, err := s.GetProdutoByID(id)
	if err != nil {
		return nil, err
	}

	// Validações básicas
	if err := s.validateProdutoInput(input); err != nil {
		return nil, err
	}

	// Verificar se código de barras já existe (se fornecido e diferente do atual)
	if input.CodigoBarras != nil && *input.CodigoBarras != "" {
		if existing.CodigoBarras == nil || *existing.CodigoBarras != *input.CodigoBarras {
			exists, err := s.codigoBarrasExists(*input.CodigoBarras, id)
			if err != nil {
				return nil, fmt.Errorf("erro ao verificar código de barras: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("código de barras já existe")
			}
		}
	}

	// Verificar se categoria existe (se fornecida)
	if input.CategoriaID != nil {
		exists, err := s.categoriaExists(*input.CategoriaID)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar categoria: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("categoria não encontrada")
		}
	}

	// Calcular margem de lucro se não fornecida
	if input.MargemLucro == nil && input.PrecoCusto != nil {
		margem := ((input.PrecoVenda - *input.PrecoCusto) / *input.PrecoCusto) * 100
		input.MargemLucro = &margem
	}

	// Definir valores padrão
	ativo := existing.Ativo
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	// Atualizar produto
	query := `
		UPDATE produtos SET
			nome = $2, descricao = $3, codigo_barras = $4, categoria_id = $5,
			preco_custo = $6, preco_venda = $7, margem_lucro = $8, unidade_medida = $9,
			peso = $10, dimensoes = $11, ativo = $12, imagem_url = $13, observacoes = $14,
			updated_at = $15
		WHERE id = $1
	`

	_, err = s.db.Exec(query,
		id, input.Nome, input.Descricao, input.CodigoBarras, input.CategoriaID,
		input.PrecoCusto, input.PrecoVenda, input.MargemLucro, input.UnidadeMedida,
		input.Peso, input.Dimensoes, ativo, input.ImagemURL, input.Observacoes,
		time.Now(),
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar produto: %w", err)
	}

	// Retornar produto atualizado
	return s.GetProdutoByID(id)
}

// DeleteProduto deletes a produto
func (s *ProdutoService) DeleteProduto(id string) error {
	// Verificar se produto existe
	_, err := s.GetProdutoByID(id)
	if err != nil {
		return err
	}

	// Verificar se produto tem estoque (futuro)
	// TODO: Implementar verificação de estoque quando módulo estiver pronto

	// Deletar produto
	query := "DELETE FROM produtos WHERE id = $1"
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar deleção: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	return nil
}

// GetCategorias returns all categorias
func (s *ProdutoService) GetCategorias() ([]models.Categoria, error) {
	query := `
		SELECT id, nome, descricao, ativo, created_at, updated_at
		FROM categorias
		WHERE ativo = true
		ORDER BY nome ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}
	defer rows.Close()

	var categorias []models.Categoria
	for rows.Next() {
		var c models.Categoria
		err := rows.Scan(&c.ID, &c.Nome, &c.Descricao, &c.Ativo, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear categoria: %w", err)
		}
		categorias = append(categorias, c)
	}

	return categorias, nil
}

// CreateCategoria creates a new categoria
func (s *ProdutoService) CreateCategoria(input *models.CategoriaInput) (*models.Categoria, error) {
	// Validações básicas
	if strings.TrimSpace(input.Nome) == "" {
		return nil, fmt.Errorf("nome da categoria é obrigatório")
	}

	// Verificar se nome já existe
	exists, err := s.categoriaNomeExists(input.Nome, "")
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar nome da categoria: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("já existe uma categoria com este nome")
	}

	// Definir valores padrão
	ativo := true
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	// Criar categoria
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO categorias (id, nome, descricao, ativo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = s.db.Exec(query, id, input.Nome, input.Descricao, ativo, now, now)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar categoria: %w", err)
	}

	// Retornar categoria criada
	return &models.Categoria{
		ID:        id,
		Nome:      input.Nome,
		Descricao: input.Descricao,
		Ativo:     ativo,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Helper functions

func (s *ProdutoService) validateProdutoInput(input *models.ProdutoInput) error {
	if strings.TrimSpace(input.Nome) == "" {
		return fmt.Errorf("nome do produto é obrigatório")
	}

	if len(input.Nome) < 2 || len(input.Nome) > 255 {
		return fmt.Errorf("nome deve ter entre 2 e 255 caracteres")
	}

	if input.PrecoVenda <= 0 {
		return fmt.Errorf("preço de venda deve ser maior que zero")
	}

	if input.PrecoCusto != nil && *input.PrecoCusto < 0 {
		return fmt.Errorf("preço de custo não pode ser negativo")
	}

	if input.PrecoCusto != nil && *input.PrecoCusto >= input.PrecoVenda {
		return fmt.Errorf("preço de custo deve ser menor que preço de venda")
	}

	if strings.TrimSpace(input.UnidadeMedida) == "" {
		return fmt.Errorf("unidade de medida é obrigatória")
	}

	if input.Peso != nil && *input.Peso < 0 {
		return fmt.Errorf("peso não pode ser negativo")
	}

	return nil
}

func (s *ProdutoService) codigoBarrasExists(codigoBarras, excludeID string) (bool, error) {
	query := "SELECT COUNT(*) FROM produtos WHERE codigo_barras = $1"
	args := []interface{}{codigoBarras}

	if excludeID != "" {
		query += " AND id != $2"
		args = append(args, excludeID)
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	return count > 0, err
}

func (s *ProdutoService) categoriaExists(categoriaID string) (bool, error) {
	query := "SELECT COUNT(*) FROM categorias WHERE id = $1 AND ativo = true"
	var count int
	err := s.db.QueryRow(query, categoriaID).Scan(&count)
	return count > 0, err
}

func (s *ProdutoService) categoriaNomeExists(nome, excludeID string) (bool, error) {
	query := "SELECT COUNT(*) FROM categorias WHERE LOWER(nome) = LOWER($1)"
	args := []interface{}{nome}

	if excludeID != "" {
		query += " AND id != $2"
		args = append(args, excludeID)
	}

	var count int
	err := s.db.QueryRow(query, args...).Scan(&count)
	return count > 0, err
}