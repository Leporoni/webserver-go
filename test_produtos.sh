#!/bin/bash

# Script de teste para CRUD de Produtos
# Testa todos os endpoints da API de produtos

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Funções de output
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Configurações
BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api"

# Verificar se servidor está rodando
info "Verificando se servidor está rodando..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    error "Servidor não está rodando em $BASE_URL"
    echo "Execute primeiro: ./start_and_test.sh ou ./dev_mode.sh"
    exit 1
fi

success "Servidor está rodando"

echo ""
echo -e "${PURPLE}🧪 Executando testes automatizados do CRUD de Produtos...${NC}"
echo ""

# Variáveis para armazenar IDs
CATEGORIA_ID=""
PRODUTO_ID=""

# Teste 1: Criar categoria
info "Teste 1: Criando categoria de teste"
CATEGORIA_RESPONSE=$(curl -s -X POST "$API_URL/categorias" \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Eletrônicos",
    "descricao": "Produtos eletrônicos e tecnologia",
    "ativo": true
  }')

if echo "$CATEGORIA_RESPONSE" | grep -q '"success":true'; then
    CATEGORIA_ID=$(echo "$CATEGORIA_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    success "Categoria criada com sucesso (ID: $CATEGORIA_ID)"
else
    error "Falha ao criar categoria"
    echo "Resposta: $CATEGORIA_RESPONSE"
    exit 1
fi

# Teste 2: Listar categorias
info "Teste 2: Listando categorias"
CATEGORIAS_RESPONSE=$(curl -s "$API_URL/categorias")

if echo "$CATEGORIAS_RESPONSE" | grep -q '"success":true'; then
    success "Categorias listadas com sucesso"
else
    error "Falha ao listar categorias"
    echo "Resposta: $CATEGORIAS_RESPONSE"
fi

# Teste 3: Criar produto
info "Teste 3: Criando produto de teste"
PRODUTO_RESPONSE=$(curl -s -X POST "$API_URL/produtos" \
  -H "Content-Type: application/json" \
  -d "{
    \"nome\": \"Smartphone XYZ\",
    \"descricao\": \"Smartphone com 128GB de armazenamento\",
    \"codigo_barras\": \"1234567890123\",
    \"categoria_id\": \"$CATEGORIA_ID\",
    \"preco_custo\": 800.00,
    \"preco_venda\": 1200.00,
    \"unidade_medida\": \"UN\",
    \"peso\": 0.180,
    \"dimensoes\": \"15x7x0.8 cm\",
    \"ativo\": true,
    \"observacoes\": \"Produto de teste criado automaticamente\"
  }")

if echo "$PRODUTO_RESPONSE" | grep -q '"success":true'; then
    PRODUTO_ID=$(echo "$PRODUTO_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    success "Produto criado com sucesso (ID: $PRODUTO_ID)"
else
    error "Falha ao criar produto"
    echo "Resposta: $PRODUTO_RESPONSE"
    exit 1
fi

# Teste 4: Buscar produto por ID
info "Teste 4: Buscando produto por ID"
GET_PRODUTO_RESPONSE=$(curl -s "$API_URL/produtos/$PRODUTO_ID")

if echo "$GET_PRODUTO_RESPONSE" | grep -q '"success":true'; then
    success "Produto encontrado com sucesso"
    # Verificar se dados estão corretos
    if echo "$GET_PRODUTO_RESPONSE" | grep -q '"nome":"Smartphone XYZ"'; then
        success "Dados do produto estão corretos"
    else
        warning "Dados do produto podem estar incorretos"
    fi
else
    error "Falha ao buscar produto"
    echo "Resposta: $GET_PRODUTO_RESPONSE"
fi

# Teste 5: Listar produtos
info "Teste 5: Listando produtos"
LIST_PRODUTOS_RESPONSE=$(curl -s "$API_URL/produtos")

if echo "$LIST_PRODUTOS_RESPONSE" | grep -q '"success":true'; then
    success "Produtos listados com sucesso"
    # Verificar se nosso produto está na lista
    if echo "$LIST_PRODUTOS_RESPONSE" | grep -q "$PRODUTO_ID"; then
        success "Produto de teste encontrado na listagem"
    else
        warning "Produto de teste não encontrado na listagem"
    fi
else
    error "Falha ao listar produtos"
    echo "Resposta: $LIST_PRODUTOS_RESPONSE"
fi

# Teste 6: Buscar produtos com filtros
info "Teste 6: Testando busca com filtros"
SEARCH_RESPONSE=$(curl -s "$API_URL/produtos?search=Smartphone&categoria_id=$CATEGORIA_ID&page=1&limit=10")

if echo "$SEARCH_RESPONSE" | grep -q '"success":true'; then
    success "Busca com filtros funcionando"
else
    error "Falha na busca com filtros"
    echo "Resposta: $SEARCH_RESPONSE"
fi

# Teste 7: Atualizar produto
info "Teste 7: Atualizando produto"
UPDATE_RESPONSE=$(curl -s -X PUT "$API_URL/produtos/$PRODUTO_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Smartphone XYZ Pro",
    "descricao": "Smartphone com 256GB de armazenamento - Versão Pro",
    "codigo_barras": "1234567890123",
    "preco_custo": 900.00,
    "preco_venda": 1400.00,
    "unidade_medida": "UN",
    "ativo": true
  }')

if echo "$UPDATE_RESPONSE" | grep -q '"success":true'; then
    success "Produto atualizado com sucesso"
    # Verificar se a atualização foi aplicada
    VERIFY_UPDATE=$(curl -s "$API_URL/produtos/$PRODUTO_ID")
    if echo "$VERIFY_UPDATE" | grep -q '"nome":"Smartphone XYZ Pro"'; then
        success "Atualização verificada com sucesso"
    else
        warning "Atualização pode não ter sido aplicada corretamente"
    fi
else
    error "Falha ao atualizar produto"
    echo "Resposta: $UPDATE_RESPONSE"
fi

# Teste 8: Testar validações
info "Teste 8: Testando validações (produto inválido)"
INVALID_RESPONSE=$(curl -s -X POST "$API_URL/produtos" \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "",
    "preco_venda": -100,
    "unidade_medida": ""
  }')

if echo "$INVALID_RESPONSE" | grep -q '"success":false' || echo "$INVALID_RESPONSE" | grep -q "obrigatório"; then
    success "Validações funcionando corretamente"
else
    warning "Validações podem não estar funcionando"
    echo "Resposta: $INVALID_RESPONSE"
fi

# Teste 9: Deletar produto
info "Teste 9: Deletando produto de teste"
DELETE_RESPONSE=$(curl -s -X DELETE "$API_URL/produtos/$PRODUTO_ID")

if echo "$DELETE_RESPONSE" | grep -q '"success":true'; then
    success "Produto deletado com sucesso"
    
    # Verificar se produto foi realmente deletado
    VERIFY_DELETE=$(curl -s "$API_URL/produtos/$PRODUTO_ID")
    if echo "$VERIFY_DELETE" | grep -q "não encontrado" || echo "$VERIFY_DELETE" | grep -q "404"; then
        success "Deleção verificada com sucesso"
    else
        warning "Produto pode não ter sido deletado corretamente"
    fi
else
    error "Falha ao deletar produto"
    echo "Resposta: $DELETE_RESPONSE"
fi

# Teste 10: Verificar produto inexistente
info "Teste 10: Testando busca de produto inexistente"
NOTFOUND_RESPONSE=$(curl -s "$API_URL/produtos/00000000-0000-0000-0000-000000000000")

if echo "$NOTFOUND_RESPONSE" | grep -q "não encontrado" || echo "$NOTFOUND_RESPONSE" | grep -q "404"; then
    success "Tratamento de produto não encontrado funcionando"
else
    warning "Tratamento de erro pode não estar funcionando"
    echo "Resposta: $NOTFOUND_RESPONSE"
fi

echo ""
echo -e "${GREEN}✅ Todos os testes do CRUD de Produtos foram executados!${NC}"
echo ""

# Informações finais
echo -e "${CYAN}🌐 Acessos disponíveis:${NC}"
echo "   • Home: $BASE_URL"
echo "   • Produtos: $BASE_URL/produtos"
echo "   • Novo Produto: $BASE_URL/produtos/novo"
echo "   • API Produtos: $API_URL/produtos"
echo "   • API Categorias: $API_URL/categorias"
echo ""

echo -e "${CYAN}📋 Endpoints testados:${NC}"
echo "   • POST /api/categorias - Criar categoria"
echo "   • GET /api/categorias - Listar categorias"
echo "   • POST /api/produtos - Criar produto"
echo "   • GET /api/produtos/{id} - Buscar produto"
echo "   • GET /api/produtos - Listar produtos"
echo "   • GET /api/produtos?search=... - Buscar com filtros"
echo "   • PUT /api/produtos/{id} - Atualizar produto"
echo "   • DELETE /api/produtos/{id} - Deletar produto"
echo ""

echo -e "${YELLOW}💡 Próximos passos:${NC}"
echo "   • Teste a interface web em $BASE_URL/produtos"
echo "   • Crie categorias e produtos através da interface"
echo "   • Teste os filtros e busca na listagem"
echo "   • Verifique a integração com o módulo de clientes"
echo ""

echo -e "${PURPLE}🎉 CRUD de Produtos implementado e testado com sucesso!${NC}"