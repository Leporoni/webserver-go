#!/bin/sh

# Script de inicialização para o container

echo "Iniciando servidor Go + Nginx..."

# Iniciar o servidor Go em background
echo "Iniciando aplicação Go na porta 8080..."
./main &

# Aguardar um momento para o servidor Go inicializar
sleep 2

# Iniciar o Nginx
echo "Iniciando Nginx na porta 80..."
nginx -g "daemon off;" &

# Aguardar ambos os processos
wait