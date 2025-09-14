package models

import (
	"time"
)

// Estoque representa o estoque de um produto
type Estoque struct {
	ID                string     `json:"id" db:"id"`
	ProdutoID         string     `json:"produto_id" db:"produto_id"`
	QuantidadeAtual   int        `json:"quantidade_atual" db:"quantidade_atual"`
	QuantidadeMinima  int        `json:"quantidade_minima" db:"quantidade_minima"`
	QuantidadeMaxima  *int       `json:"quantidade_maxima" db:"quantidade_maxima"`
	Localizacao       *string    `json:"localizacao" db:"localizacao"`
	Lote              *string    `json:"lote" db:"lote"`
	DataValidade      *time.Time `json:"data_validade" db:"data_validade"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relacionamentos
	Produto           *Produto   `json:"produto,omitempty"`
}

// MovimentacaoEstoque representa uma movimentação de estoque
type MovimentacaoEstoque struct {
	ID                  string    `json:"id" db:"id"`
	ProdutoID           string    `json:"produto_id" db:"produto_id"`
	TipoMovimentacao    string    `json:"tipo_movimentacao" db:"tipo_movimentacao"`
	Quantidade          int       `json:"quantidade" db:"quantidade"`
	QuantidadeAnterior  int       `json:"quantidade_anterior" db:"quantidade_anterior"`
	QuantidadeNova      int       `json:"quantidade_nova" db:"quantidade_nova"`
	Motivo              *string   `json:"motivo" db:"motivo"`
	Observacoes         *string   `json:"observacoes" db:"observacoes"`
	Usuario             *string   `json:"usuario" db:"usuario"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	
	// Relacionamentos
	Produto             *Produto  `json:"produto,omitempty"`
}

// EstoqueInput representa os dados de entrada para criar/atualizar estoque
type EstoqueInput struct {
	ProdutoID         string     `json:"produto_id" validate:"required,uuid"`
	QuantidadeAtual   int        `json:"quantidade_atual" validate:"gte=0"`
	QuantidadeMinima  int        `json:"quantidade_minima" validate:"gte=0"`
	QuantidadeMaxima  *int       `json:"quantidade_maxima" validate:"omitempty,gte=0"`
	Localizacao       *string    `json:"localizacao" validate:"omitempty,max=100"`
	Lote              *string    `json:"lote" validate:"omitempty,max=50"`
	DataValidade      *time.Time `json:"data_validade"`
}

// MovimentacaoEstoqueInput representa os dados de entrada para movimentação
type MovimentacaoEstoqueInput struct {
	ProdutoID        string  `json:"produto_id" validate:"required,uuid"`
	TipoMovimentacao string  `json:"tipo_movimentacao" validate:"required,oneof=entrada saida ajuste"`
	Quantidade       int     `json:"quantidade" validate:"required,gt=0"`
	Motivo           *string `json:"motivo" validate:"omitempty,max=255"`
	Observacoes      *string `json:"observacoes"`
	Usuario          *string `json:"usuario" validate:"omitempty,max=100"`
}

// EstoqueResumo representa um resumo do estoque
type EstoqueResumo struct {
	TotalProdutos        int `json:"total_produtos"`
	ProdutosEstoqueBaixo int `json:"produtos_estoque_baixo"`
	ProdutosSemEstoque   int `json:"produtos_sem_estoque"`
	ValorTotalEstoque    float64 `json:"valor_total_estoque"`
}