package main

import (
	"log"
	"net/http"
	"time"

	"github.com/SALutHere/kozyrev-arch/internal/handler"
	"github.com/SALutHere/kozyrev-arch/internal/repository"
	"github.com/SALutHere/kozyrev-arch/internal/service"
)

const inventoryPath = "data/inventory.csv"

func main() {
	// Создание зависимостей
	repo := repository.NewPartRepository()
	service := service.NewPartService(repo)
	handler := handler.NewHandler(service)

	// Загрузка данных из CSV
	if err := repo.LoadFromCSV(inventoryPath); err != nil {
		log.Printf("Не удалось загрузить данные из CSV: %v", err)
	}

	log.Println("Сервер запущен на http://localhost:8080")

	// Регистрация маршрутов
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // максимальное время на чтение всего запроса (заголовок + тело)
		WriteTimeout: 10 * time.Second, // максимальное время на запись ответа клиенту
		IdleTimeout:  60 * time.Second, // максимальное время ожидания следующего запроса при keep-alive соединении
	}
	log.Fatal(server.ListenAndServe())
}
