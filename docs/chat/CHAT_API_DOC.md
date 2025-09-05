# Документация по CHAT API

## Чат-система (HTTP + WebSocket)
Базовый URL: `/api/v1`

> **Важно:** Все эндпойнты требуют аутентификации через Bearer токен в заголовке Authorization.

---

## HTTP API

### Пользователи

#### POST /api/v1/user
**Описание:** Создание профиля пользователя в чат-системе.
**Content-Type:** `multipart/form-data`
**Параметры:**
- `json` (form field) - JSON с данными пользователя
- `photo` (file, optional) - фото профиля

**JSON данные:**
```json
{
  "username": "string",
  "name": "string", 
  "description": "string"
}
```

**Возможные коды ответов:**
- 200 OK — пользователь успешно создан
- 400 Bad Request — некорректные данные запроса
- 401 Unauthorized — пользователь не авторизован
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "status": "ok"
}
```

#### GET /api/v1/user
**Описание:** Получение информации о текущем пользователе.
**Возможные коды ответов:**
- 200 OK — информация получена
- 401 Unauthorized — пользователь не авторизован
- 404 Not Found — пользователь не найден
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "id": 1,
  "name": "string",
  "username": "string",
  "description": "string"
}
```

#### GET /api/v1/user/{id}
**Описание:** Получение информации о другом пользователе.
**Параметры URL:**
- `id` - ID пользователя

**Возможные коды ответов:**
- 200 OK — информация получена
- 400 Bad Request — некорректный ID
- 404 Not Found — пользователь не найден
- 500 Internal Server Error — внутренняя ошибка сервера

---

### Приватные чаты

#### POST /api/v1/chat/{id}
**Описание:** Создание приватного чата с пользователем.
**Параметры URL:**
- `id` - ID пользователя для создания чата

**Возможные коды ответов:**
- 200 OK — чат создан
- 400 Bad Request — некорректный ID пользователя
- 401 Unauthorized — пользователь не авторизован
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "status": "ok",
  "chat": {
    "id": 1
  }
}
```

#### GET /api/v1/chat/{id}/info
**Описание:** Получение информации о приватном чате.
**Параметры URL:**
- `id` - ID чата

**Возможные коды ответов:**
- 200 OK — информация получена
- 400 Bad Request — некорректный ID
- 404 Not Found — чат не найден
- 500 Internal Server Error — внутренняя ошибка сервера

---

### Групповые чаты

#### POST /api/v1/group
**Описание:** Создание группового чата.
**Content-Type:** `multipart/form-data`
**Параметры:**
- `json` (form field) - JSON с данными группы
- `photo` (file, optional) - фото группы

**JSON данные:**
```json
{
  "name": "string",
  "description": "string",
  "users": [
    {
      "id": 1,
      "name": "string"
    }
  ]
}
```

**Возможные коды ответов:**
- 200 OK — группа создана
- 400 Bad Request — некорректные данные
- 401 Unauthorized — пользователь не авторизован
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "status": "ok",
  "group": {
    "id": 1,
    "name": "string"
  }
}
```

#### DELETE /api/v1/group/{id}
**Описание:** Удаление группы (только владелец).
**Параметры URL:**
- `id` - ID группы

**Возможные коды ответов:**
- 200 OK — группа удалена
- 400 Bad Request — некорректный ID
- 401 Unauthorized — пользователь не авторизован
- 403 Forbidden — пользователь не владелец группы
- 500 Internal Server Error — внутренняя ошибка сервера

#### GET /api/v1/group/{id}/info
**Описание:** Получение информации о группе.
**Параметры URL:**
- `id` - ID группы

**Возможные коды ответов:**
- 200 OK — информация получена
- 400 Bad Request — некорректный ID
- 404 Not Found — группа не найдена
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "id": 1,
  "name": "string",
  "description": "string"
}
```

#### GET /api/v1/group/{id}/users
**Описание:** Получение списка пользователей группы.
**Параметры URL:**
- `id` - ID группы

**Возможные коды ответов:**
- 200 OK — список получен
- 400 Bad Request — некорректный ID
- 404 Not Found — пользователи не найдены
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "users": [
    {
      "id": 1,
      "name": "string"
    }
  ]
}
```

#### POST /api/v1/group/{id}/add
**Описание:** Добавление пользователей в группу (только владелец).
**Параметры URL:**
- `id` - ID группы

**Тело запроса:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "string"
    }
  ]
}
```

**Возможные коды ответов:**
- 200 OK — пользователи добавлены
- 400 Bad Request — некорректные данные
- 401 Unauthorized — пользователь не авторизован
- 403 Forbidden — пользователь не владелец группы
- 409 Conflict — пользователь уже в группе
- 500 Internal Server Error — внутренняя ошибка сервера

#### DELETE /api/v1/group/{id}/delete
**Описание:** Удаление пользователей из группы (только владелец).
**Параметры URL:**
- `id` - ID группы

**Тело запроса:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "string"
    }
  ]
}
```

