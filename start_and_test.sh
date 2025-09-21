#!/bin/bash

# Script completo para iniciar e testar o Sistema de Gestão
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
║                    🏢 SISTEMA DE GESTÃO                      ║
║                                                              ║
║              🚀 Inicializador Automático v1.0               ║
║                                                              ║
║  • PostgreSQL + Go + Nginx                                  ║
║  • CRUD de Clientes Completo                                ║
║  • Interface Web + API REST                                 ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Verificar se Docker está disponível e detectar se precisa de sudo
log "Verificando dependências..."

if ! command -v docker &> /dev/null; then
    error "Docker não está instalado ou não está no PATH"
    info "Instale o Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    error "Docker Compose não está instalado ou não está no PATH"
    info "Instale o Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

if ! command -v go &> /dev/null; then
    error "Go não está instalado ou não está no PATH"
    info "Instale o Go: https://golang.org/doc/install"
    exit 1
fi

# Detectar se precisa de sudo para Docker e tentar iniciar se necessário
DOCKER_CMD="docker"
DOCKER_COMPOSE_CMD="docker-compose"

if ! docker ps &>/dev/null; then
    if sudo docker ps &>/dev/null; then
        warning "Docker requer sudo - usando sudo para comandos Docker"
        DOCKER_CMD="sudo docker"
        DOCKER_COMPOSE_CMD="sudo docker-compose"
    else
        warning "Docker não está funcionando - tentando iniciar..."
        
        # Tentar iniciar o Docker
        if sudo systemctl start docker; then
            success "Docker iniciado com sucesso"
            sleep 3  # Aguardar Docker inicializar
            
            # Verificar novamente
            if ! docker ps &>/dev/null; then
                if sudo docker ps &>/dev/null; then
                    warning "Docker requer sudo - usando sudo para comandos Docker"
                    DOCKER_CMD="sudo docker"
                    DOCKER_COMPOSE_CMD="sudo docker-compose"
                else
                    error "Docker ainda não está funcionando após inicialização"
                    info "Tente manualmente: sudo systemctl start docker"
                    info "Ou verifique se o Docker está instalado corretamente"
                    exit 1
                fi
            else
                success "Docker funcionando sem sudo"
            fi
        else
            error "Falha ao iniciar Docker"
            info "Comandos para tentar manualmente:"
            info "  sudo systemctl start docker"
            info "  sudo systemctl enable docker"
            info "  sudo usermod -aG docker \$USER && newgrp docker"
            exit 1
        fi
    fi
else
    success "Docker funcionando sem sudo"
fi

success "Todas as dependências estão disponíveis"

# Verificar se estamos no diretório correto
if [ ! -f "main.go" ] || [ ! -f "docker-compose.yml" ]; then
    error "Execute este script no diretório raiz do projeto (onde estão main.go e docker-compose.yml)"
    exit 1
fi

# Parar containers existentes (se houver)
log "Parando containers existentes..."
$DOCKER_COMPOSE_CMD down --remove-orphans 2>/dev/null || true
success "Containers parados"

# Limpar volumes antigos se necessário
if [ "$1" = "--clean" ]; then
    warning "Modo de limpeza ativado - removendo volumes..."
    $DOCKER_COMPOSE_CMD down -v 2>/dev/null || true
    $DOCKER_CMD volume prune -f 2>/dev/null || true
    success "Volumes limpos"
fi

# Compilar aplicação Go
log "Compilando aplicação Go..."
go mod tidy
if ! go build -o webserver-go .; then
    error "Falha na compilação da aplicação Go"
    exit 1
fi
success "Aplicação compilada com sucesso"

# Iniciar PostgreSQL
log "Iniciando PostgreSQL..."
if ! $DOCKER_COMPOSE_CMD up -d postgres; then
    error "Falha ao iniciar PostgreSQL"
    exit 1
fi
success "PostgreSQL iniciado"

# Aguardar PostgreSQL ficar pronto
log "Aguardando PostgreSQL ficar pronto..."
POSTGRES_READY=false
for i in {1..30}; do
    if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U gestao_user -d gestao_db &>/dev/null; then
        POSTGRES_READY=true
        break
    fi
    echo -n "."
    sleep 2
done
echo ""

if [ "$POSTGRES_READY" = false ]; then
    error "PostgreSQL não ficou pronto em 60 segundos"
    log "Verificando logs do PostgreSQL..."
    $DOCKER_COMPOSE_CMD logs postgres
    exit 1
fi
success "PostgreSQL está pronto e aceitando conexões"

# Verificar se as tabelas foram criadas
log "Verificando estrutura do banco de dados..."
TABLES_COUNT=$($DOCKER_COMPOSE_CMD exec -T postgres psql -U gestao_user -d gestao_db -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | tr -d ' \n' || echo "0")

if [ "$TABLES_COUNT" -eq "0" ]; then
    warning "Tabelas não encontradas, executando migrations manualmente..."
    if [ -f "migrations/001_initial_schema.up.sql" ]; then
        $DOCKER_COMPOSE_CMD exec -T postgres psql -U gestao_user -d gestao_db < migrations/001_initial_schema.up.sql
        success "Migrations executadas manualmente"
    else
        error "Arquivo de migration não encontrado"
        exit 1
    fi
else
    success "Banco de dados já possui $TABLES_COUNT tabela(s)"
fi

# Iniciar aplicação Go em background
log "Iniciando aplicação Go..."
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=gestao_db
export DB_USER=gestao_user
export DB_PASSWORD=gestao_pass

# Matar processo anterior se existir
pkill -f "webserver-go" 2>/dev/null || true
sleep 2

# Iniciar aplicação
./webserver-go &
APP_PID=$!
sleep 3

