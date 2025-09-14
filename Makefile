# Makefile para o projeto webserver-go

.PHONY: build run test clean docker-build docker-run install-deps dev

# Variáveis
APP_NAME=webserver-go
DOCKER_IMAGE=webserver-go:latest
PORT=8080

# Compilar a aplicação
build:
	@echo "Compilando aplicação..."
	go build -o bin/$(APP_NAME) .

# Executar a aplicação
run:
	@echo "Executando aplicação na porta $(PORT)..."
	go run main.go

# Executar testes
test:
	@echo "Executando testes..."
	go test -v ./...

# Limpar arquivos compilados
clean:
	@echo "Limpando arquivos..."
	rm -rf bin/
	go clean

# Instalar dependências
install-deps:
	@echo "Instalando dependências..."
	go mod tidy
	go mod download

# Modo desenvolvimento com hot reload (requer air)
dev:
	@echo "Iniciando modo desenvolvimento..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air não encontrado. Instalando..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Construir imagem Docker
docker-build:
	@echo "Construindo imagem Docker..."
	docker build -t $(DOCKER_IMAGE) .

# Executar com Docker
docker-run:
	@echo "Executando com Docker..."
	docker run -p 80:80 -p 8080:8080 $(DOCKER_IMAGE)

# Executar com Docker Compose
docker-compose-up:
	@echo "Iniciando com Docker Compose..."
	docker-compose up --build

# Parar Docker Compose
docker-compose-down:
	@echo "Parando Docker Compose..."
	docker-compose down

# Verificar saúde da aplicação
health-check:
	@echo "Verificando saúde da aplicação..."
	@curl -f http://localhost:$(PORT)/health || echo "Aplicação não está respondendo"

# Instalar Nginx (Ubuntu/Debian)
install-nginx:
	@echo "Instalando Nginx..."
	sudo apt update
	sudo apt install -y nginx

# Configurar Nginx
setup-nginx:
	@echo "Configurando Nginx..."
	sudo cp nginx.conf /etc/nginx/sites-available/webserver-go
	sudo ln -sf /etc/nginx/sites-available/webserver-go /etc/nginx/sites-enabled/
	sudo nginx -t
	sudo systemctl reload nginx

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  build           - Compilar a aplicação"
	@echo "  run             - Executar a aplicação"
	@echo "  test            - Executar testes"
	@echo "  clean           - Limpar arquivos compilados"
	@echo "  install-deps    - Instalar dependências"
	@echo "  dev             - Modo desenvolvimento com hot reload"
	@echo "  docker-build    - Construir imagem Docker"
	@echo "  docker-run      - Executar com Docker"
	@echo "  docker-compose-up   - Iniciar com Docker Compose"
	@echo "  docker-compose-down - Parar Docker Compose"
	@echo "  health-check    - Verificar saúde da aplicação"
	@echo "  install-nginx   - Instalar Nginx"
	@echo "  setup-nginx     - Configurar Nginx"
	@echo "  help            - Mostrar esta ajuda"