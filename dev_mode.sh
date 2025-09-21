#!/bin/bash

# Script para modo de desenvolvimento com hot reload
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
║                    🔥 MODO DESENVOLVIMENTO                   ║
║                                                              ║
║              Hot Reload + PostgreSQL                        ║
║                                                              ║
║  • Reinicialização automática ao salvar                     ║
║  • Logs em tempo real                                       ║
║  • PostgreSQL em container                                  ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Verificar se Air está instalado
if ! command -v air &> /dev/null; then
    warning "Air (hot reload) não está instalado"
    info "Instalando Air..."
    if go install github.com/cosmtrek/air@latest; then
        success "Air instalado com sucesso"
    else
        error "Falha ao instalar Air"
        info "Instale manualmente: go install github.com/cosmtrek/air@latest"
        exit 1
    fi
fi

# Verificar se estamos no diretório correto
if [ ! -f "main.go" ] || [ ! -f "docker-compose.yml" ]; then
    error "Execute este script no diretório raiz do projeto"
    exit 1
fi

# Detectar se precisa de sudo para Docker e tentar iniciar se necessário
DOCKER_COMPOSE_CMD="docker-compose"
if ! docker ps &>/dev/null; then
    if sudo docker ps &>/dev/null; then
        warning "Docker requer sudo - usando sudo para comandos Docker"
        DOCKER_COMPOSE_CMD="sudo docker-compose"
    else
        warning "Docker não está funcionando - tentando iniciar..."
        
        if sudo systemctl start docker; then
            success "Docker iniciado com sucesso"
            sleep 3
            
            if ! docker ps &>/dev/null; then
                if sudo docker ps &>/dev/null; then
                    warning "Docker requer sudo - usando sudo para comandos Docker"
                    DOCKER_COMPOSE_CMD="sudo docker-compose"
                else
                    error "Docker ainda não está funcionando"
                    exit 1
                fi
            fi
        else
            error "Falha ao iniciar Docker"
            info "Tente: sudo systemctl start docker"
            exit 1
        fi
    fi
fi

# Parar containers existentes
info "Parando containers existentes..."
$DOCKER_COMPOSE_CMD down --remove-orphans 2>/dev/null || true

# Iniciar PostgreSQL
info "Iniciando PostgreSQL..."
if ! $DOCKER_COMPOSE_CMD up -d postgres; then
    error "Falha ao iniciar PostgreSQL"
    exit 1
fi
success "PostgreSQL iniciado"

# Aguardar PostgreSQL
info "Aguardando PostgreSQL ficar pronto..."
for i in {1..30}; do
    if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U gestao_user -d gestao_db &>/dev/null; then
        break
    fi
    echo -n "."
    sleep 2
done
echo ""
success "PostgreSQL está pronto"

# Configurar variáveis de ambiente
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=gestao_db
export DB_USER=gestao_user
export DB_PASSWORD=gestao_pass

# Verificar se .air.toml existe
if [ ! -f ".air.toml" ]; then
    warning "Arquivo .air.toml não encontrado, criando configuração padrão..."
    cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "static", "templates"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
EOF
    success "Arquivo .air.toml criado"
fi

# Criar diretório tmp se não existir
mkdir -p tmp

echo ""
success "🔥 Iniciando modo de desenvolvimento..."
echo ""
info "📝 Comandos úteis durante desenvolvimento:"
echo "   • Ctrl+C: Parar o servidor"
echo "   • Salvar arquivo .go: Reinicialização automática"
echo "   • http://localhost:8080: Acessar aplicação"
echo "   • http://localhost:8080/clientes: Gestão de clientes"
echo ""
warning "Logs da aplicação aparecerão abaixo:"
echo ""

# Função para cleanup ao sair
cleanup() {
    echo ""
    info "Parando modo de desenvolvimento..."
    $DOCKER_COMPOSE_CMD down --remove-orphans 2>/dev/null || true
    success "Modo de desenvolvimento parado"
    exit 0
}

# Configurar trap para cleanup
trap cleanup SIGINT SIGTERM

# Iniciar Air (hot reload)
air

# Se chegou aqui, Air foi interrompido
cleanup