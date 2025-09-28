package services

import (
	"fmt"
	"strings"

	"webserver-go/internal/database"
	"webserver-go/internal/models"

	"gorm.io/gorm"
)

type CategoriaService struct {
	db *gorm.DB
}

func NewCategoriaService() *CategoriaService {
	return &CategoriaService{
		db: database.GormDB,
	}
}

// CreateCategoria cria uma nova categoria usando GORM
func (s *CategoriaService) CreateCategoria(input models.CategoriaGormInput) (*models.CategoriaGorm, error) {
	// Validações básicas
	if err := s.validateCategoriaInput(input); err != nil {
		return nil, err
	}

	// Valores padrão
	ativo := true
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	categoria := &models.CategoriaGorm{
		Nome:      input.Nome,
		Descricao: input.Descricao,
		Ativo:     ativo,
	}

	// GORM automaticamente gera ID, CreatedAt e UpdatedAt
	if err := s.db.Create(categoria).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("categoria com este nome já existe")
		}
		return nil, fmt.Errorf("erro ao criar categoria: %w", err)
	}

	return categoria, nil
}

// GetAllCategorias retorna todas as categorias
func (s *CategoriaService) GetAllCategorias(apenasAtivas bool) ([]models.CategoriaGorm, error) {
	var categorias []models.CategoriaGorm
	
	query := s.db.Order("nome ASC")
	
	if apenasAtivas {
		query = query.Where("ativo = ?", true)
	}

	if err := query.Find(&categorias).Error; err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}

	return categorias, nil
}

// GetCategoriaByID retorna uma categoria pelo ID
func (s *CategoriaService) GetCategoriaByID(id string) (*models.CategoriaGorm, error) {
	var categoria models.CategoriaGorm
	
	if err := s.db.First(&categoria, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("categoria não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	return &categoria, nil
}

// UpdateCategoria atualiza uma categoria
func (s *CategoriaService) UpdateCategoria(id string, input models.CategoriaGormInput) (*models.CategoriaGorm, error) {
	// Validações
	if err := s.validateCategoriaInput(input); err != nil {
		return nil, err
	}

	// Buscar categoria existente
	categoria, err := s.GetCategoriaByID(id)
	if err != nil {
		return nil, err
	}

	// Valores padrão
	ativo := categoria.Ativo
	if input.Ativo != nil {
		ativo = *input.Ativo
	}

	// Atualizar campos
	updates := models.CategoriaGorm{
		Nome:      input.Nome,
		Descricao: input.Descricao,
		Ativo:     ativo,
	}

	if err := s.db.Model(categoria).Updates(updates).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("categoria com este nome já existe")
		}
		return nil, fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	// Retornar categoria atualizada
	return s.GetCategoriaByID(id)
}

// DeleteCategoria deleta uma categoria
func (s *CategoriaService) DeleteCategoria(id string) error {
	// Verificar se categoria existe
	categoria, err := s.GetCategoriaByID(id)
	if err != nil {
		return err
	}

	// Verificar se há produtos usando esta categoria
	var count int64
	if err := s.db.Model(&models.ProdutoGorm{}).Where("categoria_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("erro ao verificar produtos: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("não é possível deletar categoria que possui produtos associados")
	}

	// Deletar categoria
	if err := s.db.Delete(categoria).Error; err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}

	return nil
}

// GetCategoriasWithProductCount retorna categorias com contagem de produtos
func (s *CategoriaService) GetCategoriasWithProductCount() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	rows, err := s.db.Raw(`
		SELECT 
			c.id, c.nome, c.descricao, c.ativo, c.created_at, c.updated_at,
			COALESCE(COUNT(p.id), 0) as produto_count
		FROM categorias c
		LEFT JOIN produtos p ON c.id = p.categoria_id AND p.deleted_at IS NULL
		GROUP BY c.id, c.nome, c.descricao, c.ativo, c.created_at, c.updated_at
		ORDER BY c.nome ASC
	`).Rows()
	
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias com contagem: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var categoria models.CategoriaGorm
		var produtoCount int
		
		if err := rows.Scan(
			&categoria.ID, &categoria.Nome, &categoria.Descricao, 
			&categoria.Ativo, &categoria.CreatedAt, &categoria.UpdatedAt,
			&produtoCount,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear resultado: %w", err)
		}
		
		result := map[string]interface{}{
			"categoria":      categoria,
			"produto_count": produtoCount,
		}
		results = append(results, result)
	}

	return results, nil
}

// validateCategoriaInput valida os dados de entrada
func (s *CategoriaService) validateCategoriaInput(input models.CategoriaGormInput) error {
	if strings.TrimSpace(input.Nome) == "" {
		return fmt.Errorf("nome é obrigatório")
	}

	if len(input.Nome) < 2 || len(input.Nome) > 255 {
		return fmt.Errorf("nome deve ter entre 2 e 255 caracteres")
	}

	return nil
}