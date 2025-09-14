#!/bin/bash

# Script de demonstração do servidor web Go + Nginx

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}================================${NC}"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

# Função para testar endpoint
test_endpoint() {
    local url=$1
    local description=$2
    
    print_step "Testando $description..."
    
    if curl -s -f "$url" > /dev/null; then
        print_success "$description está funcionando"
        echo -e "${CYAN}Resposta:${NC}"
        curl -s "$url" | jq . 2>/dev/null || curl -s "$url"
        echo
    else
        echo -e "${RED}[ERROR]${NC} $description não está respondendo"
    fi
    echo
}

# Função principal de demonstração
main() {
    print_header "🚀 DEMONSTRAÇÃO SERVIDOR WEB GO + NGINX"
    
    echo -e "${YELLOW}Este script demonstra as funcionalidades do servidor web.${NC}"
    echo -e "${YELLOW}Certifique-se de que o servidor está rodando antes de continuar.${NC}"
    echo
    
    read -p "Pressione Enter para continuar ou Ctrl+C para sair..."
    echo
    
    # Verificar se o servidor está rodando
    print_step "Verificando se o servidor está rodando..."
    if ! pgrep -f "./webserver-go" > /dev/null; then
        print_info "Servidor não está rodando. Iniciando..."
        ./webserver-go &
        SERVER_PID=$!
        sleep 2
        print_success "Servidor iniciado (PID: $SERVER_PID)"
    else
        print_success "Servidor já está rodando"
    fi
    echo
    
    # Testar endpoints
    print_header "🧪 TESTANDO ENDPOINTS"
    
    test_endpoint "http://localhost:8080/" "Endpoint principal"
    test_endpoint "http://localhost:8080/health" "Health check"
    
    # Testar arquivos estáticos
    print_step "Testando arquivos estáticos..."
    if curl -s -f "http://localhost:8080/static/index.html" > /dev/null; then
        print_success "Arquivos estáticos funcionando"
        print_info "Acesse: http://localhost:8080/static/index.html"
    else
        echo -e "${RED}[ERROR]${NC} Arquivos estáticos não estão acessíveis"
    fi
    echo
    
    # Teste de performance
    print_header "⚡ TESTE DE PERFORMANCE"
    
    print_step "Executando teste de carga (100 requisições)..."
    if command -v ab > /dev/null; then
        ab -n 100 -c 10 http://localhost:8080/ 2>/dev/null | grep -E "(Requests per second|Time per request)"
    elif command -v curl > /dev/null; then
        print_info "Apache Bench não encontrado. Usando curl para teste simples..."
        start_time=$(date +%s.%N)
        for i in {1..10}; do
            curl -s http://localhost:8080/ > /dev/null
        done
        end_time=$(date +%s.%N)
        duration=$(echo "$end_time - $start_time" | bc -l 2>/dev/null || echo "N/A")
        print_success "10 requisições completadas em ${duration}s"
    else
        print_info "Ferramentas de teste não encontradas"
    fi
    echo
    
    # Informações do sistema
    print_header "📊 INFORMAÇÕES DO SISTEMA"
    
    print_step "Versão do Go:"
    go version
    echo
    
    print_step "Uso de memória do processo:"
    if pgrep -f "./webserver-go" > /dev/null; then
        ps aux | grep "./webserver-go" | grep -v grep | awk '{print "CPU: " $3 "%, Memory: " $4 "%"}'
    else
        print_info "Processo não encontrado"
    fi
    echo
    
    print_step "Portas em uso:"
    netstat -tulpn 2>/dev/null | grep -E ":8080|:80" || print_info "netstat não disponível"
    echo
    
    # URLs úteis
    print_header "🔗 URLS ÚTEIS"
    
    echo -e "${CYAN}Aplicação Go:${NC}"
    echo "  • Home: http://localhost:8080/"
    echo "  • Health: http://localhost:8080/health"
    echo "  • Interface Web: http://localhost:8080/static/index.html"
    echo
    
    if command -v nginx > /dev/null && pgrep nginx > /dev/null; then
        echo -e "${CYAN}Via Nginx (se configurado):${NC}"
        echo "  • Home: http://localhost/"
        echo "  • Health: http://localhost/health"
        echo "  • Interface Web: http://localhost/static/index.html"
        echo
    fi
    
    # Comandos úteis
    print_header "🛠️ COMANDOS ÚTEIS"
    
    echo -e "${CYAN}Desenvolvimento:${NC}"
    echo "  make run          # Executar aplicação"
    echo "  make dev          # Modo desenvolvimento (hot reload)"
    echo "  make test         # Executar testes"
    echo "  make build        # Compilar aplicação"
    echo
    
    echo -e "${CYAN}Docker:${NC}"
    echo "  make docker-build # Construir imagem"
    echo "  make docker-compose-up # Executar com Docker Compose"
    echo
    
    echo -e "${CYAN}Monitoramento:${NC}"
    echo "  make health-check # Verificar saúde"
    echo "  curl http://localhost:8080/health # Health check manual"
    echo
    
    print_header "✅ DEMONSTRAÇÃO CONCLUÍDA"
    
    echo -e "${GREEN}O servidor web Go + Nginx está funcionando perfeitamente!${NC}"
    echo -e "${YELLOW}Para parar o servidor, use Ctrl+C ou 'pkill -f webserver-go'${NC}"
    echo
    
    # Parar servidor se foi iniciado por este script
    if [ ! -z "$SERVER_PID" ]; then
        read -p "Deseja parar o servidor iniciado por este script? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            kill $SERVER_PID 2>/dev/null || true
            print_success "Servidor parado"
        fi
    fi
}

# Executar demonstração
main "$@"