package models

import (
	"time"
)

// Cliente representa um cliente no sistema
type Cliente struct {
	ID          string    `json:"id" db:"id"`
	Nome        string    `json:"nome" db:"nome"`
	Email       *string   `json:"email" db:"email"`
	Telefone    *string   `json:"telefone" db:"telefone"`
	Endereco    *string   `json:"endereco" db:"endereco"`
	Cidade      *string   `json:"cidade" db:"cidade"`
	Estado      *string   `json:"estado" db:"estado"`
	CEP         *string   `json:"cep" db:"cep"`
	CPFCNPJ     *string   `json:"cpf_cnpj" db:"cpf_cnpj"`
	TipoPessoa  string    `json:"tipo_pessoa" db:"tipo_pessoa"`
	Ativo       bool      `json:"ativo" db:"ativo"`
	Observacoes *string   `json:"observacoes" db:"observacoes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ClienteInput representa os dados de entrada para criar/atualizar cliente
type ClienteInput struct {
	Nome        string  `json:"nome" validate:"required,min=2,max=255"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Telefone    *string `json:"telefone" validate:"omitempty,min=10,max=20"`
	Endereco    *string `json:"endereco"`
	Cidade      *string `json:"cidade" validate:"omitempty,max=100"`
	Estado      *string `json:"estado" validate:"omitempty,len=2"`
	CEP         *string `json:"cep" validate:"omitempty,len=8"`
	CPFCNPJ     *string `json:"cpf_cnpj" validate:"omitempty,max=20"`
	TipoPessoa  string  `json:"tipo_pessoa" validate:"required,oneof=fisica juridica"`
	Ativo       *bool   `json:"ativo"`
	Observacoes *string `json:"observacoes"`
}