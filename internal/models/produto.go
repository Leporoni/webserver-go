package models

import (
	"time"
)

// Categoria representa uma categoria de produto
type Categoria struct {
	ID        string    `json:"id" db:"id"`
	Nome      string    `json:"nome" db:"nome"`
	Descricao *string   `json:"descricao" db:"descricao"`
	Ativo     bool      `json:"ativo" db:"ativo"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Produto representa um produto no sistema
type Produto struct {
	ID             string     `json:"id" db:"id"`
	Nome           string     `json:"nome" db:"nome"`
	Descricao      *string    `json:"descricao" db:"descricao"`
	CodigoBarras   *string    `json:"codigo_barras" db:"codigo_barras"`
	CategoriaID    *string    `json:"categoria_id" db:"categoria_id"`
	PrecoCusto     *float64   `json:"preco_custo" db:"preco_custo"`
	PrecoVenda     float64    `json:"preco_venda" db:"preco_venda"`
	MargemLucro    *float64   `json:"margem_lucro" db:"margem_lucro"`
	UnidadeMedida  string     `json:"unidade_medida" db:"unidade_medida"`
	Peso           *float64   `json:"peso" db:"peso"`
	Dimensoes      *string    `json:"dimensoes" db:"dimensoes"`
	Ativo          bool       `json:"ativo" db:"ativo"`
	ImagemURL      *string    `json:"imagem_url" db:"imagem_url"`
	Observacoes    *string    `json:"observacoes" db:"observacoes"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relacionamentos
	Categoria      *Categoria `json:"categoria,omitempty"`
}

// ProdutoInput representa os dados de entrada para criar/atualizar produto
type ProdutoInput struct {
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

// CategoriaInput representa os dados de entrada para criar/atualizar categoria
type CategoriaInput struct {
	Nome      string  `json:"nome" validate:"required,min=2,max=255"`
	Descricao *string `json:"descricao"`
	Ativo     *bool   `json:"ativo"`
}