.PHONY: run-all run-mock

# .env не читается автоматически: os.Getenv смотрит в окружение процесса.
# set -a включает автоэкспорт, точка (source) исполняет файл в текущем шелле.
run-all:
	set -a && . ./.env && set +a && go run ./cmd/cart

# Моку окружение не нужно: порт берётся из MOCK_ADDR, по умолчанию :8081.
run-mock:
	go run ./cmd/productmock