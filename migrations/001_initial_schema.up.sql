-- Initial schema for gestão system
-- Created: $(date)

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Clientes table
CREATE TABLE IF NOT EXISTS clientes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    telefone VARCHAR(20),
    endereco TEXT,
    cidade VARCHAR(100),
    estado VARCHAR(2),
    cep VARCHAR(10),
    cpf_cnpj VARCHAR(20) UNIQUE,
    tipo_pessoa VARCHAR(10) CHECK (tipo_pessoa IN ('fisica', 'juridica')) DEFAULT 'fisica',
    ativo BOOLEAN DEFAULT true,
    observacoes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Categorias table
CREATE TABLE IF NOT EXISTS categorias (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL UNIQUE,
    descricao TEXT,
    ativo BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Produtos table
CREATE TABLE IF NOT EXISTS produtos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    codigo_barras VARCHAR(50) UNIQUE,
    categoria_id UUID REFERENCES categorias(id),
    preco_custo DECIMAL(10,2),
    preco_venda DECIMAL(10,2) NOT NULL,
    margem_lucro DECIMAL(5,2),
    unidade_medida VARCHAR(10) DEFAULT 'UN',
    peso DECIMAL(8,3),
    dimensoes VARCHAR(50),
    ativo BOOLEAN DEFAULT true,
    imagem_url VARCHAR(500),
    observacoes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Estoque table
CREATE TABLE IF NOT EXISTS estoque (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    produto_id UUID REFERENCES produtos(id) ON DELETE CASCADE,
    quantidade_atual INTEGER NOT NULL DEFAULT 0,
    quantidade_minima INTEGER DEFAULT 0,
    quantidade_maxima INTEGER,
    localizacao VARCHAR(100),
    lote VARCHAR(50),
    data_validade DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(produto_id, lote)
);

-- Movimentações de estoque table
CREATE TABLE IF NOT EXISTS movimentacoes_estoque (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    produto_id UUID REFERENCES produtos(id),
    tipo_movimentacao VARCHAR(20) CHECK (tipo_movimentacao IN ('entrada', 'saida', 'ajuste')) NOT NULL,
    quantidade INTEGER NOT NULL,
    quantidade_anterior INTEGER NOT NULL,
    quantidade_nova INTEGER NOT NULL,
    motivo VARCHAR(255),
    observacoes TEXT,
    usuario VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_clientes_nome ON clientes(nome);
CREATE INDEX idx_clientes_email ON clientes(email);
CREATE INDEX idx_clientes_cpf_cnpj ON clientes(cpf_cnpj);
CREATE INDEX idx_produtos_nome ON produtos(nome);
CREATE INDEX idx_produtos_codigo_barras ON produtos(codigo_barras);
CREATE INDEX idx_produtos_categoria ON produtos(categoria_id);
CREATE INDEX idx_estoque_produto ON estoque(produto_id);
CREATE INDEX idx_movimentacoes_produto ON movimentacoes_estoque(produto_id);
CREATE INDEX idx_movimentacoes_data ON movimentacoes_estoque(created_at);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_clientes_updated_at BEFORE UPDATE ON clientes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_categorias_updated_at BEFORE UPDATE ON categorias FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_produtos_updated_at BEFORE UPDATE ON produtos FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_estoque_updated_at BEFORE UPDATE ON estoque FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();