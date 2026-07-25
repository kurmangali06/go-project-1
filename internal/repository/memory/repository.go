package memory

import "sync"

type Repository struct {
	mu      sync.RWMutex
	storage map[int64]map[int64]uint16
}

func NewRepository() *Repository {
	return &Repository{
		storage: make(map[int64]map[int64]uint16),
	}
}

func (r *Repository) AddItem(userID, skuID int64, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[userID]; !ok {
		r.storage[userID] = make(map[int64]uint16)
	}

	r.storage[userID][skuID] += count
	return nil
}
func (r *Repository) DeleteItem(userID, skuID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.storage[userID]; !ok {
		return nil
	}

	delete(r.storage[userID], skuID)
	if len(r.storage[userID]) == 0 {
		delete(r.storage, userID)
	}
	return nil
}

func (r *Repository) GetItems(userID int64) map[int64]uint16 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cart := r.storage[userID]
	result := make(map[int64]uint16, len(cart))

	for sku, count := range cart {
		result[sku] = count
	}
	return result
}
func (r *Repository) Clear(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.storage, userID)
	return nil
}
