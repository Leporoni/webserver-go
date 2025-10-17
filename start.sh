#!/bin/sh

# Script de inicialização para o container

echo "Iniciando servidor Go + Nginx..."

# Verificar se o binário existe e tem permissões
if [ ! -f "/app/main" ]; then
    echo "ERRO: Binário /app/main não encontrado!"
    exit 1
fi

if [ ! -x "/app/main" ]; then
    echo "ERRO: Binário /app/main não tem permissões de execução!"
    chmod +x /app/main
fi

# Iniciar o servidor Go em background
echo "Iniciando aplicação Go na porta 8080..."
cd /app
./main &
GO_PID=$!

# Aguardar o servidor Go inicializar e verificar se está funcionando
echo "Aguardando aplicação Go inicializar..."
sleep 5

# Verificar se o processo Go ainda está rodando
if ! kill -0 $GO_PID 2>/dev/null; then
    echo "ERRO: Aplicação Go falhou ao iniciar!"
    exit 1
fi

# Testar se a aplicação está respondendo
echo "Testando aplicação Go..."
for i in 1 2 3 4 5; do
    if wget -q --spider http://localhost:8080/health; then
        echo "✅ Aplicação Go está respondendo!"
        break
    fi
    echo "Tentativa $i/5: Aguardando aplicação..."
    sleep 2
done

# Iniciar o Nginx
echo "Iniciando Nginx na porta 80..."
nginx -t && nginx -g "daemon off;" &
NGINX_PID=$!

echo "✅ Ambos os serviços iniciados!"
echo "Go PID: $GO_PID"
echo "Nginx PID: $NGINX_PID"

# Função para cleanup
cleanup() {
    echo "Parando serviços..."
    kill $GO_PID $NGINX_PID 2>/dev/null
    exit 0
}

# Capturar sinais para cleanup
trap cleanup TERM INT

# Aguardar ambos os processos
wait