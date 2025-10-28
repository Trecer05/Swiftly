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

#### PUT /api/v1/user
**Описание:** Создание профиля пользователя в чат-системе.
**Content-Type:** `application/json`

**JSON данные:**
```json
{
  "name": "string", 
  "description": "string"
}
```
В json один из двух параметров либо оба

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
  "owner_id": 1,
  "users": [
    {
      "id": 1,
      "name": "string"
    }
  ]
}
```

**Примечание:** Поле `owner_id` автоматически устанавливается из JWT токена текущего пользователя.

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

#### DELETE /api/v1/chat/{id}
**Описание:** Удаление чата.
**Параметры URL:**
- `id` - ID группы

**Возможные коды ответов:**
- 200 OK — чат удален
- 400 Bad Request — некорректный ID
- 401 Unauthorized — пользователь не авторизован
- 500 Internal Server Error — внутренняя ошибка сервера

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
  "description": "string",
  "owner_id": 1
}
```

**Примечание:** Поле `owner_id` показывает ID владельца группы.

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

#### DELETE /api/v1/group/{id}/leave
**Описание:** Выход пользователя из группы (самостоятельно).
**Параметры URL:**
- `id` - ID группы

**Возможные коды ответов:**
- 200 OK — пользователь вышел из группы
- 400 Bad Request — некорректный ID группы
- 401 Unauthorized — пользователь не авторизован
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
{
  "status": "ok"
}
```

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

## Файлы и медиа

### Загрузка файлов в приватные чаты

#### POST /api/v1/chat/{id}/img
**Описание:** Загрузка изображений в приватный чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID чата

**Параметры:**
- `photos` (files) - массив изображений

**Поддерживаемые форматы:** jpg, png, gif, webp, svg, bmp

**MIME-типы:** image/jpeg, image/png, image/gif, image/webp, image/svg+xml, image/bmp

**Возможные коды ответов:**
- 200 OK — файлы загружены
- 400 Bad Request — некорректные данные
- 401 Unauthorized — пользователь не авторизован
- 500 Internal Server Error — внутренняя ошибка сервера

**Ответ (успех):**
```json
[
  "/chat/1/img/1234567890_photo1.jpg",
  "/chat/1/img/1234567891_photo2.png"
]
```

#### POST /api/v1/chat/{id}/video
**Описание:** Загрузка видео в приватный чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID чата

**Параметры:**
- `videos` (files) - массив видеофайлов

**Поддерживаемые форматы:** mp4, mpeg, quicktime, webm, avi, wmv

**MIME-типы:** video/mp4, video/mpeg, video/quicktime, video/webm, video/x-msvideo, video/x-ms-wmv

#### POST /api/v1/chat/{id}/audio
**Описание:** Загрузка аудио в приватный чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID чата

**Параметры:**
- `photos` (files) - массив аудиофайлов

**Поддерживаемые форматы:** mp3, wav, ogg, webm, aac

**MIME-типы:** audio/mpeg, audio/wav, audio/ogg, audio/webm, audio/aac

#### POST /api/v1/chat/{id}/file
**Описание:** Загрузка документов в приватный чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID чата

**Параметры:**
- `photos` (files) - массив документов

**Поддерживаемые форматы:** pdf, doc, docx, xls, xlsx, txt, csv, zip, rar

**MIME-типы:** application/pdf, application/msword, application/vnd.openxmlformats-officedocument.wordprocessingml.document, application/vnd.ms-excel, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, text/plain, text/csv, application/zip, application/x-rar-compressed

#### POST /api/v1/chat/{id}/imgvid
**Описание:** Смешанная загрузка изображений и видео в приватный чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID чата

**Параметры:**
- `files` (files) - массив файлов (автоматически сортируются по типу)

### Загрузка файлов в групповые чаты

#### POST /api/v1/group/{id}/img
**Описание:** Загрузка изображений в групповой чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID группы

**Параметры:**
- `photos` (files) - массив изображений

#### POST /api/v1/group/{id}/video
**Описание:** Загрузка видео в групповой чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID группы

**Параметры:**
- `videos` (files) - массив видеофайлов

#### POST /api/v1/group/{id}/audio
**Описание:** Загрузка аудио в групповой чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID группы

**Параметры:**
- `photos` (files) - массив аудиофайлов

#### POST /api/v1/group/{id}/file
**Описание:** Загрузка документов в групповой чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID группы

**Параметры:**
- `photos` (files) - массив документов

#### POST /api/v1/group/{id}/imgvid
**Описание:** Смешанная загрузка изображений и видео в групповой чат.
**Content-Type:** `multipart/form-data`
**Параметры URL:**
- `id` - ID группы

**Параметры:**
- `files` (files) - массив файлов (автоматически сортируются по типу)

### Получение файлов

#### GET /api/v1/chat/{id}/files
**Описание:** Получение списка всех файлов приватного чата.
**Параметры URL:**
- `id` - ID чата

**Возможные коды ответов:**
- 200 OK — список получен
- 400 Bad Request — некорректный ID
- 401 Unauthorized — пользователь не авторизован
- 500 Internal Server Error — внутренняя ошибка сервера

#### GET /api/v1/group/{id}/files
**Описание:** Получение списка всех файлов группового чата.
**Параметры URL:**
- `id` - ID группы

#### GET /api/v1/chat/{id}/img/{url}
**Описание:** Скачивание изображения из приватного чата.
**Параметры URL:**
- `id` - ID чата
- `url` - имя файла

#### GET /api/v1/chat/{id}/video/{url}
**Описание:** Скачивание видео из приватного чата.
**Параметры URL:**
- `id` - ID чата
- `url` - имя файла

#### GET /api/v1/chat/{id}/audio/{url}
**Описание:** Скачивание аудио из приватного чата.
**Параметры URL:**
- `id` - ID чата
- `url` - имя файла

#### GET /api/v1/chat/{id}/file/{url}
**Описание:** Скачивание документа из приватного чата.
**Параметры URL:**
- `id` - ID чата
- `url` - имя файла

#### GET /api/v1/group/{id}/img/{url}
**Описание:** Скачивание изображения из группового чата.
**Параметры URL:**
- `id` - ID группы
- `url` - имя файла

#### GET /api/v1/group/{id}/video/{url}
**Описание:** Скачивание видео из группового чата.
**Параметры URL:**
- `id` - ID группы
- `url` - имя файла

#### GET /api/v1/group/{id}/audio/{url}
**Описание:** Скачивание аудио из группового чата.
**Параметры URL:**
- `id` - ID группы
- `url` - имя файла

#### GET /api/v1/group/{id}/file/{url}
**Описание:** Скачивание документа из группового чата.
**Параметры URL:**
- `id` - ID группы
- `url` - имя файла

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

#### Сообщение с файлами
```json
{
  "type": "with_files",
  "text": "string",
  "author": {
    "id": 1,
    "name": "string"
  },
  "file_urls": [
    "/chat/1/img/1234567890_photo.jpg",
    "/chat/1/video/1234567891_video.mp4"
  ],
  "file_name": "photo.jpg",
  "file_mime": "image/jpeg",
  "file_type": "photo",
  "file_size": 1024000
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

#### Статус прочтения сообщения
```json
{
  "type": "read",
  "id": 123,
  "read": true,
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
      "read": false,
      "author": {
        "id": 1,
        "name": "string"
      },
      "time": "2023-01-01T00:00:00Z",
      "file_url": "string",
      "file_urls": ["string"],
      "file_name": "string",
      "file_mime": "string",
      "file_type": "photo",
      "file_size": 1024000
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
  "read": false,
  "author": {
    "id": 1,
    "name": "string"
  },
  "time": "2023-01-01T00:00:00Z",
  "file_url": "string",
  "file_urls": ["string"],
  "file_name": "string",
  "file_mime": "string",
  "file_type": "photo",
  "file_size": 1024000
}
```

---

## Модели данных

### Message
```json
{
  "id": 1,
  "chat_id": 1,
  "type": "message|with_files|typing|stop_typing|read|last_message",
  "text": "string",
  "read": false,
  "author": {
    "id": 1,
    "name": "string"
  },
  "time": "2023-01-01T00:00:00Z",
  "file_url": "string",
  "file_urls": ["string"],
  "file_name": "string",
  "file_mime": "string",
  "file_type": "photo|video|audio|file|other",
  "file_size": 1024000
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
  "description": "string",
  "owner_id": 1
}
```

---

## Примечания

### Аутентификация и безопасность
- Все эндпойнты требуют JWT аутентификации через Bearer токен в заголовке Authorization
- WebSocket подключения также требуют JWT аутентификации
- Валидация владельца группы для операций управления участниками

### Real-time общение
- При подключении к чату/группе автоматически отправляется история сообщений
- Сообщения передаются в real-time через Redis Pub/Sub
- Поддержка индикаторов печати (typing/stop_typing)
- Статусы прочтения сообщений (read/unread)
- Graceful закрытие WebSocket соединений с очисткой ресурсов

### Файлы и медиа
- Поддержка всех типов файлов: изображения, видео, аудио, документы
- Автоматическое определение MIME-типов файлов
- Организованное хранение файлов по папкам (photos, videos, audios, files)
- Смешанная загрузка изображений и видео в одном запросе
- Максимальный размер файла: 32MB для загрузки
- Автоматическое создание папок для хранения файлов чатов и групп

### Управление группами
- Только владелец группы может добавлять/удалять участников
- Пользователи могут самостоятельно выходить из группы
- Поддержка загрузки фото группы при создании
- Валидация прав доступа для всех операций с группами

### База данных
- Использование PostgreSQL с миграциями
- Поддержка транзакций для целостности данных
- Индексы для оптимизации запросов
- Cascade удаление связанных записей

---

## Звонки (WebRTC)

### Подключение к звонку в приватном чате
**URL:** `ws://host/api/v1/chat/{id}/call`
**Параметры URL:**
- `id` - ID чата

**Заголовки:** 
- `Authorization: Bearer {access_token}`

**Описание:** WebSocket соединение для WebRTC звонков в приватном чате. Поддерживает аудио и видео звонки.

### Подключение к звонку в групповом чате
**URL:** `ws://host/api/v1/group/{id}/call`
**Параметры URL:**
- `id` - ID группы

**Заголовки:** 
- `Authorization: Bearer {access_token}`

**Описание:** WebSocket соединение для WebRTC звонков в групповом чате. Поддерживает аудио и видео звонки.

### Типы сообщений WebRTC

#### Присоединение к комнате
```json
{
  "type": "join",
  "roomId": 1,
  "sessionId": "unique-session-id"
}
```

#### WebRTC Offer
```json
{
  "type": "offer",
  "payload": {
    "type": "offer",
    "sdp": "v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\n..."
  },
  "roomId": 1,
  "sessionId": "unique-session-id"
}
```

#### WebRTC Answer
```json
{
  "type": "answer",
  "payload": {
    "type": "answer",
    "sdp": "v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\n..."
  },
  "roomId": 1,
  "sessionId": "unique-session-id"
}
```

#### ICE Candidate
```json
{
  "type": "ice",
  "payload": {
    "candidate": "candidate:1 1 UDP 2113667326 192.168.1.100 54400 typ host",
    "sdpMLineIndex": 0,
    "sdpMid": "0"
  },
  "roomId": 1,
  "sessionId": "unique-session-id"
}
```

#### Покидание комнаты
```json
{
  "type": "leave",
  "roomId": 1,
  "sessionId": "unique-session-id"
}
```

## Дополнения к маршрутам общим
### Пользователи (дополнения)

#### POST /api/v1/user/avatar

**Описание:** Загрузка аватара для текущего пользователя.  

**Content-Type:** `multipart/form-data`  

**Параметры:**  

- `photo` (file) - фото аватара  

**Возможные коды ответов:**  

- 200 OK — аватар загружен  

- 400 Bad Request — некорректные данные запроса  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "url": "string"  // Имя файла аватара  

}  

