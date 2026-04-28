package service

import (
	e "github.com/SALutHere/kozyrev-arch/internal/errors"
	"github.com/SALutHere/kozyrev-arch/internal/model"
)

// PartRepository - интерфейс для работы с хранилищем деталей
// Go-идиома: интерфейс лежит там, где используется
type PartRepository interface {
	GetByID(id int64) (model.Part, error)
	GetAll() []model.Part
	Create(part model.Part) model.Part
	Withdraw(id int64, quantity int) error
}

// partService - бизнес логика для работы с деталями
// Приватная структура: создание только через NewPartService
type partService struct {
	repo PartRepository // интерфейс, не конкретный тип
}

// NewPartService создаёт новый сервис
func NewPartService(repo PartRepository) *partService {
	return &partService{repo: repo}
}

// GetAllParts возвращает все детали
func (s *partService) GetAllParts() []model.Part {
	return s.repo.GetAll()
}

// CreatePart создаёт новую деталь (без проверок)
func (s *partService) CreatePart(name, partType string, quantity int, weight float64) (model.Part, error) {
	part := model.Part{
		Name:     name,
		Type:     partType,
		Quantity: quantity,
		Weight:   weight,
	}
	part = s.repo.Create(part)

	return part, nil
}

// Withdraw списывает детали со склада
func (s *partService) WithdrawPart(id int64, quantity int) error {
	part, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if part.Quantity < quantity {
		return e.ErrNotEnoughParts
	}

	return s.repo.Withdraw(id, quantity)
}
