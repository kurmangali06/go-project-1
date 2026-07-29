package kata

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

// УРОВЕНЬ 3 — errors.Is и обёртки.
//
// Catalog.Lookup ищет название товара во внешнем источнике.
//
// Контракт:
//   - товара нет      → ошибка, для которой errors.Is(err, ErrNotFound) == true
//   - источник упал   → ошибка источника должна оставаться доступной
//     через errors.Is, а к сообщению добавляется контекст с sku
//   - товар найден    → название и nil
//
// Задача: проверять ошибки через errors.Is, а не сравнением == и не по
// тексту сообщения. Тест должен показать, чем эти способы отличаются.
type Catalog struct {
	source func(sku int64) (string, error)
}

func NewCatalog(source func(sku int64) (string, error)) *Catalog {
	return &Catalog{source: source}
}

func (c *Catalog) Lookup(sku int64) (string, error) {
	name, err := c.source(sku)
	if err != nil {
		return "", fmt.Errorf("lookup sku %d: %w", sku, err)
	}
	return name, nil
}