```

#### DELETE /api/v1/user/avatar/{url}

**Описание:** Удаление аватара текущего пользователя по URL.  

**Параметры URL:**  

- `url` - имя файла аватара  

**Возможные коды ответов:**  

- 200 OK — аватар удален  

- 400 Bad Request — некорректный URL  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "status": "ok"  

}  

```

#### GET /api/v1/user/avatar/urls

**Описание:** Получение списка URL аватаров текущего пользователя.  

**Возможные коды ответов:**  

- 200 OK — список получен  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

[  

  "string",  // Имя файла 1  

  "string"   // Имя файла 2  

]  

```

#### GET /api/v1/user/avatar/{url}

**Описание:** Получение аватара текущего пользователя по URL.  

**Параметры URL:**  

- `url` - имя файла аватара  

**Возможные коды ответов:**  

- 200 OK — аватар получен (бинарные данные)  

- 400 Bad Request — некорректный URL  

- 401 Unauthorized — пользователь не авторизован  

- 404 Not Found — аватар не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

#### GET /api/v1/user/{id}/avatar/urls

**Описание:** Получение списка URL аватаров другого пользователя.  

**Параметры URL:**  

- `id` - ID пользователя  

**Возможные коды ответов:**  

- 200 OK — список получен  

- 400 Bad Request — некорректный ID  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

[  

  "string",  // Имя файла 1  

  "string"   // Имя файла 2  

]  

```

