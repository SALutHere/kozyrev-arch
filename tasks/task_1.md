# Акт 1: Разделение Handler/Service/Repository

**Ветка:** `act_1`

## Задача

CLI-утилита понравилась, но кладовщик хочет веб-интерфейс.

**Вводные:**
- Нужен HTTP API вместо CLI
- CRUD операции с деталями
- Данные не должны теряться между запросами

**API endpoints:**
- `GET /parts` — список всех деталей
- `POST /parts` — создать деталь
- `DELETE /parts/{id}` — удалить деталь

## Пример использования

```bash
# Получить все детали
curl http://localhost:8080/parts

# Создать деталь
curl -X POST http://localhost:8080/parts \
  -H "Content-Type: application/json" \
  -d '{"name":"Ионный двигатель","type":"engine","quantity":10,"weight":420.0}'

# Удалить деталь
curl -X DELETE http://localhost:8080/parts/5
```

## Ключевая мысль

> Появился HTTP → появились слои.
> Handler/Service/Repository — минимальное разделение ответственности.

```
HTTP Request
     ↓
┌─────────────┐
│   Handler   │  ← парсинг JSON, HTTP ответы
├─────────────┤
│   Service   │  ← бизнес-логика (пока пустая)
├─────────────┤
│ Repository  │  ← хранение данных (in-memory)
└─────────────┘
```

Каждый слой отвечает за своё. Можно заменить in-memory на PostgreSQL — Service не изменится.