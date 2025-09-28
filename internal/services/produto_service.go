package services

import (
	"fmt"
	"strings"

	"webserver-go/internal/database"
	"webserver-go/internal/models"

	"gorm.io/gorm"
)

type ProdutoService struct {
	db *gorm.DB
}

func NewProdutoService() *ProdutoService {
	return &ProdutoService{
		db: database.GormDB,
	}
}

// CreateProduto cria um novo produto usando GORM
func (s *ProdutoService) CreateProduto(input models.ProdutoGormInput) (*models.ProdutoGorm, error) {
	// Validações básicas
	if err := s.validateProdutoInput(input); err != nil {
		return nil, err
	}

	// Valores padrão
	ativo := true
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	unidadeMedida := "UN"
	if input.UnidadeMedida != "" {
		unidadeMedida = input.UnidadeMedida
	}

	produto := &models.ProdutoGorm{
		Nome:           input.Nome,
		Descricao:      input.Descricao,
		CodigoBarras:   input.CodigoBarras,
		CategoriaID:    input.CategoriaID,
		PrecoCusto:     input.PrecoCusto,
		PrecoVenda:     input.PrecoVenda,
		MargemLucro:    input.MargemLucro,
		UnidadeMedida:  unidadeMedida,
		Peso:           input.Peso,
		Dimensoes:      input.Dimensoes,
		Ativo:          ativo,
		ImagemURL:      input.ImagemURL,
		Observacoes:    input.Observacoes,
	}

	// GORM automaticamente gera ID, CreatedAt e UpdatedAt
	if err := s.db.Create(produto).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "codigo_barras") {
				return nil, fmt.Errorf("código de barras já está em uso")
			}
		}
		return nil, fmt.Errorf("erro ao criar produto: %w", err)
	}

	// Carregar categoria se existir
	if produto.CategoriaID != nil {
		s.db.Preload("Categoria").First(produto, produto.ID)
	}

	return produto, nil
}

// GetAllProdutos retorna todos os produtos com paginação
func (s *ProdutoService) GetAllProdutos(page, limit int, search string, categoriaID string) ([]models.ProdutoGorm, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var produtos []models.ProdutoGorm
	var total int64

	// Query base
	query := s.db.Model(&models.ProdutoGorm{}).Preload("Categoria")

	// Filtros
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nome ILIKE ? OR codigo_barras ILIKE ?", searchPattern, searchPattern)
	}

	if categoriaID != "" {
		query = query.Where("categoria_id = ?", categoriaID)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao contar produtos: %w", err)
	}

	// Buscar produtos com paginação
	if err := query.Order("nome ASC").Limit(limit).Offset(offset).Find(&produtos).Error; err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar produtos: %w", err)
	}

	return produtos, total, nil
}

// GetProdutoByID retorna um produto pelo ID
func (s *ProdutoService) GetProdutoByID(id string) (*models.ProdutoGorm, error) {
	var produto models.ProdutoGorm
	
	if err := s.db.Preload("Categoria").First(&produto, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("produto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	return &produto, nil
}

// UpdateProduto atualiza um produto
func (s *ProdutoService) UpdateProduto(id string, input models.ProdutoGormInput) (*models.ProdutoGorm, error) {
	// Validações
	if err := s.validateProdutoInput(input); err != nil {
		return nil, err
	}

	// Buscar produto existente
	produto, err := s.GetProdutoByID(id)
	if err != nil {
		return nil, err
	}

	// Valores padrão
	ativo := produto.Ativo
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	unidadeMedida := produto.UnidadeMedida
	if input.UnidadeMedida != "" {
		unidadeMedida = input.UnidadeMedida
	}

	// Atualizar campos
	updates := models.ProdutoGorm{
		Nome:           input.Nome,
		Descricao:      input.Descricao,
		CodigoBarras:   input.CodigoBarras,
		CategoriaID:    input.CategoriaID,
		PrecoCusto:     input.PrecoCusto,
		PrecoVenda:     input.PrecoVenda,
		MargemLucro:    input.MargemLucro,
		UnidadeMedida:  unidadeMedida,
		Peso:           input.Peso,
		Dimensoes:      input.Dimensoes,
		Ativo:          ativo,
		ImagemURL:      input.ImagemURL,
		Observacoes:    input.Observacoes,
	}

	if err := s.db.Model(produto).Updates(updates).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "codigo_barras") {
				return nil, fmt.Errorf("código de barras já está em uso")
			}
		}
		return nil, fmt.Errorf("erro ao atualizar produto: %w", err)
	}

	// Retornar produto atualizado
	return s.GetProdutoByID(id)
}

// DeleteProduto deleta um produto (soft delete)
func (s *ProdutoService) DeleteProduto(id string) error {
	// Verificar se produto existe
	produto, err := s.GetProdutoByID(id)
	if err != nil {
		return err
	}

	// Verificar se há estoque para este produto
	var count int64
	if err := s.db.Model(&models.EstoqueGorm{}).Where("produto_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("erro ao verificar estoque: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("não é possível deletar produto que possui movimentações de estoque")
	}

	// Soft delete (GORM automaticamente usa deleted_at)
	if err := s.db.Delete(produto).Error; err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}

	return nil
}