#### GET /api/v1/user/{id}/avatar/{url}

**Описание:** Получение аватара другого пользователя по URL.  

**Параметры URL:**  

- `id` - ID пользователя  

- `url` - имя файла аватара  

**Возможные коды ответов:**  

- 200 OK — аватар получен (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 401 Unauthorized — пользователь не авторизован  

- 404 Not Found — аватар не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

---

### Приватные чаты (дополнения)

#### DELETE /api/v1/chat/{id}

**Описание:** Удаление приватного чата.  

**Параметры URL:**  

- `id` - ID чата  

**Возможные коды ответов:**  

- 200 OK — чат удален  

- 400 Bad Request — некорректный ID  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "status": "ok"  

}  

```

#### GET /api/v1/chat/{id}/messages

**Описание:** Получение сообщений приватного чата (с пагинацией).  

**Параметры запроса (query):**  

- `limit` (optional) - лимит сообщений  

- `offset` (optional) - смещение  

**Возможные коды ответов:**  

- 200 OK — сообщения получены  

- 400 Bad Request — некорректные параметры  

- 401 Unauthorized — пользователь не авторизован  

- 404 Not Found — сообщений нет  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "messages": [  

    {  

      "id": 1,  

      "chat_id": 1,  

      "user_id": 1,  

      "content": "string",  

      "type": "text",  

      "timestamp": "timestamp"  

    }  

  ]  

}  

```

