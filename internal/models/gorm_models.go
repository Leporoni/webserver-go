package models

import (
	"time"
	"gorm.io/gorm"
)

// CategoriaGorm representa uma categoria usando GORM
type CategoriaGorm struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Nome      string    `gorm:"not null;size:255;uniqueIndex" json:"nome"`
	Descricao *string   `gorm:"type:text" json:"descricao"`
	Ativo     bool      `gorm:"default:true" json:"ativo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relacionamentos
	Produtos []ProdutoGorm `gorm:"foreignKey:CategoriaID" json:"produtos,omitempty"`
}

// TableName especifica o nome da tabela
func (CategoriaGorm) TableName() string {
	return "categorias"
}

// ProdutoGorm representa um produto usando GORM
type ProdutoGorm struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Nome           string         `gorm:"not null;size:255" json:"nome"`
	Descricao      *string        `gorm:"type:text" json:"descricao"`
	CodigoBarras   *string        `gorm:"size:50;uniqueIndex" json:"codigo_barras"`
	CategoriaID    *string        `gorm:"type:uuid" json:"categoria_id"`
	PrecoCusto     *float64       `gorm:"type:decimal(10,2)" json:"preco_custo"`
	PrecoVenda     float64        `gorm:"type:decimal(10,2);not null" json:"preco_venda"`
	MargemLucro    *float64       `gorm:"type:decimal(5,2)" json:"margem_lucro"`
	UnidadeMedida  string         `gorm:"size:10;default:'UN'" json:"unidade_medida"`
	Peso           *float64       `gorm:"type:decimal(8,3)" json:"peso"`
	Dimensoes      *string        `gorm:"size:50" json:"dimensoes"`
	Ativo          bool           `gorm:"default:true" json:"ativo"`
	ImagemURL      *string        `gorm:"size:500" json:"imagem_url"`
	Observacoes    *string        `gorm:"type:text" json:"observacoes"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
	
	// Relacionamentos
	Categoria      *CategoriaGorm `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
	Estoque        []EstoqueGorm  `gorm:"foreignKey:ProdutoID" json:"estoque,omitempty"`
}

// TableName especifica o nome da tabela
func (ProdutoGorm) TableName() string {
	return "produtos"
}

// EstoqueGorm representa o estoque usando GORM
type EstoqueGorm struct {
	ID                string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ProdutoID         string     `gorm:"type:uuid;not null" json:"produto_id"`
	QuantidadeAtual   int        `gorm:"default:0" json:"quantidade_atual"`
	QuantidadeMinima  int        `gorm:"default:0" json:"quantidade_minima"`
	QuantidadeMaxima  *int       `json:"quantidade_maxima"`
	Localizacao       *string    `gorm:"size:100" json:"localizacao"`
	Lote              *string    `gorm:"size:50" json:"lote"`
	DataValidade      *time.Time `gorm:"type:date" json:"data_validade"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	
	// Relacionamentos
	Produto           *ProdutoGorm `gorm:"foreignKey:ProdutoID" json:"produto,omitempty"`
}

// TableName especifica o nome da tabela
func (EstoqueGorm) TableName() string {
	return "estoque"
}

// MovimentacaoEstoqueGorm representa uma movimentação de estoque usando GORM
type MovimentacaoEstoqueGorm struct {
	ID                  string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	ProdutoID           string    `gorm:"type:uuid;not null" json:"produto_id"`
	TipoMovimentacao    string    `gorm:"size:20;not null;check:tipo_movimentacao IN ('entrada','saida','ajuste')" json:"tipo_movimentacao"`
	Quantidade          int       `gorm:"not null" json:"quantidade"`
	QuantidadeAnterior  int       `gorm:"not null" json:"quantidade_anterior"`
	QuantidadeNova      int       `gorm:"not null" json:"quantidade_nova"`
	Motivo              *string   `gorm:"size:255" json:"motivo"`
	Observacoes         *string   `gorm:"type:text" json:"observacoes"`
	Usuario             *string   `gorm:"size:100" json:"usuario"`
	CreatedAt           time.Time `json:"created_at"`
	
	// Relacionamentos
	Produto             *ProdutoGorm `gorm:"foreignKey:ProdutoID" json:"produto,omitempty"`
}

// TableName especifica o nome da tabela
func (MovimentacaoEstoqueGorm) TableName() string {
	return "movimentacoes_estoque"
}

// Inputs para criação/atualização (mantendo compatibilidade)
type CategoriaGormInput struct {
	Nome      string  `json:"nome" validate:"required,min=2,max=255"`
	Descricao *string `json:"descricao"`
	Ativo     *bool   `json:"ativo"`
}

type ProdutoGormInput struct {
	Nome           string   `json:"nome" validate:"required,min=2,max=255"`
	Descricao      *string  `json:"descricao"`
	CodigoBarras   *string  `json:"codigo_barras" validate:"omitempty,max=50"`
	CategoriaID    *string  `json:"categoria_id" validate:"omitempty,uuid"`
	PrecoCusto     *float64 `json:"preco_custo" validate:"omitempty,gte=0"`
	PrecoVenda     float64  `json:"preco_venda" validate:"required,gt=0"`
	MargemLucro    *float64 `json:"margem_lucro" validate:"omitempty,gte=0,lte=100"`
	UnidadeMedida  string   `json:"unidade_medida" validate:"required,max=10"`
	Peso           *float64 `json:"peso" validate:"omitempty,gte=0"`
	Dimensoes      *string  `json:"dimensoes" validate:"omitempty,max=50"`
	Ativo          *bool    `json:"ativo"`
	ImagemURL      *string  `json:"imagem_url" validate:"omitempty,url,max=500"`
	Observacoes    *string  `json:"observacoes"`
}