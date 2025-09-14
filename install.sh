#!/bin/bash

# Script de instalação para o projeto webserver-go

set -e

echo "🚀 Instalando Servidor Web Go + Nginx"
echo "======================================"

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar se Go está instalado
check_go() {
    print_status "Verificando instalação do Go..."
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_success "Go encontrado: $GO_VERSION"
    else
        print_error "Go não encontrado. Por favor, instale Go 1.19+ primeiro."
        echo "Visite: https://golang.org/dl/"
        exit 1
    fi
}

# Instalar dependências do sistema
install_system_deps() {
    print_status "Instalando dependências do sistema..."
    
    if command -v apt &> /dev/null; then
        # Ubuntu/Debian
        sudo apt update
        sudo apt install -y curl wget make
        print_success "Dependências instaladas via apt"
    elif command -v yum &> /dev/null; then
        # CentOS/RHEL
        sudo yum install -y curl wget make
        print_success "Dependências instaladas via yum"
    elif command -v brew &> /dev/null; then
        # macOS
        brew install curl wget make
        print_success "Dependências instaladas via brew"
    else
        print_warning "Gerenciador de pacotes não reconhecido. Instale manualmente: curl, wget, make"
    fi
}

# Instalar Nginx
install_nginx() {
    print_status "Verificando instalação do Nginx..."
    
    if command -v nginx &> /dev/null; then
        NGINX_VERSION=$(nginx -v 2>&1 | awk '{print $3}')
        print_success "Nginx já instalado: $NGINX_VERSION"
        return 0
    fi
    
    print_status "Instalando Nginx..."
    
    if command -v apt &> /dev/null; then
        sudo apt install -y nginx
    elif command -v yum &> /dev/null; then
        sudo yum install -y nginx
    elif command -v brew &> /dev/null; then
        brew install nginx
    else
        print_error "Não foi possível instalar Nginx automaticamente."
        print_warning "Por favor, instale Nginx manualmente."
        return 1
    fi
    
    print_success "Nginx instalado com sucesso"
}

# Instalar Docker (opcional)
install_docker() {
    print_status "Verificando instalação do Docker..."
    
    if command -v docker &> /dev/null; then
        print_success "Docker já instalado"
        return 0
    fi
    
    read -p "Deseja instalar Docker? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Pulando instalação do Docker"
        return 0
    fi
    
    print_status "Instalando Docker..."
    
    if command -v apt &> /dev/null; then
        # Ubuntu/Debian
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        sudo usermod -aG docker $USER
        rm get-docker.sh
    else
        print_warning "Instalação automática do Docker disponível apenas para Ubuntu/Debian"
        print_warning "Visite: https://docs.docker.com/get-docker/"
        return 1
    fi
    
    print_success "Docker instalado. Faça logout/login para usar sem sudo"
}

# Instalar dependências Go
install_go_deps() {
    print_status "Instalando dependências Go..."
    go mod tidy
    go mod download
    print_success "Dependências Go instaladas"
}

# Instalar Air para hot reload
install_air() {
    print_status "Instalando Air para hot reload..."
    go install github.com/cosmtrek/air@latest
    print_success "Air instalado"
}

# Compilar aplicação
build_app() {
    print_status "Compilando aplicação..."
    go build -o webserver-go .
    print_success "Aplicação compilada com sucesso"
}

# Configurar Nginx
setup_nginx() {
    if ! command -v nginx &> /dev/null; then
        print_warning "Nginx não instalado. Pulando configuração."
        return 0
    fi
    
    read -p "Deseja configurar Nginx? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Pulando configuração do Nginx"
        return 0
    fi
    
    print_status "Configurando Nginx..."
    
    # Backup da configuração existente
    if [ -f /etc/nginx/sites-enabled/default ]; then
        sudo mv /etc/nginx/sites-enabled/default /etc/nginx/sites-enabled/default.backup
        print_status "Backup da configuração padrão criado"
    fi
    
    # Copiar nossa configuração
    sudo cp nginx.conf /etc/nginx/sites-available/webserver-go
    sudo ln -sf /etc/nginx/sites-available/webserver-go /etc/nginx/sites-enabled/
    
    # Testar configuração
    if sudo nginx -t; then
        sudo systemctl reload nginx
        print_success "Nginx configurado e recarregado"
    else
        print_error "Erro na configuração do Nginx"
        return 1
    fi
}

# Criar diretórios necessários
create_dirs() {
    print_status "Criando diretórios necessários..."
    mkdir -p logs tmp bin
    print_success "Diretórios criados"
}

# Função principal
main() {
    echo
    print_status "Iniciando instalação..."
    echo
    
    check_go
    install_system_deps
    install_nginx
    install_docker
    install_go_deps
    install_air
    create_dirs
    build_app
    setup_nginx
    
    echo
    print_success "🎉 Instalação concluída com sucesso!"
    echo
    echo "Para executar a aplicação:"
    echo "  make run          # Executar diretamente"
    echo "  make dev          # Modo desenvolvimento"
    echo "  make docker-compose-up  # Com Docker"
    echo
    echo "Acesse:"
    echo "  http://localhost:8080     # Aplicação Go"
    echo "  http://localhost          # Via Nginx (se configurado)"
    echo
}

# Executar função principal
main "$@"