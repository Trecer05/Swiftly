# Domain Models

Сущности предметной области и DTO.

## Содержит:
- Структуры БД (с тегами для SQL/JSON)
- Валидационные методы
- Преобразования между форматами

Пример:
```go
type User struct {
    ID    string `json:"id" db:"id"`
    Email string `json:"email" db:"email"`
}
```