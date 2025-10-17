#!/bin/bash

echo "🔄 Reconstruindo e testando o projeto..."

# Parar containers existentes
echo "🛑 Parando containers existentes..."
docker-compose down

# Remover imagens antigas para forçar rebuild
echo "🗑️  Removendo imagens antigas..."
docker rmi webserver-go-webserver 2>/dev/null || true

# Rebuild e start
echo "🔨 Reconstruindo e iniciando containers..."
docker-compose up --build -d

# Aguardar containers iniciarem
echo "⏳ Aguardando containers iniciarem..."
sleep 10

# Verificar status dos containers
echo "📊 Status dos containers:"
docker-compose ps

# Testar health check
echo "🧪 Testando health check..."
for i in {1..10}; do
    echo "Tentativa $i/10..."
    if curl -s http://localhost/health | grep -q "healthy"; then
        echo "✅ Health check passou!"
        break
    elif [ $i -eq 10 ]; then
        echo "❌ Health check falhou após 10 tentativas"
        echo "📋 Logs do container:"
        docker-compose logs webserver
        exit 1
    fi
    sleep 3
done

echo "✅ Projeto reconstruído e testado com sucesso!"
echo "🌐 Acesse: http://localhost"
echo "🔍 Health: http://localhost/health"