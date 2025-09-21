# 🚀 Servidor Web Go + Nginx

Um servidor web moderno construído com Go (Golang) e Nginx como proxy reverso.

## 📋 Características

- **Backend em Go**: Servidor HTTP rápido e eficiente
- **Proxy Reverso Nginx**: Load balancing e serving de arquivos estáticos
- **API RESTful**: Endpoints JSON para integração
- **Health Check**: Monitoramento de saúde da aplicação
- **Docker Support**: Containerização completa
- **Hot Reload**: Desenvolvimento com recarga automática

## 🛠️ Pré-requisitos

- Go 1.19+ 
- Nginx (opcional, pode usar Docker)
- Docker & Docker Compose (opcional)
- Make (opcional, para usar o Makefile)

## 🚀 Instalação e Execução

### Método 1: Execução Direta

1. **Instalar dependências:**
   ```bash
   make install-deps
   # ou
   go mod tidy
   ```

2. **Executar a aplicação:**
   ```bash
   make run
   # ou
   go run main.go
   ```

3. **Acessar a aplicação:**
   - Aplicação Go: http://localhost:8080
   - Interface Web: http://localhost:8080/static/index.html
   - Health Check: http://localhost:8080/health

### Método 2: Com Nginx (Recomendado)

1. **Instalar Nginx:**
   ```bash
   make install-nginx
   # ou
   sudo apt update && sudo apt install -y nginx
   ```

2. **Configurar Nginx:**
   ```bash
   make setup-nginx
   ```

3. **Executar a aplicação:**
   ```bash
   make run
   ```

4. **Acessar via Nginx:**
   - http://localhost (porta 80)

### Método 3: Docker (Mais Fácil)

1. **Usando Docker Compose:**
   ```bash
   make docker-compose-up
   # ou
   docker-compose up --build
   ```

2. **Acessar a aplicação:**
   - http://localhost (Nginx na porta 80)
   - http://localhost:8080 (Go direto)

## 📁 Estrutura do Projeto

```
webserver-go/
├── main.go              # Aplicação principal Go
├── static/              # Arquivos estáticos
│   └── index.html       # Interface web
├── nginx.conf           # Configuração do Nginx
├── Dockerfile           # Imagem Docker
├── docker-compose.yml   # Orquestração Docker
├── Makefile            # Comandos automatizados
├── start.sh            # Script de inicialização
├── go.mod              # Dependências Go
└── README.md           # Este arquivo
```

## 🔧 Comandos Disponíveis

```bash
# Desenvolvimento
make run                 # Executar aplicação
make dev                 # Modo desenvolvimento (hot reload)
make build              # Compilar aplicação
make test               # Executar testes
make clean              # Limpar arquivos

# Docker
make docker-build       # Construir imagem
make docker-run         # Executar container
make docker-compose-up  # Iniciar com compose
make docker-compose-down # Parar compose

# Nginx
make install-nginx      # Instalar Nginx
make setup-nginx        # Configurar Nginx

# Utilitários
make health-check       # Verificar saúde
make help              # Mostrar ajuda
```

## 🌐 Endpoints da API

| Endpoint | Método | Descrição |
|----------|--------|-----------|
| `/` | GET | Página inicial (JSON) |
| `/health` | GET | Health check |
| `/static/*` | GET | Arquivos estáticos |

### Exemplos de Resposta

**GET /**
```json
{
  "message": "Bem-vindo ao servidor web Go!",
  "timestamp": "2024-01-15T10:30:00Z",
  "status": "success"
}
```

**GET /health**
```json
{
  "message": "Servidor funcionando perfeitamente",
  "timestamp": "2024-01-15T10:30:00Z",
  "status": "healthy"
}
```

## 🐳 Docker

### Construir e Executar

```bash
# Construir imagem
docker build -t webserver-go .

# Executar container
docker run -p 80:80 -p 8080:8080 webserver-go

# Ou usar Docker Compose
docker-compose up --build
```

### Variáveis de Ambiente

- `ENV`: Ambiente de execução (development/production)

## 🔧 Configuração do Nginx

O arquivo `nginx.conf` está configurado para:

- **Proxy Reverso**: Redireciona requisições para o servidor Go
- **Arquivos Estáticos**: Servidos diretamente pelo Nginx
- **Cache**: Headers de cache para arquivos estáticos
- **Load Balancing**: Pronto para múltiplas instâncias

### Configuração Manual do Nginx

```bash
# Copiar configuração
sudo cp nginx.conf /etc/nginx/sites-available/webserver-go

# Ativar site
sudo ln -s /etc/nginx/sites-available/webserver-go /etc/nginx/sites-enabled/

# Testar configuração
sudo nginx -t

# Recarregar Nginx
sudo systemctl reload nginx
```

## 🚀 Desenvolvimento

### Hot Reload

Para desenvolvimento com recarga automática:

```bash
# Instalar Air (se não tiver)
go install github.com/cosmtrek/air@latest

# Executar em modo desenvolvimento
make dev
```

### Adicionando Novas Rotas

1. Adicione o handler em `main.go`:
```go
func newHandler(w http.ResponseWriter, r *http.Request) {
    // Sua lógica aqui
}
```

2. Registre a rota:
```go
http.HandleFunc("/nova-rota", loggingMiddleware(newHandler))
```

## 📊 Monitoramento

### Health Check

```bash
# Verificar saúde da aplicação
curl http://localhost:8080/health

# Ou usar o Makefile
make health-check
```

### Logs

- **Aplicação Go**: Logs no console
- **Nginx**: `/var/log/nginx/access.log` e `/var/log/nginx/error.log`
- **Docker**: `docker-compose logs`

## 🔒 Segurança

- Servidor roda como usuário não-root no Docker
- Headers de segurança configurados no Nginx
- Timeouts configurados para evitar ataques DoS

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.

## 🆘 Suporte

Se encontrar problemas:

1. Verifique os logs: `docker-compose logs`
2. Teste a conectividade: `make health-check`
3. Verifique se as portas estão livres: `netstat -tulpn | grep :8080`

---

**Desenvolvido com ❤️ usando Go e Nginx**