# Акт 2: Интерфейсы и изоляция бизнес-логики

**Ветка:** `act_2`

## Задача

Кто-то списал 500 двигателей в минус. На складе оказалось -327 штук.

**Вводные:**
- Нельзя списать больше, чем есть на складе
- Понятные сообщения об ошибках
- Бизнес-логика должна быть тестируемой без HTTP

**Новый endpoint:**
- `POST /parts/{id}/withdraw` — списать детали со склада

## Пример использования

```bash
# Списать 5 деталей (успех)
curl -X POST http://localhost:8080/parts/1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"quantity": 5}'
# → 200 OK

# Списать больше, чем есть (ошибка)
curl -X POST http://localhost:8080/parts/1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"quantity": 9999}'
# → 400: недостаточно деталей: доступно 12, запрошено 9999

# Деталь не существует
curl -X POST http://localhost:8080/parts/999/withdraw \
  -H "Content-Type: application/json" \
  -d '{"quantity": 10}'
# → 404: деталь не найдена

# Некорректное количество
curl -X POST http://localhost:8080/parts/1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"quantity": 0}'
# → 400: количество должно быть больше 0
```

## Ключевая мысль

> Появилась бизнес-логика → появились интерфейсы.
> Интерфейс определяется потребителем (Go-идиома).

**Два типа валидации:**
- **Handler** — структурная (quantity > 0, JSON валидный)
- **Service** — бизнес-логика (хватает ли деталей)

```go
// Service определяет интерфейс с нужными методами
type PartRepository interface {
    GetByID(id int64) (Part, error)
    Withdraw(id int64, qty int) error
}

type PartService struct {
    repo PartRepository  // интерфейс, не конкретный тип
}
```

Теперь Service можно тестировать с моком, без реальной БД.