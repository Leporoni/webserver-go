-- Seed default categories
-- Requires extension uuid-ossp (for uuid_generate_v4). If not installed, run: CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
    -- Ensure extension exists
    PERFORM 1 FROM pg_extension WHERE extname = 'uuid-ossp';
    IF NOT FOUND THEN
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    END IF;
END $$;

INSERT INTO categorias (id, nome, descricao, ativo, created_at, updated_at) VALUES
    (uuid_generate_v4(), 'Eletrônicos', 'Dispositivos e acessórios', true, NOW(), NOW()),
    (uuid_generate_v4(), 'Acessórios', 'Itens complementares', true, NOW(), NOW()),
    (uuid_generate_v4(), 'Vestuário', 'Roupas e moda', true, NOW(), NOW()),
    (uuid_generate_v4(), 'Alimentos', 'Comidas e bebidas', true, NOW(), NOW()),
    (uuid_generate_v4(), 'Casa e Jardim', 'Itens para casa e jardim', true, NOW(), NOW());
