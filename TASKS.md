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
- [ ] Токен из ENV (сейчас передаётся в `New` — подключим при вызове из main)

### T2. AddItem проверяет существование товара ⏳ ТЕКУЩАЯ
- [ ] Перед добавлением звать Product Service; если товара нет → 412 Precondition Failed
- [ ] Пробросить `ProductService` в `cart.Service` через интерфейс

### T3. ListCart отдаёт name/price + total_price
- [ ] По каждому sku тянуть name/price из Product Service
- [ ] В ответ добавить `items[i].name`, `items[i].price`, посчитать `total_price`
- [ ] Формат ответа — как в примере ТЗ (мы остановились ровно здесь)

### T4. Clear → 204 No Content
- [ ] Сейчас возвращает 200 (`handler.go:90`) — по ТЗ должно быть 204

### T5. Финализация cart.http
- [ ] Дописать сценарии под name/price/total_price, поправить рассинхрон sku (5002/5000)

## ДОП. ЗАДАНИЕ (10 баллов)
- [ ] D1. Middleware логирования входящих HTTP-запросов
- [ ] D2. Валидация входящих структур через open-source библиотеку
- [ ] D3. Client Middleware ретраев в Product Service: на 420/429, 3 ретрая, потом ошибка

---
## Текущий шаг
Мы дошли до **total_price** → он упирается в Product Service.
Порядок: **T1 → T2 → T3 → T4 → T5**, затем доп. D1–D3.
