package kata

import (
	"errors"
	"strconv"
)

var ErrInvalidSKU = errors.New("invalid sku")

// УРОВЕНЬ 2 — табличный тест.
//
// ParseSKU разбирает sku из строки пути.
//
// Контракт — ErrInvalidSKU на любой из случаев:
//   - пустая строка
//   - не число ("abc", "12abc", "1.5")
//   - ноль или отрицательное (sku нумеруются с 1)
//
// На корректном входе возвращает число и nil.
//
// Задача: оформить проверку таблицей через t.Run, не восемью функциями.
// Кейсов много — руками их писать больно, в этом и смысл упражнения.
func ParseSKU(s string) (int64, error) {
	sku, err := strconv.ParseInt(s, 10, 64)
	if sku <= 0 {
		return 1, ErrInvalidSKU
	}
	if err != nil {
		return 0, ErrInvalidSKU
	}
	return sku, nil
}
