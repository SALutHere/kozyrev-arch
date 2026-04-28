package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	e "github.com/SALutHere/kozyrev-arch/internal/errors"
	"github.com/SALutHere/kozyrev-arch/internal/model"
)

// PartService - интерфейс сервиса для работы с деталями
// Go-идиома: интерфейс определяется потребителем
type PartService interface {
	GetAllParts() []model.Part
	CreatePart(name, partType string, quantity int, weight float64) (model.Part, error)
	WithdrawPart(id int64, quantity int) error
}

// handler - HTTP-обработчики
// Приватная структура: создание только через NewHandler
type handler struct {
	service PartService
}

// NewHandler создаёт новый обработчик
func NewHandler(service PartService) *handler {
	return &handler{service: service}
}

// RegisterRoutes регистрирует маршруты
func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /parts", h.GetParts)
	mux.HandleFunc("POST /parts", h.CreatePart)
	mux.HandleFunc("POST /parts/{id}/withdraw", h.WithdrawPart)
}

// GetParts возвращает список всех деталей
func (h *handler) GetParts(w http.ResponseWriter, r *http.Request) {
	parts := h.service.GetAllParts()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(parts); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

// CreatePart создаёт новую деталь
func (h *handler) CreatePart(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string  `json:"name"`
		Type     string  `json:"type"`
		Quantity int     `json:"quantity"`
		Weight   float64 `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	part, err := h.service.CreatePart(input.Name, input.Type, input.Quantity, input.Weight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(part); err != nil {
		log.Printf("Ошибка кодирования JSON: %v", err)
	}
}

func (h *handler) WithdrawPart(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}

	var input struct {
		Quantity int `json:"quantity"`
	}
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}

	// Структурная валидация
	if input.Quantity <= 0 {
		http.Error(w, "количество должно быть больше 0", http.StatusBadRequest)
		return
	}

	// Бизнес валидация уже произведена в сервисе
	if err = h.service.WithdrawPart(id, input.Quantity); err != nil {
		switch {
		case errors.Is(err, e.ErrNotFound):
			http.Error(w, "деталь не найдена", http.StatusNotFound)
		case errors.Is(err, e.ErrNotEnoughParts):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
