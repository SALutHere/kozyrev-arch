package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const inventoryPath = "data/inventory.csv"

// Part - деталь со склада
type Part struct {
	ID       int64
	Name     string
	Type     string
	Quantity int
	Weight   float64
}

// TypeStats - статистика по типу деталей
type TypeStats struct {
	Count       int
	TotalWeight float64
}

func main() {
	parts, err := readInventory(inventoryPath)
	if err != nil {
		log.Fatalf("Ошибка чтения инвентаря: %v", err)
	}

	// Служебный вывод, что загружено
	printLoadInfo(parts)

	stats := calcStatsByType(parts)
	missing := findMissing(parts)

	printStats(stats)
	printMissing(missing)
}

// readInventory читает CSV-файл с инвентарём и возвращает список деталей
func readInventory(path string) ([]Part, error) {
	// Валидация пути для предотвращения path traversal
	cleanPath := filepath.Clean(path)
	if filepath.IsAbs(cleanPath) {
		return nil, fmt.Errorf("недопустимый путь к файлу: %s", path)
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("ошибка при закрытии файла: %w", err)
			os.Exit(1)
		}
	}()

	reader := csv.NewReader(file)

	// Пропускаем заголовок
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("не удалось прочитать заголовок: %w", err)
	}

	var parts []Part
	for {
		var records []string
		if records, err = reader.Read(); err != nil {
			break
		}

		var id int64
		var quantity int
		var weight float64

		id, err = strconv.ParseInt(records[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать ID: %w", err)
		}

		quantity, err = strconv.Atoi(records[3])
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать количество: %w", err)
		}

		weight, err = strconv.ParseFloat(records[4], 64)
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать вес: %w", err)
		}

		parts = append(parts, Part{
			ID:       id,
			Name:     records[1],
			Type:     records[2],
			Quantity: quantity,
			Weight:   weight,
		})
	}

	return parts, nil
}

// printLoadInfo выводит информацию о загруженных данных
func printLoadInfo(parts []Part) {
	// Собираем уникальные типы
	typesSet := make(map[string]struct{})
	for _, p := range parts {
		typesSet[p.Type] = struct{}{}
	}

	types := make([]string, 0, len(typesSet))
	for t := range typesSet {
		types = append(types, t)
	}
	sort.Strings(types)

	log.Printf("Загружено %d деталей из %s\n", len(parts), inventoryPath)
	log.Printf("Типы: %v\n", types)
	log.Println()
}

// calcStatsByType подсчитывает статистику по типам деталей
func calcStatsByType(parts []Part) map[string]TypeStats {
	stats := make(map[string]TypeStats)

	for _, p := range parts {
		s := stats[p.Type]
		s.Count += p.Quantity
		s.TotalWeight += p.Weight * float64(p.Quantity)
		stats[p.Type] = s
	}

	return stats
}

// findMissing находит детали с нулевым остатком
func findMissing(parts []Part) []string {
	var missing []string

	for _, p := range parts {
		if p.Quantity == 0 {
			missing = append(missing, p.Name)
		}
	}

	return missing
}

// printStats выводит статистику по типам
func printStats(stats map[string]TypeStats) {
	log.Println("=== Статистика склада ===")
	log.Println()
	log.Println("По типам:")

	// Сортируем типы для стабильного вывода
	types := make([]string, 0, len(stats))
	for t := range stats {
		types = append(types, t)
	}
	sort.Strings(types)

	var totalWeight float64
	for _, t := range types {
		s := stats[t]
		log.Printf(" %-10s %3d шт, %8.1f кг\n", t+":", s.Count, s.TotalWeight)
		totalWeight += s.TotalWeight
	}

	log.Println()
	log.Printf("Общий вес: %.1f кг\n", totalWeight)
}

// printMissing выводит список отсутствующих деталей
func printMissing(missing []string) {
	log.Println()
	log.Println("=== Отсутствуют на складе ===")

	if len(missing) == 0 {
		log.Println("  Все детали в наличии!")
		return
	}

	for _, name := range missing {
		log.Printf("  %s\n", name)
	}
}