---

### Групповые чаты (новый раздел, если не в truncated)

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

  "owner_id": 0  // Автоматически подставляется  

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

    "name": "string",  

    "url": "string"  // Имя фото  

  }  

}  

```

#### PUT /api/v1/group/{id}

**Описание:** Обновление группового чата.  

**Content-Type:** `application/json`  

**Параметры URL:**  

- `id` - ID группы  

**JSON данные:**  

```json

{  

  "name": "string"  // Или другие поля  

}  

```

**Возможные коды ответов:**  

- 200 OK — группа обновлена  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 404 Not Found — группа не найдена  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "status": "ok"  

}  

```

#### DELETE /api/v1/group/{id}

**Описание:** Удаление группового чата (только владельцем).  

**Параметры URL:**  

- `id` - ID группы  

**Возможные коды ответов:**  

- 200 OK — группа удалена  

- 400 Bad Request — некорректный ID  

- 401 Unauthorized — пользователь не авторизован  

- 403 Forbidden — не владелец  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "status": "ok"  

}  

```

#### POST /api/v1/group/{id}/exit

**Описание:** Выход из группового чата.  

**Параметры URL:**  

- `id` - ID группы  

**Возможные коды ответов:**  

- 200 OK — выход выполнен  

- 400 Bad Request — некорректный ID  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "status": "ok"  

}  

```

#### GET /api/v1/group/{id}/info

**Описание:** Получение информации о групповом чате.  

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

  "description": "string",  

  "owner_id": 1  

}  

```

#### GET /api/v1/group/{id}/users

**Описание:** Получение списка пользователей в группе.  

**Параметры URL:**  

- `id` - ID группы  

**Возможные коды ответов:**  

- 200 OK — список получен  

- 400 Bad Request — некорректный ID  

- 404 Not Found — пользователей нет  

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

#### GET /api/v1/group/{id}/messages

**Описание:** Получение сообщений группового чата (с пагинацией).  

**Параметры запроса (query):**  

- `limit` (optional) - лимит сообщений  

- `offset` (optional) - смещение  

**Возможные коды ответов:**  

- 200 OK — сообщения получены  

- 400 Bad Request — некорректные параметры  

- 401 Unauthorized — пользователь не авторизован  

- 404 Not Found — сообщений нет  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "messages": [  

    {  

      "id": 1,  

      "group_id": 1,  

      "user_id": 1,  

      "content": "string",  

      "type": "text",  

      "timestamp": "timestamp"  

    }  

  ]  

}  

```

