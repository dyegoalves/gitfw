# GitFW CLI - Makefile for macOS

BINARY_NAME=gitfw
VERSION=0.1.8-beta
INSTALL_PATH=/usr/local/bin

.PHONY: all build install uninstall clean help

all: build

build:
	@echo "🔨 Construindo $(BINARY_NAME) v$(VERSION)..."
	go build -o $(BINARY_NAME) main.go

install: build
	@echo "🚀 Instalando $(BINARY_NAME) no sistema ($(INSTALL_PATH))..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Instalação concluída! Tente rodar 'gitfw help'"

uninstall:
	@echo "🗑️ Removendo $(BINARY_NAME) do sistema..."
	@sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Removido com sucesso."

clean:
	@echo "🧹 Limpando arquivos temporários..."
	rm -f $(BINARY_NAME)

help:
	@echo "GitFW CLI Makefile"
	@echo "Comandos:"
	@echo "  make build     - Compila o binário localmente"
	@echo "  make install   - Instala o binário globalmente no macOS (requer sudo)"
	@echo "  make uninstall - Remove o binário do sistema"
	@echo "  make clean     - Remove o binário local"
