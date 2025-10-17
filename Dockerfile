# Multi-stage build para otimizar o tamanho da imagem
FROM golang:1.21-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Imagem final
FROM alpine:latest

# Instalar certificados SSL, nginx e wget para healthcheck
RUN apk --no-cache add ca-certificates nginx wget

# Definir diretório de trabalho
WORKDIR /app

# Copiar binário da aplicação
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/nginx.conf /etc/nginx/nginx.conf

# Criar diretórios necessários para o nginx
RUN mkdir -p /var/log/nginx /var/lib/nginx/tmp /var/cache/nginx /run/nginx

# Criar página de erro personalizada
RUN mkdir -p /usr/share/nginx/html && \
    echo '<h1>Service Temporarily Unavailable</h1>' > /usr/share/nginx/html/50x.html

# Definir permissões
RUN chmod +x /app/main

# Expor portas
EXPOSE 80 8080

# Script de inicialização
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Comando padrão
CMD ["/start.sh"]