#### GET /api/v1/chats

**Описание:** Получение списка чатов/групп текущего пользователя (с пагинацией).  

**Параметры запроса (query):**  

- `limit` (optional) - лимит комнат  

- `offset` (optional) - смещение  

**Возможные коды ответов:**  

- 200 OK — список получен  

- 400 Bad Request — некорректные параметры  

- 401 Unauthorized — пользователь не авторизован  

- 404 Not Found — комнат нет  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "rooms": [  

    {  

      "id": 1,  

      "type": "private/group",  

      "name": "string",  

      "last_message": {  

        "content": "string"  

      }  

    }  

  ]  

}  

```

---

### Файлы и медиа (расширение существующего раздела)

#### POST /api/v1/chat/{id}/img или /api/v1/group/{id}/img

**Описание:** Загрузка изображений в чат/группу.  

**Content-Type:** `multipart/form-data`  

**Параметры URL:**  

- `id` - ID чата/группы  

**Параметры:**  

- `photos` (files) - изображения  

**Возможные коды ответов:**  

- 200 OK — изображения загружены  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

[  

  "/chat/1/img/filename1",  

  "/chat/1/img/filename2"  

]  

```

#### POST /api/v1/chat/{id}/video или /api/v1/group/{id}/video

**Описание:** Загрузка видео в чат/группу.  

**Content-Type:** `multipart/form-data`  

**Параметры URL:**  

- `id` - ID чата/группы  

**Параметры:**  

- `videos` (files) - видео  

**Возможные коды ответов:**  

- 200 OK — видео загружены  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

[  

  "/chat/1/video/filename1"  

]  

```

#### POST /api/v1/chat/{id}/audio или /api/v1/group/{id}/audio

**Описание:** Загрузка аудио в чат/группу.  

**Content-Type:** `multipart/form-data`  

**Параметры URL:**  

- `id` - ID чата/группы  

**Параметры:**  

- `audio` (files) - аудиофайлы  

**Возможные коды ответов:**  

- 200 OK — аудио загружены  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

[  

  "/chat/1/audio/filename1"  

]  

```

#### POST /api/v1/chat/{id}/file или /api/v1/group/{id}/file

**Описание:** Загрузка документов в чат/группу.  

**Content-Type:** `multipart/form-data`  

**Параметры URL:**  

- `id` - ID чата/группы  

**Параметры:**  

- `files` (files) - документы  

**Возможные коды ответов:**  

- 200 OK — документы загружены  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

[  

  "/chat/1/file/filename1"  

]  

```

#### GET /api/v1/chat/{id}/img/{url} или /api/v1/group/{id}/img/{url}

**Описание:** Получение изображения по URL.  

**Параметры URL:**  

- `id` - ID чата/группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — изображение получено (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 404 Not Found — файл не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

#### GET /api/v1/chat/{id}/video/{url} или /api/v1/group/{id}/video/{url}

**Описание:** Получение видео по URL.  

**Параметры URL:**  

- `id` - ID чата/группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — видео получено (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 404 Not Found — файл не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

#### GET /api/v1/chat/{id}/audio/{url} или /api/v1/group/{id}/audio/{url}

**Описание:** Получение аудио по URL.  

**Параметры URL:**  

- `id` - ID чата/группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — аудио получено (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 404 Not Found — файл не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

#### GET /api/v1/chat/{id}/file/{url} или /api/v1/group/{id}/file/{url}

**Описание:** Получение документа по URL.  

**Параметры URL:**  

- `id` - ID чата/группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — документ получен (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 404 Not Found — файл не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

#### GET /api/v1/chat/{id}/files или /api/v1/group/{id}/files

**Описание:** Получение списка медиафайлов (photos/videos) в чате/группе.  

**Параметры URL:**  

- `id` - ID чата/группы  

**Возможные коды ответов:**  

- 200 OK — список получен  

- 400 Bad Request — некорректный ID  

- 404 Not Found — файлов нет  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

[  

  "filename1",  

  "filename2"  

]  

```

#### POST /api/v1/group/{id}/avatar

**Описание:** Загрузка аватара группы.  

**Content-Type:** `multipart/form-data`  

**Параметры URL:**  

- `id` - ID группы  

**Параметры:**  

- `photo` (file) - фото аватара  

**Возможные коды ответов:**  

