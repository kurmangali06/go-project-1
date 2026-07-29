package kata

import (
	"errors"
	"testing"
)

func TestCatalog_Lookup(t *testing.T) {
	c := NewCatalog(func(sku int64) (string, error) {
		return "", ErrNotFound
	})

	_, err := c.Lookup(123)

	// Контракт: ошибка должна быть. Жалуемся, когда её НЕТ.
	if err == nil {
		t.Fatal("ошибки нет, а контракт требует ErrNotFound")
	}

	// Контракт: errors.Is должен её находить. Жалуемся, когда НЕ находит.
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("errors.Is(err, ErrNotFound) = false, err = %v", err)
	}

	// Та же ошибка, проверенная по тексту сообщения, — пройдёт.
	// Сравни, какая из двух проверок поймала баг.
	if err.Error() != "lookup sku 123: not found" {
		t.Errorf("текст ошибки = %q", err.Error())
	}
}
