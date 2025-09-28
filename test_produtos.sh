#!/bin/bash

# Script para testar o CRUD de produtos (GORM)
echo "🧪 Testando CRUD de Produtos (GORM)..."

# Verificar se o servidor está rodando
echo "📡 Verificando se o servidor está rodando..."
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Servidor não está rodando. Inicie com: go run main.go"
    exit 1
fi

echo "✅ Servidor está rodando!"

# Testar criação de categoria primeiro
echo ""
echo "📝 Testando criação de categoria..."
CATEGORIA_ID=$(curl -s -X POST http://localhost:8080/api/categorias \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Eletrônicos Teste",
    "descricao": "Categoria para testes"
  }' | jq -r '.data.id' 2>/dev/null)

if [ "$CATEGORIA_ID" != "null" ] && [ "$CATEGORIA_ID" != "" ]; then
    echo "✅ Categoria criada com ID: $CATEGORIA_ID"
else
    echo "⚠️ Categoria não criada, continuando sem categoria..."
    CATEGORIA_ID=""
fi

# Testar criação de produto
echo ""
echo "📝 Testando criação de produto..."
PRODUTO_JSON='{
    "nome": "Notebook Dell Teste",
    "descricao": "Notebook Dell Inspiron 15 para testes",
    "codigo_barras": "7891234567890",
    "preco_custo": 1500.00,
    "preco_venda": 2000.00,
    "margem_lucro": 33.33,
    "unidade_medida": "UN",
    "peso": 2.5,
    "dimensoes": "35x25x2 cm",
    "ativo": true
}'

if [ "$CATEGORIA_ID" != "" ]; then
    PRODUTO_JSON=$(echo $PRODUTO_JSON | jq --arg cat_id "$CATEGORIA_ID" '. + {categoria_id: $cat_id}')
fi

PRODUTO_ID=$(curl -s -X POST http://localhost:8080/api/produtos \
  -H "Content-Type: application/json" \
  -d "$PRODUTO_JSON" | jq -r '.data.id' 2>/dev/null)

if [ "$PRODUTO_ID" != "null" ] && [ "$PRODUTO_ID" != "" ]; then
    echo "✅ Produto criado com ID: $PRODUTO_ID"
else
    echo "❌ Falha ao criar produto"
    exit 1
fi

# Testar busca por ID
echo ""
echo "🔍 Testando busca por ID..."
NOME_PRODUTO=$(curl -s http://localhost:8080/api/produtos/$PRODUTO_ID | jq -r '.data.nome' 2>/dev/null)
if [ "$NOME_PRODUTO" = "Notebook Dell Teste" ]; then
    echo "✅ Produto encontrado: $NOME_PRODUTO"
else
    echo "❌ Falha ao buscar produto"
fi

# Testar listagem
echo ""
echo "📋 Testando listagem de produtos..."
TOTAL_PRODUTOS=$(curl -s http://localhost:8080/api/produtos | jq -r '.data.total' 2>/dev/null)
echo "✅ Total de produtos: $TOTAL_PRODUTOS"

# Testar busca
echo ""
echo "🔍 Testando busca por nome..."
PRODUTOS_ENCONTRADOS=$(curl -s "http://localhost:8080/api/produtos?search=Dell" | jq -r '.data.total' 2>/dev/null)
echo "✅ Produtos encontrados com 'Dell': $PRODUTOS_ENCONTRADOS"

# Testar atualização
echo ""
echo "✏️ Testando atualização de produto..."
curl -s -X PUT http://localhost:8080/api/produtos/$PRODUTO_ID \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Notebook Dell Atualizado",
    "preco_venda": 1800.00,
    "ativo": true
  }' > /dev/null

NOME_ATUALIZADO=$(curl -s http://localhost:8080/api/produtos/$PRODUTO_ID | jq -r '.data.nome' 2>/dev/null)
if [ "$NOME_ATUALIZADO" = "Notebook Dell Atualizado" ]; then
    echo "✅ Produto atualizado: $NOME_ATUALIZADO"
else
    echo "❌ Falha ao atualizar produto"
fi

# Testar deleção (soft delete)
echo ""
echo "🗑️ Testando deleção de produto (soft delete)..."
DELETE_RESULT=$(curl -s -X DELETE http://localhost:8080/api/produtos/$PRODUTO_ID | jq -r '.success' 2>/dev/null)
if [ "$DELETE_RESULT" = "true" ]; then
    echo "✅ Produto deletado com sucesso (soft delete)"
else
    echo "❌ Falha ao deletar produto"
fi

# Verificar se foi realmente deletado (soft delete)
echo ""
echo "🔍 Verificando se produto foi deletado (soft delete)..."
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/produtos/$PRODUTO_ID)
if [ "$STATUS_CODE" = "404" ]; then
    echo "✅ Produto não encontrado (soft delete funcionando)"
else
    echo "❌ Produto ainda existe (erro no soft delete)"
fi

# Limpar categoria se foi criada
if [ "$CATEGORIA_ID" != "" ]; then
    echo ""
    echo "🧹 Limpando categoria de teste..."
    curl -s -X DELETE http://localhost:8080/api/categorias/$CATEGORIA_ID > /dev/null
    echo "✅ Categoria removida"
fi

echo ""
echo "🎉 Testes de produtos (GORM) concluídos!"
echo ""
echo "🌐 Acesse a interface web em:"
echo "   - Home: http://localhost:8080"
echo "   - Produtos: http://localhost:8080/produtos"
echo "   - Novo Produto: http://localhost:8080/produtos/novo"
echo "   - Clientes (SQL): http://localhost:8080/clientes"
echo ""
echo "📊 Comparação disponível:"
echo "   - Produtos (GORM): Implementação moderna com ORM"
echo "   - Clientes (SQL): Implementação tradicional com SQL puro"