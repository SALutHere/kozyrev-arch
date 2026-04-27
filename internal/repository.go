package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var ErrNotFound = errors.New("деталь не найдена")

// partRepository - in-memory хранилище деталей
// Приватная структура: создание только через NewPartRepository
type partRepository struct {
	mu      sync.Mutex
	storage map[int64]Part
	nextID  int64
}

// NewPartRepository создаёт новый репозиторий
func NewPartRepository() *partRepository {
	return &partRepository{
		storage: make(map[int64]Part),
		nextID:  1,
	}
}

// LoadFromCSV загружает детали из CSV файла
func (r *partRepository) LoadFromCSV(path string) error {
	cleanPath := filepath.Clean(path)
	if filepath.IsAbs(cleanPath) {
		return fmt.Errorf("недопустимый путь к файлу: %s", path)
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("Ошибка закрытия файла: %v", err)
		}
	}()

	reader := csv.NewReader(file)

	// Пропускаем заголовок
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("не удалось прочитать заголовок: %w", err)
	}

	var (
		record   []string
		id       int64
		quantity int
		weight   float64
	)

	for {
		record, err = reader.Read()
		if err != nil {
			break
		}

		id, err = strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			return fmt.Errorf("не удалось преобразовать ID: %w", err)
		}

		quantity, err = strconv.Atoi(record[3])
		if err != nil {
			return fmt.Errorf("не удалось преобразовать количество: %w", err)
		}

		weight, err = strconv.ParseFloat(record[4], 64)
		if err != nil {
			return fmt.Errorf("не удалось преобразовать вес: %w", err)
		}

		r.storage[id] = Part{
			ID:       id,
			Name:     record[1],
			Type:     record[2],
			Quantity: quantity,
			Weight:   weight,
		}

		if id >= r.nextID {
			r.nextID = id + 1
		}
	}

	return nil
}

// GetAll возвращает все детали
func (r *partRepository) GetAll() []Part {
	r.mu.Lock()
	defer r.mu.Unlock()

	parts := make([]Part, 0, len(r.storage))
	for _, p := range r.storage {
		parts = append(parts, p)
	}

	return parts
}

// Create сохраняет новую деталь и возвращает её с присвоением ID
func (r *partRepository) Create(part Part) Part {
	r.mu.Lock()
	defer r.mu.Unlock()

	part.ID = r.nextID
	r.storage[r.nextID] = part
	r.nextID++

	return part
}

// Delete удаляет деталь по ID
func (r *partRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[id]; !ok {
		return ErrNotFound
	}
	delete(r.storage, id)

	return nil
}