- 200 OK — аватар загружен  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "url": "string"  

}  

```

#### GET /api/v1/group/{id}/avatar/url

**Описание:** Получение URL последнего аватара группы.  

**Параметры URL:**  

- `id` - ID группы  

**Возможные коды ответов:**  

- 200 OK — URL получен  

- 400 Bad Request — некорректный ID  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "url": "string"  

}  

```

#### GET /api/v1/group/{id}/avatar/{url}

**Описание:** Получение аватара группы по URL.  

**Параметры URL:**  

- `id` - ID группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — аватар получен (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 404 Not Found — аватар не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

#### DELETE /api/v1/group/{id}/avatar/{url}

**Описание:** Удаление аватара группы по URL.  

**Параметры URL:**  

- `id` - ID группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — аватар удален  

- 400 Bad Request — некорректные параметры  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "status": "ok"  

}  

```

#### POST /api/v1/chat/{id}/message/audio или /api/v1/group/{id}/message/audio

**Описание:** Загрузка аудиосообщения в чат/группу.  

**Content-Type:** `multipart/form-data`  

**Параметры URL:**  

- `id` - ID чата/группы  

**Параметры:**  

- `audio` (file) - аудиофайл  

**Возможные коды ответов:**  

- 200 OK — аудио загружено  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "url": "string"  

}  

```

#### GET /api/v1/chat/{id}/message/audio/{url} или /api/v1/group/{id}/message/audio/{url}

**Описание:** Получение аудиосообщения по URL.  

**Параметры URL:**  

- `id` - ID чата/группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — аудио получено (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 404 Not Found — файл не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

#### POST /api/v1/chat/{id}/message/video или /api/v1/group/{id}/message/video

**Описание:** Загрузка видеосообщения в чат/группу.  

**Content-Type:** `multipart/form-data`  

**Параметры URL:**  

- `id` - ID чата/группы  

**Параметры:**  

- `video` (file) - видеофайл  

**Возможные коды ответов:**  

- 200 OK — видео загружено  

- 400 Bad Request — некорректные данные  

- 401 Unauthorized — пользователь не авторизован  

- 500 Internal Server Error — внутренняя ошибка сервера  

**Ответ (успех):**  

```json

{  

  "url": "string"  

}  

```

#### GET /api/v1/chat/{id}/message/video/{url} или /api/v1/group/{id}/message/video/{url}

**Описание:** Получение видеосообщения по URL.  

**Параметры URL:**  

- `id` - ID чата/группы  

- `url` - имя файла  

**Возможные коды ответов:**  

- 200 OK — видео получено (бинарные данные)  

- 400 Bad Request — некорректные параметры  

- 404 Not Found — файл не найден  

- 500 Internal Server Error — внутренняя ошибка сервера  

---

### Real-time общение (дополнения к существующему разделу)

#### WS /api/v1/main

**Описание:** WebSocket для основного соединения (получение уведомлений, статусов и последних сообщений по всем чатам пользователя).  

**Заголовки:**  

- `Authorization: Bearer {access_token}`  

**Описание:** Автоматически подписывается на все комнаты пользователя, отправляет последние сообщения, статусы и уведомления.  

**Типы сообщений (входящие/исходящие):** Как в существующем разделе (messages, status, notifications).  

#### WS /api/v1/chat/{id}

**Описание:** WebSocket для приватного чата (real-time сообщения, статусы, история).  

**Параметры URL:**  

- `id` - ID чата  

**Заголовки:**  

- `Authorization: Bearer {access_token}`  

**Параметры запроса (query):**  

- `limit` (optional) - лимит истории сообщений  

- `offset` (optional) - смещение  

**Описание:** Отправляет историю сообщений при подключении, real-time обновления.  

#### WS /api/v1/group/{id}

**Описание:** WebSocket для группового чата (real-time сообщения, статусы, история).  

**Параметры URL:**  

- `id` - ID группы  

**Заголовки:**  

- `Authorization: Bearer {access_token}`  

**Параметры запроса (query):**  

- `limit` (optional) - лимит истории сообщений  

- `offset` (optional) - смещение  

**Описание:** Отправляет историю сообщений при подключении, real-time обновления.  

### Real-time общение

#### WS /api/v1/main

**Описание:** WebSocket для основного соединения (получение уведомлений, статусов и последних сообщений по всем чатам пользователя).  

**Заголовки:**  

- `Authorization: Bearer {access_token}`  

**Описание:** Автоматически подписывается на все комнаты пользователя, отправляет последние сообщения, статусы и уведомления.  

#### WS /api/v1/chat/{id}

**Описание:** WebSocket для приватного чата (real-time сообщения, статусы, история).  

**Параметры URL:**  

- `id` - ID чата  

**Заголовки:**  

- `Authorization: Bearer {access_token}`  

**Параметры запроса (query):**  

- `limit` (optional) - лимит истории сообщений  

- `offset` (optional) - смещение  

**Описание:** Отправляет историю сообщений при подключении, real-time обновления.  

#### WS /api/v1/group/{id}

**Описание:** WebSocket для группового чата (real-time сообщения, статусы, история).  

**Параметры URL:**  

- `id` - ID группы  

**Заголовки:**  

- `Authorization: Bearer {access_token}`  

**Параметры запроса (query):**  

- `limit` (optional) - лимит истории сообщений  

- `offset` (optional) - смещение  

**Описание:** Отправляет историю сообщений при подключении, real-time обновления.  

#### Типы сообщений WebSocket

Все сообщения оборачиваются в конверт (Envelope) с типом (message, notification, status). Ниже приведены примеры для каждого типа.

##### Message (EnvelopeType: "message")

Сообщения для чатов, включая текст, файлы, аудио/видео-сообщения, действия (typing, read и т.д.).

- **Typing (входящее/исходящее):** Пользователь печатает.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "typing",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **StopTyping (входящее/исходящее):** Пользователь прекратил печатать.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "stop_typing",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **Default (входящее/исходящее):** Обычное текстовое сообщение.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "message",  

    "text": "Hello, world!",  

    "read": true,  

    "status": {  

      "type": "string",  

      "user_id": 1,  

      "online": true  

    },  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z",  

    "edited": false  

  }  

}  

