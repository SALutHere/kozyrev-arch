package internal

import (
	"encoding/json"
	"errors"
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
	mux.HandleFunc("DELETE /parts/{id}", h.DeletePart)
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

// DeletePart удаляет деталь
func (h *handler) DeletePart(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}

	if err = h.service.DeletePart(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "деталь не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
