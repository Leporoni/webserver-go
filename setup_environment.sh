#!/bin/bash

# Script para configurar o ambiente do Sistema de Gestão
# Autor: Sistema de Gestão Go
# Versão: 1.0

set -e  # Parar em caso de erro

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Função para log colorido
log() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')]${NC} $1"
}

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
    echo -e "${CYAN}ℹ️  $1${NC}"
}

# Banner
echo -e "${PURPLE}"
cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                    🔧 CONFIGURAÇÃO DO AMBIENTE               ║
║                                                              ║
║              Sistema de Gestão - Setup v1.0                 ║
║                                                              ║
║  • Verificação de dependências                              ║
║  • Configuração do Docker                                   ║
║  • Preparação do ambiente                                   ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Verificar se é root
if [ "$EUID" -eq 0 ]; then
    error "Não execute este script como root (sudo)"
    info "Execute como usuário normal: ./setup_environment.sh"
    exit 1
fi

log "Verificando sistema operacional..."
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    success "Sistema Linux detectado"
    OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    success "Sistema macOS detectado"
    OS="macos"
else
    warning "Sistema operacional não reconhecido: $OSTYPE"
    OS="unknown"
fi

# Verificar dependências
log "Verificando dependências..."

# Verificar Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    success "Go $GO_VERSION instalado"
else
    error "Go não está instalado"
    info "Instale Go: https://golang.org/doc/install"
    if [ "$OS" = "linux" ]; then
        info "Ubuntu/Debian: sudo apt install golang-go"
        info "Ou baixe de: https://golang.org/dl/"
    fi
    exit 1
fi

# Verificar Docker
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    success "Docker $DOCKER_VERSION instalado"
    
    # Verificar se Docker está rodando
    if ! docker ps &>/dev/null; then
        if sudo docker ps &>/dev/null; then
            warning "Docker requer sudo"
            NEEDS_DOCKER_SETUP=true
        else
            warning "Docker não está rodando - tentando iniciar..."
            if sudo systemctl start docker; then
                success "Docker iniciado"
                sleep 2
                if ! docker ps &>/dev/null && sudo docker ps &>/dev/null; then
                    NEEDS_DOCKER_SETUP=true
                fi
            else
                error "Falha ao iniciar Docker"
                exit 1
            fi
        fi
    else
        success "Docker funcionando corretamente"
    fi
else
    error "Docker não está instalado"
    info "Instale Docker: https://docs.docker.com/get-docker/"
    if [ "$OS" = "linux" ]; then
        info "Ubuntu/Debian:"
        info "  curl -fsSL https://get.docker.com -o get-docker.sh"
        info "  sudo sh get-docker.sh"
    fi
    exit 1
fi

# Verificar Docker Compose
if command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(docker-compose --version | awk '{print $3}' | sed 's/,//')
    success "Docker Compose $COMPOSE_VERSION instalado"
else
    error "Docker Compose não está instalado"
    info "Instale Docker Compose: https://docs.docker.com/compose/install/"
    if [ "$OS" = "linux" ]; then
        info "Ubuntu/Debian: sudo apt install docker-compose"
    fi
    exit 1
fi

# Configurar Docker se necessário
if [ "$NEEDS_DOCKER_SETUP" = true ]; then
    log "Configurando Docker para o usuário atual..."
    
    warning "Docker requer sudo. Configurando permissões..."
    info "Adicionando usuário $USER ao grupo docker..."
    
    if sudo usermod -aG docker $USER; then
        success "Usuário adicionado ao grupo docker"
        warning "IMPORTANTE: Você precisa fazer logout e login novamente"
        warning "Ou execute: newgrp docker"
        info "Após isso, o Docker funcionará sem sudo"
    else
        error "Falha ao adicionar usuário ao grupo docker"
    fi
    
    # Habilitar Docker para iniciar automaticamente
    if sudo systemctl enable docker; then
        success "Docker configurado para iniciar automaticamente"
    else
        warning "Falha ao configurar inicialização automática do Docker"
    fi
fi

# Verificar curl
if command -v curl &> /dev/null; then
    success "curl instalado"
else
    warning "curl não está instalado (necessário para testes)"
    if [ "$OS" = "linux" ]; then
        info "Instale com: sudo apt install curl"
    fi
fi

# Verificar jq (opcional)
if command -v jq &> /dev/null; then
    success "jq instalado (para testes JSON)"
else
    info "jq não está instalado (opcional, melhora os testes)"
    if [ "$OS" = "linux" ]; then
        info "Instale com: sudo apt install jq"
    fi
fi

# Verificar estrutura do projeto
log "Verificando estrutura do projeto..."

if [ ! -f "main.go" ]; then
    error "main.go não encontrado - execute no diretório raiz do projeto"
    exit 1
fi

if [ ! -f "docker-compose.yml" ]; then
    error "docker-compose.yml não encontrado"
    exit 1
fi

if [ ! -f "go.mod" ]; then
    error "go.mod não encontrado"
    exit 1
fi

success "Estrutura do projeto verificada"

# Verificar e baixar dependências Go
log "Verificando dependências Go..."
if go mod tidy; then
    success "Dependências Go atualizadas"
else
    error "Falha ao atualizar dependências Go"
    exit 1
fi

# Testar compilação
log "Testando compilação..."
if go build -o webserver-go-test .; then
    success "Compilação bem-sucedida"
    rm -f webserver-go-test
else
    error "Falha na compilação"
    exit 1
fi

# Dar permissões aos scripts
log "Configurando permissões dos scripts..."
chmod +x *.sh 2>/dev/null || true
success "Permissões dos scripts configuradas"

# Resumo final
echo ""
echo -e "${PURPLE}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${PURPLE}║                        ✅ CONFIGURAÇÃO CONCLUÍDA             ║${NC}"
echo -e "${PURPLE}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""

success "Ambiente configurado com sucesso!"
echo ""

if [ "$NEEDS_DOCKER_SETUP" = true ]; then
    warning "AÇÃO NECESSÁRIA:"
    echo "   1. Execute: newgrp docker"
    echo "   2. Ou faça logout/login"
    echo "   3. Depois execute: ./start_and_test.sh"
    echo ""
else
    info "🚀 Pronto para usar!"
    echo "   Execute: ./start_and_test.sh"
    echo ""
fi

info "📚 Scripts disponíveis:"
echo "   • ./start_and_test.sh - Inicializar sistema completo"
echo "   • ./dev_mode.sh - Modo desenvolvimento"
echo "   • ./stop_system.sh - Parar sistema"
echo "   • ./test_clientes.sh - Testar CRUD"
echo "   • ./help.sh - Ajuda completa"
echo ""

info "🌐 Após inicializar, acesse:"
echo "   • http://localhost:8080 - Sistema"
echo "   • http://localhost:8080/clientes - Gestão de clientes"
echo ""

success "Setup concluído! 🎉"