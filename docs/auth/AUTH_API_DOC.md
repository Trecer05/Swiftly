# Документация по AUTH_API 

## Аутентификация (HTTP)
Базовый URL: `/api/v1`

> **Важно:** В запросах login и register поля `email` и `phone` — взаимоисключающие. Нужно передавать только одно из них.

### POST /api/v1/login
**Описание:** Авторизация пользователя.
**Тело запроса:**
```json
{
  "email": "string,omitempty",
  "number": "string",
  "password": "string"
}
```
**Возможные коды ответов:**
- 200 OK — успешная авторизация
- 400 Bad Request — некорректные данные запроса
- 401 Unauthorized — неверный пароль
- 404 Not Found — пользователь не найден
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "access_token": "string",
  "refresh_token": "string"
}
```
**Примеры ошибок:**
- 401: `{"error": "invalid password"}`
- 404: `{"error": "user not found"}`
- 400: `{"error": "invalid request body"}`

### POST /api/v1/register
**Описание:** Регистрация пользователя.
**Тело запроса:**
```json
{
  "email": "string,omitempty",
  "number": "string",
  "password": "string"
}
```
**Возможные коды ответов:**
- 200 OK — успешная регистрация
- 400 Bad Request — некорректные данные запроса
- 409 Conflict — email или телефон уже зарегистрированы
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "access_token": "string",
  "refresh_token": "string"
}
```
**Примеры ошибок:**
- 409: `{"error": "email already exists"}`
- 409: `{"error": "phone already exists"}`
- 400: `{"error": "invalid request body"}`

### POST /api/v1/logout
**Описание:** Выход пользователя (инвалидация refresh токена).
**Возможные коды ответов:**
- 200 OK — успешный выход
- 401 Unauthorized — пользователь не авторизован или refresh токен истёк
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "logout": "success"
}
```
**Примеры ошибок:**
- 401: `{"error": "unauthorized"}`
- 401: `{"error": "access token expired"}` - запрос на /refresh
- 401: `{"error": "refresh token expired"}` - тогда сразу выход пользователя

### POST /api/v1/refresh
**Описание:** Обновление access токена по refresh токену.
**Тело запроса:**
```json
{
  "refresh_token": "string"
}
```
**Возможные коды ответов:**
- 200 OK — успешное обновление токена
- 400 Bad Request — некорректные данные запроса
- 401 Unauthorized — refresh токен истёк или невалиден
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "access_token": "string",
  "refresh_token": "string"
}
```
**Примеры ошибок:**
- 401: `{"error": "refresh token expired"}`
- 400: `{"error": "invalid request body"}`

---

## Примечания
- Для всех HTTP-запросов Content-Type: application/json.
- В login и register может быть либо email, либо phone (одно из двух).
- Возможны дополнительные поля в моделях пользователя (см. исходный код моделей). 
- /refresh и /logout - защищенные маршруты bearer аутентификацией в header bearer добавляется access токен