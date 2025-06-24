# Migrations

Структура баз данных каждого микросервиса

## Содержит:
- Структуры БД

## Подпакеты:
- `auth/` - Структура auth микросервиса
- `chat/` - Структура чат сервиса
- Подпакеты будут дополняться

Пример:
```sql
CREATE TABLE users {
    id serial SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
}
```