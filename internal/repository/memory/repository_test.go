package memory

import (
	"sync"
	"testing"
)

func TestAddItem_SumsCount(t *testing.T) {
	r := NewRepository()

	if err := r.AddItem(22, 2008, 3); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := r.AddItem(22, 2008, 2); err != nil {
		t.Fatalf("AddItem: %v", err)
	}

	items := r.GetItems(22)
	if items[2008] != 5 {
		t.Errorf("count = %d, ожидали 5 (3+2)", items[2008])
	}
	if len(items) != 1 {
		t.Errorf("позиций = %d, ожидали 1: повторный AddItem не должен плодить записи", len(items))
	}
}

func TestAddItem_SeparateUsers(t *testing.T) {
	r := NewRepository()

	if err := r.AddItem(22, 2008, 3); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := r.AddItem(33, 2008, 7); err != nil {
		t.Fatalf("AddItem: %v", err)
	}

	if got := r.GetItems(22)[2008]; got != 3 {
		t.Errorf("корзина 22: count = %d, ожидали 3", got)
	}
	if got := r.GetItems(33)[2008]; got != 7 {
		t.Errorf("корзина 33: count = %d, ожидали 7", got)
	}
}

func TestDeleteItem_Missing(t *testing.T) {
	r := NewRepository()

	if err := r.DeleteItem(22, 2008); err != nil {
		t.Fatalf("удаление несуществующего sku вернуло ошибку: %v", err)
	}
	if items := r.GetItems(22); len(items) != 0 {
		t.Errorf("корзина не пуста: %v", items)
	}
}

func TestDeleteItem_Existing(t *testing.T) {
	r := NewRepository()

	if err := r.AddItem(22, 2008, 3); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := r.AddItem(22, 5000, 1); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := r.DeleteItem(22, 2008); err != nil {
		t.Fatalf("DeleteItem: %v", err)
	}

	items := r.GetItems(22)
	if _, ok := items[2008]; ok {
		t.Error("sku 2008 остался в корзине после удаления")
	}
	if items[5000] != 1 {
		t.Errorf("sku 5000: count = %d, ожидали 1: удаление задело чужую позицию", items[5000])
	}
}

func TestClear(t *testing.T) {
	r := NewRepository()

	if err := r.AddItem(22, 2008, 3); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := r.AddItem(22, 2009, 4); err != nil {
		t.Fatalf("AddItem: %v", err)
	}
	if err := r.Clear(22); err != nil {
		t.Fatalf("Clear: %v", err)
	}

	if items := r.GetItems(22); len(items) != 0 {
		t.Errorf("после Clear корзина не пуста: %v", items)
	}
}

// GetItems обязан возвращать копию: если отдать внутреннюю мапу, вызывающий код
// сможет менять хранилище в обход мьютекса.
func TestGetItems_ReturnsCopy(t *testing.T) {
	r := NewRepository()

	if err := r.AddItem(22, 2008, 3); err != nil {
		t.Fatalf("AddItem: %v", err)
	}

	got := r.GetItems(22)
	got[2008] = 999
	got[7777] = 1

	fresh := r.GetItems(22)
	if fresh[2008] != 3 {
		t.Errorf("count = %d, ожидали 3: правка результата GetItems протекла в хранилище", fresh[2008])
	}
	if _, ok := fresh[7777]; ok {
		t.Error("в хранилище появился sku 7777, добавленный в копию")
	}
}

// Без мьютекса `+=` — это чтение, сложение и запись; часть инкрементов терялась бы.
// Запускать вместе с -race: go test -race ./...
func TestAddItem_Concurrent(t *testing.T) {
	r := NewRepository()

	const goroutines = 100

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.AddItem(1007, 2008, 1); err != nil {
				t.Errorf("AddItem: %v", err) // t.Fatalf из горутины вызывать нельзя
			}
		}()
	}
	wg.Wait()

	if got := r.GetItems(1007)[2008]; got != goroutines {
		t.Errorf("count = %d, ожидали %d: потеряны инкременты", got, goroutines)
	}
}

// Читатели и писатели одновременно — проверяем RWMutex в GetItems.
func TestRepository_ConcurrentReadWrite(t *testing.T) {
	r := NewRepository()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			if err := r.AddItem(1007, 2008, 1); err != nil {
				t.Errorf("AddItem: %v", err)
			}
		}()
		go func() {
			defer wg.Done()
			for sku, count := range r.GetItems(1007) {
				_, _ = sku, count
			}
		}()
	}
	wg.Wait()
}
