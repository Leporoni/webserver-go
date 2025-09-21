#!/bin/bash

# Script para testar o CRUD de clientes
echo "🧪 Testando CRUD de Clientes..."

# Verificar se o servidor está rodando
echo "📡 Verificando se o servidor está rodando..."
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Servidor não está rodando. Inicie com: go run main.go"
    exit 1
fi

echo "✅ Servidor está rodando!"

# Testar criação de cliente
echo ""
echo "📝 Testando criação de cliente..."
CLIENTE_ID=$(curl -s -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva Teste",
    "email": "joao.teste@email.com",
    "telefone": "(11) 99999-9999",
    "tipo_pessoa": "fisica",
    "endereco": "Rua das Flores, 123",
    "cidade": "São Paulo",
    "estado": "SP",
    "cep": "01234567"
  }' | jq -r '.data.id')

if [ "$CLIENTE_ID" != "null" ] && [ "$CLIENTE_ID" != "" ]; then
    echo "✅ Cliente criado com ID: $CLIENTE_ID"
else
    echo "❌ Falha ao criar cliente"
    exit 1
fi

# Testar busca por ID
echo ""
echo "🔍 Testando busca por ID..."
NOME_CLIENTE=$(curl -s http://localhost:8080/api/clientes/$CLIENTE_ID | jq -r '.data.nome')
if [ "$NOME_CLIENTE" = "João Silva Teste" ]; then
    echo "✅ Cliente encontrado: $NOME_CLIENTE"
else
    echo "❌ Falha ao buscar cliente"
fi

# Testar listagem
echo ""
echo "📋 Testando listagem de clientes..."
TOTAL_CLIENTES=$(curl -s http://localhost:8080/api/clientes | jq -r '.data.total')
echo "✅ Total de clientes: $TOTAL_CLIENTES"

# Testar atualização
echo ""
echo "✏️ Testando atualização de cliente..."
curl -s -X PUT http://localhost:8080/api/clientes/$CLIENTE_ID \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva Atualizado",
    "email": "joao.atualizado@email.com",
    "tipo_pessoa": "fisica"
  }' > /dev/null

NOME_ATUALIZADO=$(curl -s http://localhost:8080/api/clientes/$CLIENTE_ID | jq -r '.data.nome')
if [ "$NOME_ATUALIZADO" = "João Silva Atualizado" ]; then
    echo "✅ Cliente atualizado: $NOME_ATUALIZADO"
else
    echo "❌ Falha ao atualizar cliente"
fi

# Testar deleção
echo ""
echo "🗑️ Testando deleção de cliente..."
DELETE_RESULT=$(curl -s -X DELETE http://localhost:8080/api/clientes/$CLIENTE_ID | jq -r '.success')
if [ "$DELETE_RESULT" = "true" ]; then
    echo "✅ Cliente deletado com sucesso"
else
    echo "❌ Falha ao deletar cliente"
fi

# Verificar se foi realmente deletado
echo ""
echo "🔍 Verificando se cliente foi deletado..."
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/clientes/$CLIENTE_ID)
if [ "$STATUS_CODE" = "404" ]; then
    echo "✅ Cliente não encontrado (deletado corretamente)"
else
    echo "❌ Cliente ainda existe (erro na deleção)"
fi

echo ""
echo "🎉 Testes concluídos!"
echo ""
echo "🌐 Acesse a interface web em:"
echo "   - Home: http://localhost:8080"
echo "   - Clientes: http://localhost:8080/clientes"
echo "   - Novo Cliente: http://localhost:8080/clientes/novo"