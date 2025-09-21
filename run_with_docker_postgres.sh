#!/bin/bash

# Script para executar com PostgreSQL via Docker
# Corrige problema de IPv6

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo -e "${PURPLE}"
cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                    🐳 DOCKER POSTGRES                        ║
║                                                              ║
║              Executando com PostgreSQL via Docker           ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Verificar se PostgreSQL está rodando
info "Verificando PostgreSQL..."
if docker ps | grep postgres &>/dev/null; then
    success "PostgreSQL está rodando via Docker"
else
    error "PostgreSQL não está rodando"
    info "Execute primeiro: ./start_and_test.sh"
    exit 1
fi

# Configurar variáveis de ambiente para IPv4
export DB_HOST=127.0.0.1  # Usar IPv4 explicitamente
export DB_PORT=5432
export DB_NAME=gestao_db
export DB_USER=gestao_user
export DB_PASSWORD=gestao_pass

echo ""
info "🌐 Configuração do banco:"
echo "  Host: $DB_HOST (IPv4 explícito)"
echo "  Porta: $DB_PORT"
echo "  Banco: $DB_NAME"
echo "  Usuário: $DB_USER"
echo ""

# Testar conexão
info "Testando conexão com banco..."
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;" &>/dev/null; then
    success "Conexão com banco OK"
else
    warning "Testando conexão direta..."
    # Aguardar um pouco mais
    sleep 5
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;" &>/dev/null; then
        success "Conexão com banco OK (após aguardar)"
    else
        error "Falha na conexão com banco"
        info "Verifique se o PostgreSQL está realmente pronto"
        exit 1
    fi
fi

# Compilar aplicação se necessário
if [ ! -f "webserver-go" ] || [ "main.go" -nt "webserver-go" ]; then
    info "Compilando aplicação..."
    if go build -o webserver-go .; then
        success "Aplicação compilada"
    else
        error "Falha na compilação"
        exit 1
    fi
fi

echo ""
success "🚀 Iniciando Sistema de Gestão..."
echo ""
info "Acessos disponíveis:"
echo "  • Home: http://localhost:8080"
echo "  • Clientes: http://localhost:8080/clientes"
echo "  • API: http://localhost:8080/api/clientes"
echo "  • Health: http://localhost:8080/health"
echo ""
warning "Para parar: Ctrl+C"
echo ""

# Iniciar aplicação com variáveis corretas
./webserver-go