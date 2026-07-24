package cart

type cartRepository interface {
	AddItem(userID, skuID int64, count uint16)
	DeleteItem(userID, skuID int64)
	Clear(userID int64)
	GetItems(userID int64) map[int64]uint16
}
type Service struct {
	repo cartRepository
}

func NewService(repo cartRepository) *Service {
	return &Service{repo: repo}
}
func (s *Service) AddItem(userID, skuID int64, count uint16) error {
	s.repo.AddItem(userID, skuID, count)
	return nil
}
func (s *Service) DeleteItem(userID, skuID int64) {
	s.repo.DeleteItem(userID, skuID)
}

func (s *Service) Clear(userID int64) {
	s.repo.Clear(userID)
}

func (s *Service) GetItems(userID int64) map[int64]uint16 {
	return s.repo.GetItems(userID)
}
