package internal

// partService - бизнес логика для работы с деталями
// Приватная структура: создание только через NewPartService
type partService struct {
	repo *partRepository // конкретный тип, не интерфейс
}

// NewPartService создаёт новый сервис
func NewPartService(repo *partRepository) *partService {
	return &partService{repo: repo}
}

// GetAllParts возвращает все детали
func (s *partService) GetAllParts() []Part {
	return s.repo.GetAll()
}

// CreatePart создаёт новую деталь (без проверок)
func (s *partService) CreatePart(name, partType string, quantity int, weight float64) (Part, error) {
	part := Part{
		Name:     name,
		Type:     partType,
		Quantity: quantity,
		Weight:   weight,
	}
	part = s.repo.Create(part)

	return part, nil
}

// Withdraw списывает детали со склада
// БЕЗ ВАЛИДАЦИИ - валидация пока в handler (это проблема!)
func (s *partService) WithdrawPart(id int64, quantity int) error {
	return s.repo.Withdraw(id, quantity)
}

// GetPartByID возвращает деталь по ID
func (s *partService) GetPartByID(id int64) (Part, error) {
	return s.repo.GetByID(id)
}