**Возможные коды ответов:**
- 200 OK — пользователи удалены
- 400 Bad Request — некорректные данные
- 401 Unauthorized — пользователь не авторизован
- 403 Forbidden — пользователь не владелец группы
- 409 Conflict — пользователь не в группе
- 500 Internal Server Error — внутренняя ошибка сервера

---

### Список чатов

#### GET /api/v1/chats
**Описание:** Получение списка чатов пользователя.
**Параметры запроса (query):**
- `limit` (optional) - количество чатов (по умолчанию: без ограничения)
- `offset` (optional) - смещение (по умолчанию: 0)

**Возможные коды ответов:**
- 200 OK — список получен
- 400 Bad Request — некорректные параметры
- 401 Unauthorized — пользователь не авторизован
- 404 Not Found — чаты не найдены
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "rooms": [
    {
      "id": 1,
      "name": "string",
      "type": "private|group",
      "last_message": {
        "id": 1,
        "chat_id": 1,
        "type": "message",
        "text": "string",
        "author": {
          "id": 1,
          "name": "string"
        },
        "time": "2023-01-01T00:00:00Z"
      }
    }
  ]
}
```

---

## WebSocket API

### Подключение к приватному чату
**URL:** `ws://host/api/v1/chat/{id}?limit={limit}&offset={offset}`
**Параметры URL:**
- `id` - ID чата
**Параметры запроса:**
- `limit` (optional) - количество сообщений истории (по умолчанию: 100)
- `offset` (optional) - смещение истории (по умолчанию: 0)

**Заголовки:** 
- `Authorization: Bearer {access_token}`

### Подключение к групповому чату
**URL:** `ws://host/api/v1/group/{id}?limit={limit}&offset={offset}`
**Параметры URL:**
- `id` - ID группы
**Параметры запроса:**
- `limit` (optional) - количество сообщений истории (по умолчанию: 100)
- `offset` (optional) - смещение истории (по умолчанию: 0)

**Заголовки:** 
- `Authorization: Bearer {access_token}`

### Главное подключение (получение уведомлений)
**URL:** `ws://host/api/v1/main`
**Описание:** Подключение для получения последних сообщений из всех чатов пользователя.

**Заголовки:** 
- `Authorization: Bearer {access_token}`

### Типы сообщений WebSocket

#### Отправка сообщения
```json
{
  "type": "message",
  "text": "string",
  "author": {
    "id": 1,
    "name": "string"
  }
}
```

#### Индикатор печати
```json
{
  "type": "typing",
  "author": {
    "id": 1,
    "name": "string"
  }
}
```

#### Остановка индикатора печати
```json
{
  "type": "stop_typing",
  "author": {
    "id": 1,
    "name": "string"
  }
}
```

#### Получение истории (от сервера)
```json
{
  "type": "history",
  "messages": [
    {
      "id": 1,
      "chat_id": 1,
      "type": "message",
      "text": "string",
      "author": {
        "id": 1,
        "name": "string"
      },
      "time": "2023-01-01T00:00:00Z"
    }
  ],
  "error": "string" // если нет сообщений
}
```

#### Последнее сообщение (от сервера)
```json
{
  "id": 1,
  "chat_id": 1,
  "type": "last_message",
  "text": "string",
  "author": {
    "id": 1,
    "name": "string"
  },
  "time": "2023-01-01T00:00:00Z"
}
```

---

## Модели данных

### Message
```json
{
  "id": 1,
  "chat_id": 1,
  "type": "message|typing|stop_typing|last_message",
  "text": "string",
  "author": {
    "id": 1,
    "name": "string"
  },
  "time": "2023-01-01T00:00:00Z"
}
```

### ChatRoom
```json
{
  "id": 1,
  "name": "string",
  "type": "private|group",
  "last_message": {
    // Message object or null
  }
}
```

### User
```json
{
  "id": 1,
  "name": "string",
  "username": "string",
  "description": "string"
}
```

### Group
```json
{
  "id": 1,
  "name": "string",
  "description": "string"
}
```

---

## Примечания

- Все WebSocket подключения требуют JWT аутентификации через заголовок Authorization
- При подключении к чату/группе автоматически отправляется история сообщений
- Сообщения передаются в real-time через Redis Pub/Sub
- Поддержка файлов через multipart/form-data для создания пользователей и групп
- Автоматическое создание папок для хранения файлов чатов и групп
- Валидация владельца группы для операций управления участниками
- Graceful закрытие WebSocket соединений с очисткой ресурсов
