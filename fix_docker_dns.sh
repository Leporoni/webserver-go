#!/bin/bash

# Script para corrigir problema de DNS/IPv6 do Docker
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

error() {
    echo -e "${RED}❌ $1${NC}"
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
║                    🔧 CORREÇÃO DNS DOCKER                    ║
║                                                              ║
║              Resolvendo problema IPv6/DNS                   ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

info "Problema detectado: Docker não consegue acessar registry-1.docker.io via IPv6"
info "Solução: Configurar Docker para usar IPv4 e DNS específicos"
echo ""

# Backup da configuração atual (se existir)
if [ -f "/etc/docker/daemon.json" ]; then
    warning "Fazendo backup da configuração atual do Docker..."
    sudo cp /etc/docker/daemon.json /etc/docker/daemon.json.backup
    success "Backup criado: /etc/docker/daemon.json.backup"
fi

# Criar diretório se não existir
sudo mkdir -p /etc/docker

# Criar nova configuração do Docker
info "Criando nova configuração do Docker..."
sudo tee /etc/docker/daemon.json > /dev/null << 'EOF'
{
  "dns": ["8.8.8.8", "8.8.4.4", "1.1.1.1"],
  "ipv6": false,
  "ip6tables": false,
  "experimental": false,
  "live-restore": true,
  "userland-proxy": false,
  "registry-mirrors": []
}
EOF

success "Configuração do Docker criada"

# Mostrar a configuração
info "Configuração aplicada:"
cat /etc/docker/daemon.json | sed 's/^/  /'
echo ""

# Reiniciar Docker
info "Reiniciando Docker..."
if sudo systemctl restart docker; then
    success "Docker reiniciado com sucesso"
else
    error "Falha ao reiniciar Docker"
    exit 1
fi

# Aguardar Docker inicializar
info "Aguardando Docker inicializar..."
sleep 5

# Testar Docker
info "Testando Docker..."
if docker ps &>/dev/null; then
    success "Docker funcionando corretamente"
else
    error "Docker não está respondendo"
    info "Verifique os logs: sudo journalctl -u docker.service"
    exit 1
fi

# Testar conectividade com Docker Hub
info "Testando conectividade com Docker Hub..."
if docker pull hello-world &>/dev/null; then
    success "Conectividade com Docker Hub OK"
    
    # Limpar imagem de teste
    docker rmi hello-world &>/dev/null || true
    
    echo ""
    success "🎉 Problema resolvido!"
    info "Agora você pode executar: ./start_and_test.sh"
    
else
    error "Ainda há problemas de conectividade"
    
    echo ""
    warning "Soluções alternativas:"
    echo "1. Usar PostgreSQL local: ./start_local_postgres.sh"
    echo "2. Configurar proxy/VPN se estiver em rede corporativa"
    echo "3. Verificar firewall: sudo ufw status"
    
    # Testar IPv4 específico
    info "Testando conectividade IPv4 específica..."
    if ping -4 -c 1 registry-1.docker.io &>/dev/null; then
        success "IPv4 funciona - problema é configuração do Docker"
        info "Tente reiniciar o sistema: sudo reboot"
    else
        error "Problema de rede mais amplo"
        info "Verifique configurações de rede/proxy"
    fi
fi

echo ""
info "📋 Logs úteis para diagnóstico:"
echo "  • Docker: sudo journalctl -u docker.service"
echo "  • Sistema: sudo journalctl -xe"
echo "  • Rede: ip route show"