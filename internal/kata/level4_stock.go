package kata

import "errors"

var ErrNotEnough = errors.New("not enough stock")

// УРОВЕНЬ 4 — побочные эффекты.
//
// Stock — остатки товаров на складе.
//
// Контракт Reserve:
//   - остатка хватает → остаток уменьшается на count, возвращается nil
//   - остатка не хватает → ErrNotEnough, и остаток НЕ меняется
//   - count <= 0 → ErrNotEnough, остаток не меняется
//
// Задача: проверять не только возвращённую ошибку, но и состояние объекта
// после вызова. Тест, который смотрит только на err, здесь ничего не поймает.
type Stock struct {
	items map[int64]int
}

func NewStock(items map[int64]int) *Stock {
	cp := make(map[int64]int, len(items))
	for sku, n := range items {
		cp[sku] = n
	}
	return &Stock{items: cp}
}

func (s *Stock) Left(sku int64) int {
	return s.items[sku]
}

func (s *Stock) Reserve(sku int64, count int) error {
	s.items[sku] -= count
	if s.items[sku] < 0 {
		return ErrNotEnough
	}
	return nil
}
