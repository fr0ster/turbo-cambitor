# Makefile for Turbo Cambitor
# Автоматизація Git та збірки

.PHONY: help build clean test deps install-deps branch feature-branch

# Змінні
GO = go
VERSION = $(shell git describe --tags --always --dirty)

# Кольори для виводу
GREEN = \033[32m
YELLOW = \033[33m
RED = \033[31m
BLUE = \033[34m
RESET = \033[0m

# Допомога
help:
	@echo "$(BLUE)🚀 Turbo Cambitor Makefile$(RESET)"
	@echo ""
	@echo "$(GREEN)Доступні команди:$(RESET)"
	@echo "  build         - Збірка проекту"
	@echo "  clean         - Очищення"
	@echo "  test          - Запуск тестів"
	@echo "  deps          - Перевірка залежностей"
	@echo "  install-deps  - Встановлення залежностей"
	@echo "  branch        - Показати поточний бранч"
	@echo "  feature-branch - Створити feature бранч"
	@echo "  help          - Показати цю довідку"
	@echo ""
	@echo "$(GREEN)Короткі команди:$(RESET)"
	@echo "  b             - build"
	@echo "  c             - clean"
	@echo "  t             - test"
	@echo "  h             - help"

# Збірка проекту
build:
	@echo "$(GREEN)🔨 Building turbo-cambitor...$(RESET)"
	@$(GO) build ./...
	@echo "$(GREEN)✅ Build completed!$(RESET)"

# Очищення
clean:
	@echo "$(YELLOW)🧹 Cleaning...$(RESET)"
	@$(GO) clean
	@echo "$(GREEN)✅ Clean completed!$(RESET)"

# Запуск тестів
test:
	@echo "$(YELLOW)🧪 Running tests...$(RESET)"
	@$(GO) test ./...
	@echo "$(GREEN)✅ Tests completed!$(RESET)"

# Перевірка залежностей
deps:
	@echo "$(YELLOW)🔍 Checking dependencies...$(RESET)"
	@$(GO) mod verify
	@$(GO) list -m all
	@echo "$(GREEN)✅ Dependencies checked!$(RESET)"

# Встановлення залежностей
install-deps:
	@echo "$(YELLOW)📦 Installing dependencies...$(RESET)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "$(GREEN)✅ Dependencies installed!$(RESET)"

# Показати поточний бранч
branch:
	@echo "$(BLUE)🌿 Current branch:$(RESET)"
	@git branch --show-current
	@echo "$(BLUE)📊 Git status:$(RESET)"
	@git status --short

# Створити feature бранч
feature-branch:
	@echo "$(YELLOW)🌿 Creating feature branch...$(RESET)"
	@git checkout -b feature/update-to-turbo-signer-v2
	@echo "$(GREEN)✅ Feature branch created: feature/update-to-turbo-signer-v2$(RESET)"
	@echo "$(BLUE)📊 Current branch:$(RESET)"
	@git branch --show-current

# Короткі команди
b: build
c: clean
t: test
h: help

# Default target
.DEFAULT_GOAL := help