```

- **WithFiles (входящее/исходящее):** Сообщение с файлами.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "with_files",  

    "file_urls": ["/chat/1/file/filename1", "/chat/1/file/filename2"],  

    "file_name": "filename1",  

    "file_mime": "application/pdf",  

    "file_type": "file",  

    "file_size": 1024,  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **AudioMessage (входящее/исходящее):** Аудиосообщение.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "audio_message",  

    "audio_message_url": "/chat/1/message/audio/filename",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **VideoMessage (входящее/исходящее):** Видеосообщение.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "video_message",  

    "video_message_url": "/chat/1/message/video/filename",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **LastMessage (исходящее):** Последнее сообщение в чате (при подключении).  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "last_message",  

    "text": "Last message text",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **Read (входящее/исходящее):** Сообщение прочитано.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "read",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **Delete (входящее/исходящее):** Сообщение удалено.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "delete",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **Update (входящее/исходящее):** Сообщение обновлено.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "update",  

    "text": "Updated text",  

    "edited": true,  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

- **Call (входящее/исходящее):** Инициация звонка.  

```json

{  

  "type": "message",  

  "payload": {  

    "id": 1,  

    "chat_id": 1,  

    "type": "call",  

    "author": {  

      "id": 1,  

      "name": "string"  

    },  

    "time": "2023-10-01T12:00:00Z"  

  }  

}  

```

##### Status (EnvelopeType: "status")

Статус пользователя (онлайн/оффлайн).

```json

{  

  "type": "status",  

  "payload": {  

    "type": "status",  

    "user_id": 1,  

    "online": true  

  }  

}  

```

##### Notifications (EnvelopeType: "notification")

Уведомления о сообщениях или звонках.

- **MessageNotification:** Уведомление о новом сообщении.  

```json

{  

  "type": "notification",  

  "payload": {  

    "type": "message",  

    "message": {  

      "id": 1,  

      "chat_id": 1,  

      "type": "message",  

      "text": "New message"  

    }  

  }  

}  

```

- **CallNotification:** Уведомление о звонке.  

```json

{  

  "type": "notification",  

  "payload": {  

    "type": "call",  

    "message": {  

      "id": 1,  

      "chat_id": 1,  

      "type": "call",  

      "author": {  

        "id": 1,  

        "name": "string"  

      }  

    }  

  }  

}  

```

**Примечание:** Подробная документация по звонкам доступна в [CALLS_API_DOC.md](../calls/CALLS_API_DOC.md)