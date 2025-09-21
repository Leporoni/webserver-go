#!/bin/bash

# Script de teste rápido
# Testa se o sistema está funcionando

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo "🧪 Teste Rápido do Sistema"
echo ""

# Testar se PostgreSQL está rodando
info "Verificando PostgreSQL..."
if docker ps | grep postgres &>/dev/null; then
    success "PostgreSQL rodando via Docker"
    
    # Testar conexão
    if PGPASSWORD=gestao_pass psql -h 127.0.0.1 -U gestao_user -d gestao_db -c "SELECT 1;" &>/dev/null; then
        success "Conexão com banco OK"
    else
        error "Falha na conexão com banco"
    fi
else
    error "PostgreSQL não está rodando"
    info "Execute: docker-compose up -d postgres"
fi

# Testar se aplicação compila
info "Testando compilação..."
if go build -o test-build . &>/dev/null; then
    success "Compilação OK"
    rm -f test-build
else
    error "Falha na compilação"
fi

# Testar se aplicação inicia (teste rápido)
info "Testando inicialização da aplicação..."
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_NAME=gestao_db
export DB_USER=gestao_user
export DB_PASSWORD=gestao_pass

# Compilar se necessário
go build -o webserver-go . &>/dev/null

# Iniciar aplicação em background por 5 segundos
timeout 5s ./webserver-go &
APP_PID=$!

sleep 3

# Testar se está respondendo
if curl -s http://localhost:8080/health &>/dev/null; then
    success "Aplicação respondendo"
    
    # Testar endpoint de clientes
    if curl -s http://localhost:8080/api/clientes &>/dev/null; then
        success "API de clientes funcionando"
    else
        error "API de clientes não responde"
    fi
else
    error "Aplicação não está respondendo"
fi

# Parar aplicação
kill $APP_PID 2>/dev/null || true
wait $APP_PID 2>/dev/null || true

echo ""
echo "🎯 Para executar o sistema completo:"
echo "  ./run_with_docker_postgres.sh"
echo ""
echo "🌐 Acessos após iniciar:"
echo "  • http://localhost:8080"
echo "  • http://localhost:8080/clientes"