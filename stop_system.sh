#!/bin/bash

# Script para parar o Sistema de Gestão
# Autor: Sistema de Gestão Go
# Versão: 1.0

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

echo -e "${BLUE}"
cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                    🛑 PARAR SISTEMA                          ║
║                                                              ║
║              Parando todos os serviços...                   ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Parar aplicação Go (agora roda no Docker)
info "Parando aplicação Go local (se houver)..."
if pkill -f "webserver-go"; then
    success "Aplicação Go local parada"
else
    info "Nenhuma aplicação Go local encontrada"
fi

# Detectar se precisa de sudo para Docker
DOCKER_COMPOSE_CMD="docker-compose"
if ! docker-compose ps &>/dev/null 2>&1; then
    if sudo docker-compose ps &>/dev/null 2>&1; then
        DOCKER_COMPOSE_CMD="sudo docker-compose"
    fi
fi

# Parar stack completa (Nginx + Go + PostgreSQL)
info "Parando stack completa (Nginx + Go + PostgreSQL)..."
if $DOCKER_COMPOSE_CMD down --remove-orphans; then
    success "Stack completa parada (Nginx + Go + PostgreSQL)"
else
    warning "Erro ao parar containers ou nenhum container rodando"
fi

# Verificar se ainda há processos
info "Verificando processos restantes..."
REMAINING=$(pgrep -f "webserver-go" | wc -l)
if [ "$REMAINING" -eq "0" ]; then
    success "Todos os processos foram parados"
else
    warning "$REMAINING processo(s) ainda rodando"
    info "Para forçar parada: pkill -9 -f webserver-go"
fi

# Verificar containers
info "Verificando containers..."
CONTAINERS=$($DOCKER_COMPOSE_CMD ps -q | wc -l)
if [ "$CONTAINERS" -eq "0" ]; then
    success "Todos os containers foram parados"
else
    warning "$CONTAINERS container(s) ainda rodando"
    info "Para forçar parada: $DOCKER_COMPOSE_CMD down -v"
fi

echo ""
success "Sistema parado com sucesso! 🛑"
echo ""
info "Para reiniciar: ./start_and_test.sh"