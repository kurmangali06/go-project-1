# Cart Service — бэклог (модуль «Основы Go»)

> Сервис корзины пользователя. HTTP на стандартной библиотеке Go 1.22, in-memory
> хранилище, поход во внешний Product Service за товаром/ценой.
> **Дедлайн сдачи:** 1 июня 23:59 · **проверка:** 4 июня 23:59.

## Архитектура (слои)
```
cmd/cart/main.go                         — точка входа, роутинг, порт :8082
internal/handler/handler.go              — HTTP-слой
internal/service/cart/service.go         — бизнес-логика
internal/repository/memory/repository.go — in-memory хранилище
internal/clients/product/                — клиент Product Service (НЕ создан)
```

## Эндпоинты (порт :8082)
| Метод  | URI                             | Успех | Пусто/нет |
|--------|---------------------------------|-------|-----------|
| POST   | `/user/{user_id}/cart/{sku_id}` | 200   | —         |
| DELETE | `/user/{user_id}/cart/{sku_id}` | 200   | 200       |
| DELETE | `/user/{user_id}/cart`          | 204   | 204       |
| GET    | `/user/{user_id}/cart/list`     | 200   | 404       |

## Product Service
Swagger: http://route256.pavl.uk:8080/docs/
- `get_product` — req `{token string, sku uint32}` → resp `{name string, price uint32}`
- `list_skus`   — req `{token, startAfterSku, count}` → resp `{skus []uint32}`

---

## СДЕЛАНО ✅
- [x] Каркас проекта, `Makefile` (`run-all`), порт :8082
- [x] In-memory репозиторий (map + RWMutex), AddItem суммирует count
- [x] 4 эндпоинта + роутинг, слои через интерфейсы
- [x] AddItem: count из тела JSON + валидация `count == 0` → 400
- [x] Парсинг ID через общий `checkID`, 400 на кривых id/body
- [x] ListCart: 404 на пустой корзине, items сортируются по sku ↑

## ОСНОВНОЕ ЗАДАНИЕ — осталось

### T1. Клиент Product Service ✅
- [x] Пакет `internal/clients/product`: `GetProduct(ctx, sku) (Product, error)`
- [x] POST на `get_product`, таймаут 2s, `ErrProductNotFound` на 404
- [x] Токен из ENV (`main.go`: `os.Getenv("PRODUCT_TOKEN")`)

### T2. AddItem проверяет существование товара ✅
- [x] Перед добавлением звать Product Service; если товара нет → 412 Precondition Failed
- [x] Пробросить `ProductService` в `cart.Service` через интерфейс

### T3. ListCart отдаёт name/price + total_price ✅
- [x] По каждому sku тянуть name/price из Product Service
- [x] В ответ добавить `items[i].name`, `items[i].price`, посчитать `total_price`
- [x] Формат ответа — как в примере ТЗ

### T4. Clear → 204 No Content ✅
- [x] `handler.go: Clear` возвращает 204 вместо 200

### T5. Финализация cart.http ✅
- [x] Дописаны ожидания name/price/total_price (пункты 4, 6), 204 для Clear (пункт 7)
- [x] Добавлен сценарий 412 на неизвестном sku (пункт 13); рассинхрон sku не найден — уже был 5000 везде

## ДОП. ЗАДАНИЕ (10 баллов)
- [x] D1. Middleware логирования входящих HTTP-запросов
      (`internal/middleware/logger.go`, подключено в `cmd/cart/main.go`)
- [~] D2. Валидация входящих структур через open-source библиотеку — **пропущено**
- [~] D3. Client Middleware ретраев в Product Service (420/429, 3 ретрая) — **пропущено**

## ХВОСТЫ (гигиена проекта)

### H1. Убрать `.env` из индекса git
- [ ] `.env` уже закоммичен в `origin/main` (коммит 59256bf) — `.gitignore` его
      не спасает, потому что gitignore действует только на **неотслеживаемые** файлы
