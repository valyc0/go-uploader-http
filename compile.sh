#!/bin/bash

# Nome del file di output (senza estensione)
OUTPUT_NAME="webserver"

# Compila per Linux
echo "Compilazione per Linux..."
GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_NAME}_linux main.go
if [ $? -eq 0 ]; then
    echo "Compilazione per Linux completata: ${OUTPUT_NAME}_linux"
else
    echo "Errore durante la compilazione per Linux"
    exit 1
fi

# Compila per Windows
echo "Compilazione per Windows..."
GOOS=windows GOARCH=amd64 go build -o ${OUTPUT_NAME}_windows.exe main.go
if [ $? -eq 0 ]; then
    echo "Compilazione per Windows completata: ${OUTPUT_NAME}_windows.exe"
else
    echo "Errore durante la compilazione per Windows"
    exit 1
fi

echo "Compilazione terminata con successo!"

