package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// handler - HTTP-обработчики
// Приватная структура: создание только через NewHandler
type handler struct {
	service *partService
}

// NewHandler создаёт новый обработчик
func NewHandler(service *partService) *handler {
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

	// ПРОБЛЕМА: handler напрямую обращается к service.GetByID и делает бизнес-валидацию!
	part, err := h.service.GetPartByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "деталь не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ПРОБЛЕМА: бизнес-логика в handler!
	if part.Quantity < input.Quantity {
		http.Error(w, fmt.Sprintf("недостаточно деталей: доступно %d, запрошено %d", part.Quantity, input.Quantity), http.StatusBadRequest)
		return
	}

	// ПРОБЛЕМА: handler знает про внутреннюю структуру сервиса
	if err = h.service.WithdrawPart(id, input.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
