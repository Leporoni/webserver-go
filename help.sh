#!/bin/bash

# Script de ajuda do Sistema de Gestão
# Autor: Sistema de Gestão Go
# Versão: 1.0

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${PURPLE}"
cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                    📚 SISTEMA DE GESTÃO                      ║
║                                                              ║
║                    Guia de Scripts v1.0                     ║
║                                                              ║
║  • PostgreSQL + Go + Nginx                                  ║
║  • CRUD de Clientes Completo                                ║
║  • Interface Web + API REST                                 ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

echo -e "${CYAN}🚀 SCRIPTS DISPONÍVEIS:${NC}"
echo ""

echo -e "${GREEN}0. ./setup_environment.sh${NC}"
echo -e "${BLUE}   • Configura o ambiente automaticamente${NC}"
echo -e "${BLUE}   • Verifica e instala dependências${NC}"
echo -e "${BLUE}   • Configura permissões do Docker${NC}"
echo -e "${BLUE}   • Prepara o projeto para uso${NC}"
echo -e "${YELLOW}   Uso: ./setup_environment.sh${NC}"
echo -e "${YELLOW}   Execute PRIMEIRO se for a primeira vez${NC}"
echo ""

echo -e "${GREEN}1. ./start_and_test.sh${NC}"
echo -e "${BLUE}   • Inicializa todo o sistema automaticamente${NC}"
echo -e "${BLUE}   • Compila a aplicação Go${NC}"
echo -e "${BLUE}   • Inicia PostgreSQL com Docker${NC}"
echo -e "${BLUE}   • Executa testes automatizados do CRUD${NC}"
echo -e "${BLUE}   • Detecta automaticamente se precisa de sudo${NC}"
echo -e "${BLUE}   • Inicia Docker automaticamente se necessário${NC}"
echo -e "${YELLOW}   Uso: ./start_and_test.sh [--clean]${NC}"
echo -e "${YELLOW}   --clean: Remove volumes antigos${NC}"
echo ""

echo -e "${GREEN}2. ./stop_system.sh${NC}"
echo -e "${BLUE}   • Para toda a aplicação${NC}"
echo -e "${BLUE}   • Para containers Docker${NC}"
echo -e "${BLUE}   • Verifica se tudo foi parado corretamente${NC}"
echo -e "${YELLOW}   Uso: ./stop_system.sh${NC}"
echo ""

echo -e "${GREEN}3. ./dev_mode.sh${NC}"
echo -e "${BLUE}   • Modo de desenvolvimento com hot reload${NC}"
echo -e "${BLUE}   • Reinicialização automática ao salvar arquivos${NC}"
echo -e "${BLUE}   • Usa Air para hot reload${NC}"
echo -e "${BLUE}   • PostgreSQL em container${NC}"
echo -e "${YELLOW}   Uso: ./dev_mode.sh${NC}"
echo ""

echo -e "${GREEN}4. ./test_clientes.sh${NC}"
echo -e "${BLUE}   • Testa apenas o CRUD de clientes${NC}"
echo -e "${BLUE}   • Requer sistema já rodando${NC}"
echo -e "${BLUE}   • Testa todos os endpoints da API${NC}"
echo -e "${YELLOW}   Uso: ./test_clientes.sh${NC}"
echo ""

echo -e "${GREEN}5. ./help.sh${NC}"
echo -e "${BLUE}   • Este arquivo de ajuda${NC}"
echo -e "${YELLOW}   Uso: ./help.sh${NC}"
echo ""

echo -e "${CYAN}🌐 ACESSOS APÓS INICIALIZAR:${NC}"
echo ""
echo -e "${GREEN}Interface Web:${NC}"
echo "   • Home: http://localhost:8080"
echo "   • Clientes: http://localhost:8080/clientes"
echo "   • Novo Cliente: http://localhost:8080/clientes/novo"
echo ""
echo -e "${GREEN}API REST:${NC}"
echo "   • Health: http://localhost:8080/health"
echo "   • Listar Clientes: GET http://localhost:8080/api/clientes"
echo "   • Buscar Cliente: GET http://localhost:8080/api/clientes/{id}"
echo "   • Criar Cliente: POST http://localhost:8080/api/clientes"
echo "   • Atualizar Cliente: PUT http://localhost:8080/api/clientes/{id}"
echo "   • Deletar Cliente: DELETE http://localhost:8080/api/clientes/{id}"
echo ""

echo -e "${CYAN}📋 COMANDOS MANUAIS:${NC}"
echo ""
echo -e "${GREEN}Desenvolvimento:${NC}"
echo "   go run main.go                    # Executar aplicação"
echo "   go build -o webserver-go .        # Compilar"
echo "   go test ./...                     # Executar testes"
echo "   go mod tidy                       # Atualizar dependências"
echo ""
echo -e "${GREEN}Docker:${NC}"
echo "   docker-compose up -d postgres     # Só PostgreSQL"
echo "   docker-compose up -d              # Stack completa"
echo "   docker-compose down               # Parar tudo"
echo "   docker-compose logs postgres      # Ver logs"
echo ""

echo -e "${CYAN}🔧 DEPENDÊNCIAS NECESSÁRIAS:${NC}"
echo ""
echo "   • Go 1.18+ (https://golang.org/doc/install)"
echo "   • Docker (https://docs.docker.com/get-docker/)"
echo "   • Docker Compose (https://docs.docker.com/compose/install/)"
echo "   • curl (para testes)"
echo "   • jq (opcional, para testes JSON)"
echo ""

echo -e "${CYAN}🐛 SOLUÇÃO DE PROBLEMAS:${NC}"
echo ""
echo -e "${YELLOW}Erro de permissão Docker:${NC}"
echo "   sudo usermod -aG docker \$USER"
echo "   newgrp docker"
echo "   # ou use os scripts que detectam sudo automaticamente"
echo ""
echo -e "${YELLOW}PostgreSQL não conecta:${NC}"
echo "   docker-compose down"
echo "   docker-compose up -d postgres"
echo "   # aguarde 30 segundos"
echo ""
echo -e "${YELLOW}Aplicação não compila:${NC}"
echo "   go mod tidy"
echo "   go clean -cache"
echo "   go build ."
echo ""
echo -e "${YELLOW}Porta 8080 em uso:${NC}"
echo "   sudo lsof -i :8080"
echo "   kill <PID>"
echo ""

echo -e "${CYAN}📚 DOCUMENTAÇÃO:${NC}"
echo ""
echo "   • README.md - Documentação principal"
echo "   • PROJETO_RESUMO.md - Resumo do projeto"
echo "   • CRUD_CLIENTES_IMPLEMENTADO.md - Detalhes do CRUD"
echo "   • migrations/ - Scripts do banco de dados"
echo ""

echo -e "${GREEN}🎯 INÍCIO RÁPIDO:${NC}"
echo ""
echo -e "${YELLOW}1. Primeira vez (configurar ambiente):${NC}"
echo "   ./setup_environment.sh"
echo "   # Depois faça logout/login ou: newgrp docker"
echo ""
echo -e "${YELLOW}2. Para usar o sistema:${NC}"
echo "   ./start_and_test.sh"
echo ""
echo -e "${YELLOW}3. Para desenvolvimento:${NC}"
echo "   ./dev_mode.sh"
echo ""
echo -e "${YELLOW}4. Para parar tudo:${NC}"
echo "   ./stop_system.sh"
echo ""

echo -e "${PURPLE}✨ Sistema de Gestão v1.0 - Pronto para uso! ✨${NC}"