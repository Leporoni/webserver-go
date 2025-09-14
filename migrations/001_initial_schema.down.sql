-- Rollback initial schema

-- Drop triggers
DROP TRIGGER IF EXISTS update_estoque_updated_at ON estoque;
DROP TRIGGER IF EXISTS update_produtos_updated_at ON produtos;
DROP TRIGGER IF EXISTS update_categorias_updated_at ON categorias;
DROP TRIGGER IF EXISTS update_clientes_updated_at ON clientes;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_movimentacoes_data;
DROP INDEX IF EXISTS idx_movimentacoes_produto;
DROP INDEX IF EXISTS idx_estoque_produto;
DROP INDEX IF EXISTS idx_produtos_categoria;
DROP INDEX IF EXISTS idx_produtos_codigo_barras;
DROP INDEX IF EXISTS idx_produtos_nome;
DROP INDEX IF EXISTS idx_clientes_cpf_cnpj;
DROP INDEX IF EXISTS idx_clientes_email;
DROP INDEX IF EXISTS idx_clientes_nome;

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS movimentacoes_estoque;
DROP TABLE IF EXISTS estoque;
DROP TABLE IF EXISTS produtos;
DROP TABLE IF EXISTS categorias;
DROP TABLE IF EXISTS clientes;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";