- [ ] Убрать из индекса, оставив файл на диске; закоммитить
- [ ] Проверить: `git ls-files | grep env` → пусто
- [ ] Завести `.env.example` с ключами без значений, чтобы было понятно, что настраивать

### H2. Починить `cmd/productmock` → `make run-mock` ✅
- [x] `cmd/productmock/main.go`: `POST /get_product`, каталог 2008/5000, 404 на остальное
- [x] Порт через `MOCK_ADDR`, по умолчанию **:8081** — 8080 на машине занят node-процессом
- [x] `.env` / `.env.example` / шапка `cart.http` переведены на :8081
- [x] Прогнан end-to-end: 200 / 412 / total_price=5300 / 204 — всё сходится с ТЗ

### H3. Роут `/health` ✅
- [x] Решено убрать `GET /health` из `cart.http` (роута в сервисе нет и не требуется по ТЗ)

### H4. ListCart на пустой корзине отдаёт 200 вместо 404 ✅
- [x] `handler.go:119-121` теперь возвращает 404 (регрессия из a1619f5 закрыта)
- [ ] Мелочь: текст ошибки `"product not found"` → `"cart is empty"`

## T6. ТЕСТЫ (в проекте нет ни одного `*_test.go`)

### T6.1 `internal/handler/handler_test.go`
- [x] Стаб `cartService` с полями-функциями, хелпер `listRequest` (`SetPathValue`)
- [x] `TestListCart_Empty` → 404, `TestListCart_NotEmpty` → 200 + total_price, `TestListCart_BadUserID` → 400
- [ ] `TestAddItem`: таблица — `count=0` → 400, битый JSON → 400,
      `cart.ErrProductNotFound` → 412, успех → 200

### T6.2 `internal/service/cart/service_test.go` ✅
- [x] Стаб `productService` на `atomic.Int64` (GetItems ходит параллельно), реальный `memory.Repository`
- [x] 7 тестов: AddItem неизвестного sku → `ErrProductNotFound` + ничего не записано;
      сетевая ошибка НЕ подменяется на ErrProductNotFound; GetItems считает total_price,
      сортирует по sku ↑, ходит ровно N раз; пустая корзина; ошибка Product Service
      пробрасывается и данные не отдаются; DeleteItem/Clear
- [x] Мутационная проверка: убрал `sort.Slice` → падает; сломал трансляцию ошибок → падает

### T6.3 `internal/repository/memory/repository_test.go` ✅
- [x] 8 тестов: суммирование count, изоляция корзин разных юзеров, удаление
      существующего/несуществующего sku, Clear, GetItems возвращает копию
- [x] `TestAddItem_Concurrent` + `TestRepository_ConcurrentReadWrite`, `go test -race ./...` чист
- [x] Мутационная проверка: с убранным `mu.Lock()` тесты падают с DATA RACE — значит не холостые

---
## ИТОГ СЕССИИ 28–29.07.2026

**Состояние:** `go vet` чист, `go test -race ./...` зелёный, 18 тестов в 3 пакетах.

| Пакет | Покрытие |
|---|---|
| `internal/service/cart` | 100.0% |
| `internal/repository/memory` | 96.2% |
| `internal/handler` | 26.6% |
| `internal/clients/product` | 0% |
| `internal/middleware` | 0% |

**Осталось (по убыванию важности):**
- [ ] `TestAddItem` в `internal/handler` таблицей — подтянет покрытие с 26.6%
- [ ] Тесты `internal/clients/product` через `httptest.NewServer`
- [ ] Текст `"product not found"` → `"cart is empty"` в `handler.go:120`
- [ ] D2 (валидация) и D3 (ретраи) — пропущены сознательно

---
## Текущий шаг
Основное задание (T1–T5) закрыто, D1 сделан, D2/D3 сознательно пропущены.
Порядок: T1 ✅ → T2 ✅ → T3 ✅ → T4 ✅ → T5 ✅ → D1 ✅ → **H1 → H2 → H3**.
