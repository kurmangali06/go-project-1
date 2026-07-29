package kata

// УРОВЕНЬ 7 — конкурентность и -race.
//
// ViewCounter считает просмотры карточек товаров.
//
// Контракт:
//   - потокобезопасен: Add и Total можно звать из разных горутин одновременно
//   - после N вызовов Add(sku, 1) значение Total(sku) равно N
//
// Задача: написать тест, который запускает много горутин и падает
// под `go test -race ./internal/kata/`. Обрати внимание: без -race такой
// тест может случайно пройти — гонка не всегда проявляется.
//
// Подсказка по инструментам: sync.WaitGroup, чтобы дождаться всех горутин.
// И помни, что t.Fatal из горутины вызывать нельзя.
type ViewCounter struct {
	views map[int64]int
}

func NewViewCounter() *ViewCounter {
	return &ViewCounter{views: make(map[int64]int)}
}

func (c *ViewCounter) Add(sku int64, n int) {
	c.views[sku] += n
}

func (c *ViewCounter) Total(sku int64) int {
	return c.views[sku]
}
