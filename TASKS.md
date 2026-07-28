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

### H3. Роут `/health`
- [ ] В `cart.http` первый запрос — `GET /health`, но в `cmd/cart/main.go` такого
      роута нет → сейчас 404 (подтверждено прогоном). Добавить хендлер или убрать из `cart.http`

### H4. ListCart на пустой корзине отдаёт 200 вместо 404
- [ ] ТЗ (таблица эндпоинтов выше) требует **404** на пустой корзине
- [ ] `internal/handler/handler.go:119-126` отдаёт `200 {"items":[],"total_price":0}`
- [ ] Регрессия появилась в коммите a1619f5 («refactor empty cart responses»)
- [ ] Решить: чинить под ТЗ или сознательно оставить 200

---
## Текущий шаг
Основное задание (T1–T5) закрыто, D1 сделан, D2/D3 сознательно пропущены.
Порядок: T1 ✅ → T2 ✅ → T3 ✅ → T4 ✅ → T5 ✅ → D1 ✅ → **H1 → H2 → H3**.
