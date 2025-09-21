# 🏢 Sistema de Gestão - Go + PostgreSQL + Nginx

Um sistema completo de gestão construído com Go, PostgreSQL e Nginx, com CRUD de clientes implementado.

## 🎯 Características

- **🔥 CRUD de Clientes Completo**: Service layer + API REST + Interface Web
- **🗄️ PostgreSQL**: Banco de dados robusto com migrations automáticas
- **🌐 Interface Web Moderna**: Templates responsivos com busca e paginação
- **🚀 API RESTful**: Endpoints JSON completos para integração
- **🐳 Docker Ready**: Containerização completa com docker-compose
- **🔄 Hot Reload**: Desenvolvimento com recarga automática (Air)
- **✅ Testes Automatizados**: Scripts de teste do CRUD completo
- **🛡️ Validações**: Validações de entrada e regras de negócio

## 🛠️ Pré-requisitos

- **Go 1.18+** (https://golang.org/doc/install)
- **Docker & Docker Compose** (https://docs.docker.com/get-docker/)
- **curl** (para testes)
- **jq** (opcional, para testes JSON)

## 🚀 Como Executar

### 🎆 Método 1: Inicialização Automática (RECOMENDADO)
```bash
cd webserver-go
./start_and_test.sh
# Inicia tudo automaticamente + testes
# Acesse: http://localhost:8080
```

### 🔥 Método 2: Modo Desenvolvimento (Hot Reload)
```bash
cd webserver-go
./dev_mode.sh
# Hot reload automático ao salvar arquivos
# Acesse: http://localhost:8080
```

### 🐳 Método 3: Docker Completo
```bash
cd webserver-go
sudo docker-compose up -d
# Acesse: http://localhost (Nginx) ou http://localhost:8080 (Go)
```

### 🛠️ Método 4: Manual
```bash
cd webserver-go
sudo docker-compose up -d postgres  # PostgreSQL
go run main.go                       # Aplicação
# Acesse: http://localhost:8080
```

## 🔧 Scripts Disponíveis

```bash
# 🚀 Inicialização e Testes
./start_and_test.sh        # Inicia tudo + testa CRUD
./start_and_test.sh --clean # Com limpeza de volumes

# 🔥 Desenvolvimento
./dev_mode.sh              # Hot reload com Air

# 🛑 Parar Sistema
./stop_system.sh           # Para tudo (app + containers)

# 🧪 Testes
./test_clientes.sh         # Testa apenas CRUD clientes

# 📚 Ajuda
./help.sh                  # Guia completo de uso

# 🐳 Docker Manual
sudo docker-compose up -d postgres  # Só PostgreSQL
sudo docker-compose up -d           # Stack completa
sudo docker-compose down            # Parar tudo

# 🛠️ Go Manual
go run main.go             # Executar aplicação
go build -o webserver-go . # Compilar
go test ./...              # Testes unitários
```

## 📁 Estrutura do Projeto

```
webserver-go/
├── 🚀 APLICAÇÃO PRINCIPAL
│   ├── main.go                    # Servidor HTTP principal
│   ├── go.mod                     # Dependências Go
│   └── go.sum                     # Lock de dependências
│
├── 🏗️ ARQUITETURA INTERNA
│   └── internal/
│       ├── handlers/              # HTTP handlers
│       │   ├── api.go            # Handlers básicos
│       │   └── clientes.go       # CRUD de clientes
│       ├── services/              # Business logic
│       │   └── cliente_service.go # Service de clientes
│       ├── models/                # Estruturas de dados
│       │   ├── cliente.go        # Modelo Cliente
│       │   ├── produto.go        # Modelo Produto
│       │   └── estoque.go        # Modelo Estoque
│       ├── database/              # Conexão e migrations
│       │   ├── connection.go     # Conexão PostgreSQL
│       │   └── migrate.go        # Sistema de migrations
│       └── middleware/            # Middlewares HTTP
│           └── logging.go        # Middleware de logging
│
├── 🗄️ BANCO DE DADOS
│   └── migrations/
│       ├── 001_initial_schema.up.sql   # Schema inicial
│       └── 001_initial_schema.down.sql # Rollback
│
├── 🐳 CONTAINERIZAÇÃO
│   ├── Dockerfile                # Imagem Docker otimizada
│   ├── docker-compose.yml        # Orquestração completa
│   └── nginx.conf                # Configuração Nginx
│
├── 🔧 SCRIPTS DE AUTOMAÇÃO
│   ├── start_and_test.sh         # Inicialização completa
│   ├── stop_system.sh            # Parar sistema
│   ├── dev_mode.sh               # Modo desenvolvimento
│   ├── test_clientes.sh          # Testes CRUD
│   └── help.sh                   # Guia de ajuda
│
├── 📚 DOCUMENTAÇÃO
│   ├── README.md                 # Este arquivo
│   ├── PROJETO_RESUMO.md         # Resumo executivo
│   └── CRUD_CLIENTES_IMPLEMENTADO.md # Detalhes do CRUD
│
└── 🔧 CONFIGURAÇÕES
    ├── .air.toml                 # Configuração hot reload
    ├── .env.example              # Variáveis de ambiente
    └── .gitignore                # Controle de versão
```

## 🌐 Endpoints Disponíveis

### 🏠 Interface Web
- **Home**: http://localhost:8080
- **Clientes**: http://localhost:8080/clientes
- **Novo Cliente**: http://localhost:8080/clientes/novo
- **Ver Cliente**: http://localhost:8080/clientes/{id}
- **Editar Cliente**: http://localhost:8080/clientes/{id}/editar

### 🔌 API REST

#### Sistema
| Endpoint | Método | Descrição |
|----------|--------|-----------|
| `/` | GET | Home da API |
| `/health` | GET | Health check |

#### Clientes
| Endpoint | Método | Descrição |
|----------|--------|-----------|
| `/api/clientes` | GET | Listar clientes (com paginação) |
| `/api/clientes/{id}` | GET | Buscar cliente por ID |
| `/api/clientes` | POST | Criar novo cliente |
| `/api/clientes/{id}` | PUT | Atualizar cliente |
| `/api/clientes/{id}` | DELETE | Deletar cliente |

### 📋 Exemplos de Uso da API

#### Criar Cliente
```bash
curl -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva",
    "email": "joao@email.com",
    "telefone": "(11) 99999-9999",
    "tipo_pessoa": "fisica",
    "endereco": "Rua das Flores, 123",
    "cidade": "São Paulo",
    "estado": "SP",
    "cep": "01234567"
  }'
```

#### Listar Clientes
```bash
curl http://localhost:8080/api/clientes
curl "http://localhost:8080/api/clientes?page=1&limit=10&search=João"
```

#### Buscar Cliente
```bash
curl http://localhost:8080/api/clientes/{id}
```

#### Atualizar Cliente
```bash
curl -X PUT http://localhost:8080/api/clientes/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "João Silva Santos",
    "email": "joao.santos@email.com",
    "tipo_pessoa": "fisica"
  }'
```

#### Deletar Cliente
```bash
curl -X DELETE http://localhost:8080/api/clientes/{id}
```

## 🗄️ Banco de Dados

### PostgreSQL Schema
```sql
-- Clientes
CREATE TABLE clientes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    telefone VARCHAR(20),
    endereco TEXT,
    cidade VARCHAR(100),
    estado VARCHAR(2),
    cep VARCHAR(10),
    cpf_cnpj VARCHAR(20) UNIQUE,
    tipo_pessoa VARCHAR(10) CHECK (tipo_pessoa IN ('fisica', 'juridica')),
    ativo BOOLEAN DEFAULT true,
    observacoes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Funcionalidades do CRUD
- ✅ **Validações**: Nome obrigatório, email válido, tipo de pessoa
- ✅ **Paginação**: 20 itens por página (configurável)
- ✅ **Busca**: Por nome ou email
- ✅ **Ordenação**: Alfabética por nome
- ✅ **Status**: Ativo/Inativo
- ✅ **Timestamps**: Criação e atualização automáticas
- ✅ **UUIDs**: Identificadores únicos
- ✅ **Constraints**: Email e CPF/CNPJ únicos

## 🔧 Desenvolvimento

### Hot Reload
```bash
./dev_mode.sh
# Instala Air automaticamente se necessário
# Reinicialização automática ao salvar arquivos .go
```

### Adicionando Novas Funcionalidades
1. **Modelo**: Criar em `internal/models/`
2. **Service**: Implementar em `internal/services/`
3. **Handler**: Adicionar em `internal/handlers/`
4. **Rotas**: Registrar em `main.go`
5. **Migration**: Criar em `migrations/`

### Estrutura de Resposta da API
```json
{
  "success": true,
  "message": "Operação realizada com sucesso",
  "data": {
    // dados da resposta
  },
  "error": null
}
```

### Resposta Paginada
```json
{
  "success": true,
  "message": "Clientes encontrados com sucesso",
  "data": {
    "data": [...],
    "total": 100,
    "page": 1,
    "limit": 20,
    "total_pages": 5
  }
}
```

## 🧪 Testes

### Testes Automatizados
```bash
./test_clientes.sh
# Testa todos os endpoints do CRUD
# Cria, busca, atualiza e deleta cliente de teste
```

### Testes Manuais
```bash
# Health check
curl http://localhost:8080/health

# Interface web
open http://localhost:8080/clientes
```

## 📊 Status do Projeto

**🎉 CRUD DE CLIENTES IMPLEMENTADO COM SUCESSO!**

- ✅ **Service Layer** completo para clientes
- ✅ **API REST** com todos endpoints CRUD
- ✅ **Interface Web** responsiva e moderna
- ✅ **PostgreSQL** integrado e funcionando
- ✅ **Validações** de entrada e negócio
- ✅ **Paginação** e busca implementadas
- ✅ **Scripts automatizados** com detecção de sudo
- ✅ **Testes automatizados** do CRUD completo
- ✅ **Hot reload** para desenvolvimento
- ✅ **Docker** containerizado e otimizado

## 🎯 Próximos Passos

### 📦 Próxima Fase - CRUD de Produtos
- [ ] Service layer para produtos
- [ ] Handlers HTTP para produtos
- [ ] Interface web para produtos
- [ ] Gestão de categorias
- [ ] Upload de imagens

### 📊 Fase Seguinte - Controle de Estoque
- [ ] Service layer para estoque
- [ ] Movimentações de estoque
- [ ] Relatórios de estoque
- [ ] Alertas de estoque baixo

## 🔒 Segurança

- **Usuário não-root** no Docker
- **Validações de entrada** em todos endpoints
- **Sanitização** de dados SQL
- **Headers de segurança** configurados
- **Timeouts** apropriados

## 🐛 Solução de Problemas

### Erro de Permissão Docker
```bash
sudo usermod -aG docker $USER
newgrp docker
# ou use os scripts que detectam sudo automaticamente
```

### PostgreSQL não conecta
```bash
sudo docker-compose down
sudo docker-compose up -d postgres
# aguarde 30 segundos
```

### Aplicação não compila
```bash
go mod tidy
go clean -cache
go build .
```

### Porta 8080 em uso
```bash
sudo lsof -i :8080
kill <PID>
```

## 📞 Início Rápido

Para executar o projeto:

1. **🎆 Automático (Recomendado):** `./start_and_test.sh`
2. **🔥 Desenvolvimento:** `./dev_mode.sh`
3. **🛑 Parar tudo:** `./stop_system.sh`
4. **📚 Ajuda completa:** `./help.sh`

**Acessos:**
- **Home:** http://localhost:8080
- **Clientes:** http://localhost:8080/clientes
- **API:** http://localhost:8080/api/clientes
- **Health:** http://localhost:8080/health

---

**🚀 Sistema de Gestão v1.0 - Desenvolvido com Go, PostgreSQL e Docker**