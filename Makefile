.PHONY: run-all run-mock test lint check hooks

# .env не читается автоматически: os.Getenv смотрит в окружение процесса.
# set -a включает автоэкспорт, точка (source) исполняет файл в текущем шелле.
run-all:
	set -a && . ./.env && set +a && go run ./cmd/cart

# Моку окружение не нужно: порт берётся из MOCK_ADDR, по умолчанию :8081.
run-mock:
	go run ./cmd/productmock

# -count=1 отключает кеш результатов: тесты должны гоняться по-настоящему.
test:
	go test -race -count=1 ./...

lint:
	gofmt -l .
	go vet ./...

check: lint test

# Подключает хуки из .githooks (в отличие от .git/hooks они версионируются).
# Выполнить один раз после клонирования репозитория.
hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/*
	@echo "хуки подключены: .githooks"