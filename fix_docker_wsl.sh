#!/bin/bash

# Script para corrigir Docker em WSL/sistemas sem systemd
# Autor: Sistema de Gestão Go
# Versão: 1.0

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
║                    🔧 CORREÇÃO DOCKER WSL                    ║
║                                                              ║
║              Solução para WSL/sistemas sem systemd          ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Detectar ambiente
if grep -qi microsoft /proc/version; then
    info "WSL detectado"
    ENV_TYPE="wsl"
elif [ ! -d "/run/systemd/system" ]; then
    info "Sistema sem systemd detectado"
    ENV_TYPE="no-systemd"
else
    info "Sistema com systemd"
    ENV_TYPE="systemd"
fi

echo ""

# Verificar se Docker Desktop está rodando (WSL)
if [ "$ENV_TYPE" = "wsl" ]; then
    info "Verificando Docker Desktop no Windows..."
    
    if docker version &>/dev/null; then
        success "Docker Desktop está funcionando"
        
        # Testar conectividade
        info "Testando conectividade com Docker Hub..."
        if timeout 30 docker pull hello-world &>/dev/null; then
            success "Conectividade OK - problema resolvido!"
            docker rmi hello-world &>/dev/null || true
            
            echo ""
            success "🎉 Docker funcionando!"
            info "Execute: ./start_and_test.sh"
            exit 0
        else
            warning "Ainda há problemas de conectividade"
        fi
    else
        error "Docker Desktop não está rodando no Windows"
        echo ""
        warning "Soluções para WSL:"
        echo "1. Abra Docker Desktop no Windows"
        echo "2. Verifique se WSL integration está habilitada"
        echo "3. Reinicie Docker Desktop"
        echo "4. Use alternativa local: ./start_local_postgres.sh"
        exit 1
    fi
fi

# Para sistemas sem systemd, tentar reiniciar Docker manualmente
if [ "$ENV_TYPE" = "no-systemd" ]; then
    warning "Sistema sem systemd - tentando reiniciar Docker manualmente..."
    
    # Tentar parar Docker
    sudo pkill dockerd 2>/dev/null || true
    sudo pkill docker-containerd 2>/dev/null || true
    sleep 2
    
    # Tentar iniciar Docker
    info "Iniciando Docker daemon..."
    sudo dockerd --config-file=/etc/docker/daemon.json &
    DOCKER_PID=$!
    
    # Aguardar Docker inicializar
    info "Aguardando Docker inicializar..."
    for i in {1..30}; do
        if docker version &>/dev/null; then
            success "Docker iniciado"
            break
        fi
        echo -n "."
        sleep 1
    done
    echo ""
    
    if ! docker version &>/dev/null; then
        error "Docker não iniciou corretamente"
        sudo kill $DOCKER_PID 2>/dev/null || true
        exit 1
    fi
fi

# Testar conectividade final
info "Testando conectividade final..."
if timeout 30 docker pull hello-world &>/dev/null; then
    success "Conectividade com Docker Hub OK"
    docker rmi hello-world &>/dev/null || true
    
    echo ""
    success "🎉 Problema resolvido!"
    info "Execute: ./start_and_test.sh"
    
else
    error "Ainda há problemas de conectividade"
    
    echo ""
    warning "Soluções alternativas:"
    echo "1. Usar PostgreSQL local: ./start_local_postgres.sh"
    echo "2. Verificar proxy/firewall corporativo"
    echo "3. Tentar VPN diferente"
    echo "4. Usar hotspot do celular para testar"
    
    # Diagnóstico adicional
    echo ""
    info "📋 Diagnóstico adicional:"
    echo "Testando conectividade IPv4..."
    if ping -4 -c 1 8.8.8.8 &>/dev/null; then
        success "IPv4 para 8.8.8.8 OK"
    else
        error "Problema de conectividade IPv4"
    fi
    
    if ping -4 -c 1 registry-1.docker.io &>/dev/null; then
        success "IPv4 para Docker Hub OK"
    else
        error "Docker Hub inacessível via IPv4"
    fi
    
    echo ""
    info "Configuração de rede:"
    echo "DNS atual:"
    cat /etc/resolv.conf | grep nameserver | sed 's/^/  /'
    
    echo ""
    info "Rotas de rede:"
    ip route show | head -3 | sed 's/^/  /'
fi