# Verificar se a aplicação está rodando
if ! kill -0 $APP_PID 2>/dev/null; then
    error "Aplicação Go falhou ao iniciar"
    exit 1
fi

# Aguardar aplicação ficar pronta
log "Aguardando aplicação ficar pronta..."
APP_READY=false
for i in {1..15}; do
    if curl -s http://localhost:8080/health &>/dev/null; then
        APP_READY=true
        break
    fi
    echo -n "."
    sleep 2
done
echo ""

if [ "$APP_READY" = false ]; then
    error "Aplicação não ficou pronta em 30 segundos"
    log "Verificando se a aplicação ainda está rodando..."
    if kill -0 $APP_PID 2>/dev/null; then
        warning "Aplicação está rodando mas não responde no health check"
    else
        error "Aplicação parou de funcionar"
    fi
    exit 1
fi

success "Aplicação está pronta e respondendo"

# Executar testes automatizados
echo ""
log "🧪 Executando testes automatizados do CRUD de Clientes..."
echo ""

# Função para testar endpoint
test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local expected_status=$4
    local description=$5
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "%{http_code}" -X $method "$url" -H "Content-Type: application/json" -d "$data")
    else
        response=$(curl -s -w "%{http_code}" -X $method "$url")
    fi
    
    status_code="${response: -3}"
    body="${response%???}"
    
    if [ "$status_code" = "$expected_status" ]; then
        success "$description (Status: $status_code)"
        echo "$body"
    else
        error "$description (Esperado: $expected_status, Recebido: $status_code)"
        echo "$body"
        return 1
    fi
}

# Teste 1: Health Check
info "Teste 1: Health Check"
test_endpoint "GET" "http://localhost:8080/health" "" "200" "Health check"

echo ""

# Teste 2: Criar cliente
info "Teste 2: Criar Cliente"
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/clientes \
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
  }')

if echo "$CREATE_RESPONSE" | grep -q '"success":true'; then
    CLIENTE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    success "Cliente criado com ID: $CLIENTE_ID"
else
    error "Falha ao criar cliente"
    echo "$CREATE_RESPONSE"
    exit 1
fi

echo ""

# Teste 3: Buscar cliente por ID
info "Teste 3: Buscar Cliente por ID"
GET_RESPONSE=$(curl -s http://localhost:8080/api/clientes/$CLIENTE_ID)
if echo "$GET_RESPONSE" | grep -q '"nome":"João Silva Teste"'; then
    success "Cliente encontrado corretamente"
else
    error "Falha ao buscar cliente"
    echo "$GET_RESPONSE"
fi

echo ""

# Teste 4: Listar clientes
info "Teste 4: Listar Clientes"
LIST_RESPONSE=$(curl -s http://localhost:8080/api/clientes)
if echo "$LIST_RESPONSE" | grep -q '"success":true'; then
    TOTAL=$(echo "$LIST_RESPONSE" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    success "Listagem funcionando - Total: $TOTAL cliente(s)"
else
    error "Falha ao listar clientes"
    echo "$LIST_RESPONSE"
fi

echo ""

# Teste 5: Atualizar cliente
info "Teste 5: Atualizar Cliente"
UPDATE_RESPONSE=$(curl -s -X PUT http://localhost:8080/api/clientes/$CLIENTE_ID \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva Atualizado",
    "email": "joao.atualizado@email.com",
    "tipo_pessoa": "fisica"
  }')

if echo "$UPDATE_RESPONSE" | grep -q '"success":true'; then
    success "Cliente atualizado com sucesso"
else
    error "Falha ao atualizar cliente"
    echo "$UPDATE_RESPONSE"
fi

echo ""

# Teste 6: Deletar cliente
info "Teste 6: Deletar Cliente"
DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/api/clientes/$CLIENTE_ID)
if echo "$DELETE_RESPONSE" | grep -q '"success":true'; then
    success "Cliente deletado com sucesso"
else
    error "Falha ao deletar cliente"
    echo "$DELETE_RESPONSE"
fi

echo ""

# Teste 7: Verificar se foi deletado
info "Teste 7: Verificar Deleção"
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/clientes/$CLIENTE_ID)
if [ "$STATUS_CODE" = "404" ]; then
    success "Cliente não encontrado (deletado corretamente)"
else
    error "Cliente ainda existe (erro na deleção)"
fi

# Resumo final
echo ""
echo -e "${PURPLE}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${PURPLE}║                        🎉 SUCESSO!                           ║${NC}"
echo -e "${PURPLE}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
success "Sistema de Gestão está funcionando perfeitamente!"
echo ""
info "🌐 Acessos disponíveis:"
echo "   • Home: http://localhost:8080"
echo "   • Clientes: http://localhost:8080/clientes"
echo "   • Novo Cliente: http://localhost:8080/clientes/novo"
echo "   • API Health: http://localhost:8080/health"
echo "   • API Clientes: http://localhost:8080/api/clientes"
echo ""
info "📊 Serviços rodando:"
echo "   • PostgreSQL: localhost:5432"
echo "   • Aplicação Go: localhost:8080"
echo "   • PID da aplicação: $APP_PID"
echo ""
warning "Para parar o sistema:"
echo "   • Aplicação: kill $APP_PID"
echo "   • PostgreSQL: $DOCKER_COMPOSE_CMD down"
echo "   • Tudo: $DOCKER_COMPOSE_CMD down && kill $APP_PID"
echo "   • Script automático: ./stop_system.sh"
echo ""
info "📝 Logs da aplicação:"
echo "   • tail -f /tmp/webserver-go.log (se configurado)"
echo "   • ou monitore o processo PID $APP_PID"
echo ""
success "Sistema pronto para uso! 🚀"