// GetProdutosByCategoria retorna produtos de uma categoria específica
func (s *ProdutoService) GetProdutosByCategoria(categoriaID string) ([]models.ProdutoGorm, error) {
	var produtos []models.ProdutoGorm
	
	if err := s.db.Preload("Categoria").Where("categoria_id = ?", categoriaID).Order("nome ASC").Find(&produtos).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos da categoria: %w", err)
	}

	return produtos, nil
}

// GetProdutosComEstoqueBaixo retorna produtos com estoque abaixo do mínimo
func (s *ProdutoService) GetProdutosComEstoqueBaixo() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	rows, err := s.db.Raw(`
		SELECT 
			p.id, p.nome, p.preco_venda, p.unidade_medida,
			e.quantidade_atual, e.quantidade_minima,
			c.nome as categoria_nome
		FROM produtos p
		LEFT JOIN estoque e ON p.id = e.produto_id
		LEFT JOIN categorias c ON p.categoria_id = c.id
		WHERE p.deleted_at IS NULL 
		  AND p.ativo = true
		  AND (e.quantidade_atual < e.quantidade_minima OR e.quantidade_atual IS NULL)
		ORDER BY p.nome ASC
	`).Rows()
	
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos com estoque baixo: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var produto models.ProdutoGorm
		var quantidadeAtual, quantidadeMinima *int
		var categoriaNome *string
		
		if err := rows.Scan(
			&produto.ID, &produto.Nome, &produto.PrecoVenda, &produto.UnidadeMedida,
			&quantidadeAtual, &quantidadeMinima, &categoriaNome,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear resultado: %w", err)
		}
		
		result := map[string]interface{}{
			"produto":            produto,
			"quantidade_atual":   quantidadeAtual,
			"quantidade_minima":  quantidadeMinima,
			"categoria_nome":     categoriaNome,
		}
		results = append(results, result)
	}

	return results, nil
}

// GetRelatorioProdutos retorna estatísticas dos produtos
func (s *ProdutoService) GetRelatorioProdutos() (map[string]interface{}, error) {
	var totalProdutos, produtosAtivos, produtosInativos int64
	var valorTotalEstoque float64

	// Total de produtos
	if err := s.db.Model(&models.ProdutoGorm{}).Count(&totalProdutos).Error; err != nil {
		return nil, fmt.Errorf("erro ao contar produtos: %w", err)
	}

	// Produtos ativos
	if err := s.db.Model(&models.ProdutoGorm{}).Where("ativo = ?", true).Count(&produtosAtivos).Error; err != nil {
		return nil, fmt.Errorf("erro ao contar produtos ativos: %w", err)
	}

	// Produtos inativos
	if err := s.db.Model(&models.ProdutoGorm{}).Where("ativo = ?", false).Count(&produtosInativos).Error; err != nil {
		return nil, fmt.Errorf("erro ao contar produtos inativos: %w", err)
	}

	// Valor total do estoque
	row := s.db.Raw(`
		SELECT COALESCE(SUM(p.preco_custo * e.quantidade_atual), 0)
		FROM produtos p
		INNER JOIN estoque e ON p.id = e.produto_id
		WHERE p.deleted_at IS NULL AND p.ativo = true AND p.preco_custo IS NOT NULL
	`).Row()
	
	if err := row.Scan(&valorTotalEstoque); err != nil {
		valorTotalEstoque = 0 // Se não conseguir calcular, assume 0
	}

	return map[string]interface{}{
		"total_produtos":     totalProdutos,
		"produtos_ativos":    produtosAtivos,
		"produtos_inativos":  produtosInativos,
		"valor_total_estoque": valorTotalEstoque,
	}, nil
}

// validateProdutoInput valida os dados de entrada
func (s *ProdutoService) validateProdutoInput(input models.ProdutoGormInput) error {
	if strings.TrimSpace(input.Nome) == "" {
		return fmt.Errorf("nome é obrigatório")
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

	if input.MargemLucro != nil && (*input.MargemLucro < 0 || *input.MargemLucro > 100) {
		return fmt.Errorf("margem de lucro deve estar entre 0 e 100")
	}

	if input.Peso != nil && *input.Peso < 0 {
		return fmt.Errorf("peso não pode ser negativo")
	}

	// Validar categoria se fornecida
	if input.CategoriaID != nil && *input.CategoriaID != "" {
		var count int64
		if err := s.db.Model(&models.CategoriaGorm{}).Where("id = ? AND ativo = ?", *input.CategoriaID, true).Count(&count).Error; err != nil {
			return fmt.Errorf("erro ao validar categoria: %w", err)
		}
		if count == 0 {
			return fmt.Errorf("categoria não encontrada ou inativa")
		}
	}

	return nil
}