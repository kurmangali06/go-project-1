package kata

import (
	"context"
	"time"
)

// УРОВЕНЬ 8 — context и отмена.
//
// SlowCatalog имитирует медленный внешний сервис.
//
// Контракт Lookup:
//   - контекст отменён или истёк дедлайн → немедленно вернуть ctx.Err(),
//     не дожидаясь окончания работы
//   - иначе → название товара и nil
//
// Задача: написать тест с context.WithTimeout на срок заведомо меньше
// задержки и убедиться, что Lookup возвращает управление сразу.
// Проверять надо две вещи: какую ошибку вернули (errors.Is с
// context.DeadlineExceeded) и сколько времени это заняло.
//
// Второй тест — на нормальный путь, чтобы первый не оказался
// «зелёным по любой причине».
type SlowCatalog struct {
	Delay time.Duration
	Names map[int64]string
}

func (c SlowCatalog) Lookup(ctx context.Context, sku int64) (string, error) {
	time.Sleep(c.Delay)
	return c.Names[sku], nil
}
