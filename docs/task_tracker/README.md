# Task Tracker API — Полная документация (Markdown)

**Базовый URL:** `/api/v1`
**Аутентификация:** все эндпоинты (кроме `/health`) требуют заголовка

```
Authorization: Bearer {access_token}
```

**Rate limit:** 100 запросов / минута (см. middleware)

---

# Содержание

1. [Общие модели JSON](#общие-модели-json)
2. [HTTP API](#http-api)

   * [Health check](#get-apiv1health)
   * [Dashboard info](#get-apiv1teamteam_iddashboardinfo)
   * [Create task](#post-apiv1teamteam_idcolumncolumn_idtask)
   * [Get task](#get-apiv1teamteam_idtasktask_id)
   * [Delete task](#delete-apiv1teamteam_idcolumncolumn_idtasktask_id)
   * [Update task title](#put-apiv1teamteam_idtasktask_idtitle)
   * [Update task description](#put-apiv1teamteam_idtasktask_iddescription)
   * [Update task deadline](#put-apiv1teamteam_idtasktask_iddeadline)
   * [Update task developer](#put-apiv1teamteam_idtasktask_iddeveloper)
   * [Update task status](#put-apiv1teamteam_idtasktask_idstatus)
   * [Move task (column/position)](#put-apiv1teamteam_idtasktask_idmove)
3. [WebSocket API (Dashboard)](#websocket-api-dashboard)

   * [Подключение](#подключение)
   * [Формат Envelope (WS)](#формат-envelope-ws)
   * [Типы сообщений и примеры данных](#типы-сообщений-и-примеры-данных)
4. [Ошибки и коды ответа — единообразие](#ошибки-и-коды-ответа---единообразие)
5. [Примечания для реализации на фронтенде / интеграции](#примечания-для-реализации-на-фронтенде--интеграции)

---

# Общие модели JSON

Ниже — схемы используемых объектов (примерный формат).

## Priority

```json
{
  "type": "low | medium | high",
  "title": "Low | Medium | High",
  "hex_color": "#RRGGBB"
}
```

## DashboardTask

```json
{
  "id": 123,
  "title": "Task title",
  "description": "Optional description or null",
  "status": "todo | in_progress | completed",
  "priority": {
    "type": "medium",
    "title": "Medium",
    "hex_color": "#FFFF00"
  },
  "label": "frontend",
  "deadline": "2025-12-01T18:00:00Z" // ISO 8601, optional or null
}
```

## Column

```json
{
  "id": 3,
  "title": "Backlog",
  "position": 0,
  "hex_color": "#D3D3D3",
  "tasks": [ /* DashboardTask */ ]
}
```

## DashboardInfo

```json
{
  "team_id": 5,
  "columns": [ /* Column */ ]
}
```

## Task (полная)

```json
{
  "id": 123,
  "author_id": 11,
  "developer_id": 21,
  "label": "backend",
  "title": "Implement API",
  "description": "Detailed description or null",
  "status": "todo | in_progress | completed",
  "priority": {
    "type": "medium",
    "title": "Medium",
    "hex_color": "#FFFF00"
  },
  "created_at": "2025-11-30T10:00:00Z",
  "updated_at": "2025-11-30T12:00:00Z",
  "deadline": "2025-12-10T18:00:00Z",
  "completed_at": null
}
```

---

# HTTP API

> Все ответы успешные (200/201) содержат JSON и заголовок `Content-Type: application/json`.

---

### GET /api/v1/health

**Описание:** Простейшая проверка работоспособности сервиса.
**Метод:** `GET`
**Аутентификация:** **не требуется**

**Ответы:**

200 OK

```json
{
  "status": "ok"
}
```

---

### GET /api/v1/team/{team_id}/dashboard/info

**Описание:** Возвращает информацию для рендера доски: список колонок с задачами (в порядке `position`), включая приоритеты задач.
**Метод:** `GET`
**Аутентификация:** требуется

**Параметры URL:**

* `team_id` (int) — ID команды

**Успешный ответ (200):**

```json
{
  "dashboard_info": {
    "team_id": 5,
    "columns": [
      {
        "id": 1,
        "title": "Backlog",
        "position": 0,
        "hex_color": "#D3D3D3",
        "tasks": [
          {
            "id": 101,
            "title": "Design DB schema",
            "description": "Draw ER and discuss with team",
            "status": "todo",
            "priority": {
              "type": "medium",
              "title": "Medium",
              "hex_color": "#FFFF00"
            },
            "label": "architecture",
            "deadline": "2025-12-05T12:00:00Z"
          },
          {
            "id": 102,
            "title": "Research caching",
            "description": null,
            "status": "todo",
            "priority": {
              "type": "low",
              "title": "Low",
              "hex_color": "#00FF00"
            },
            "label": "research",
            "deadline": null
          }
        ]
      },
      {
        "id": 2,
        "title": "In Progress",
        "position": 1,
        "hex_color": "#E8F5FF",
        "tasks": [ /* ... */ ]
      }
    ]
  }
}
```

**Коды ошибок:**

* 400 Bad Request — `team_id` некорректен
* 401 Unauthorized — нет токена
* 500 Internal Server Error — ошибка сервера

---

### POST /api/v1/team/{team_id}/column/{column_id}/task

**Описание:** Создаёт новую задачу в колонке. Если указан `position_in_column`, механизм БД сдвинет существующие задачи (см. триггеры). Если `position_in_column` не указан, на сервере/DB ставится `max+1`.
**Метод:** `POST`
**Content-Type:** `application/json`
**Аутентификация:** требуется

**Параметры URL:**

* `team_id` — ID команды
* `column_id` — ID колонки

**JSON запрос:**

```json
{
  "title": "Implement login",
  "description": "OAuth2 + refresh tokens",
  "position_in_column": 0,            // optional, integer
  "deadline": "2025-12-01T18:00:00Z", // optional, ISO8601
  "label": "backend",                 // optional
  "priority": "high"                  // optional: "low"|"medium"|"high"
}
```

> На сервере поле `creator_id` берётся из контекста (token -> user id).
> Поле `priority` можно принимать и применять; если nil — используется дефолт (например `medium`) согласно бизнес-логике.

**Успешный ответ (201 Created):**

```json
{
  "id": 201
}
```

**Примеры ошибок:**

* 400 Bad Request — неверные данные (например, отсутствует `title` или `position_in_column` не int)

```json
{
  "error": "invalid request body",
  "details": "title is required"
}
```

* 401 Unauthorized — нет токена
* 404 Not Found — колонка или команда не найдены / пользователь не в команде
* 500 Internal Server Error

**Пример cURL:**

```bash
curl -X POST "https://host/api/v1/team/5/column/2/task" \
 -H "Authorization: Bearer $TOKEN" \
 -H "Content-Type: application/json" \
 -d '{
   "title":"Implement login",
   "description":"OAuth2",
   "position_in_column":0,
   "deadline":"2025-12-01T18:00:00Z",
   "label":"backend",
   "priority":"high"
 }'
```

---

### GET /api/v1/team/{team_id}/task/{task_id}

**Описание:** Возвращает полную информацию о задаче.
**Метод:** `GET`
**Аутентификация:** требуется

**Параметры URL:**

* `team_id`
* `task_id`

**Успешный ответ (200):**

```json
{
  "task": {
    "id": 201,
    "author_id": 11,
    "developer_id": 21,
    "label": "backend",
    "title": "Implement login",
    "description": "OAuth2 + refresh tokens",
    "status": "in_progress",
    "priority": {
      "type": "high",
      "title": "High",
      "hex_color": "#FF0000"
    },
    "created_at": "2025-11-30T10:00:00Z",
    "updated_at": "2025-11-30T12:00:00Z",
    "deadline": "2025-12-01T18:00:00Z",
    "completed_at": null,
    "column_id": 2,
    "position_in_column": 0
  }
}
```

**Коды ошибок:**

* 400 — неверные параметры
* 401 — не авторизован
* 404 — задача не найдена / не принадлежит команде
* 500 — внутренняя ошибка

---

### DELETE /api/v1/team/{team_id}/column/{column_id}/task/{task_id}

**Описание:** Удаляет задачу. После удаления сервер возвращает позицию (index) в колонке, в которую была задача, чтобы фронт мог корректно обновить клиентский state. Также через WebSocket отправляется `task_deleted`.
**Метод:** `DELETE`
**Аутентификация:** требуется

**Параметры URL:**

* `team_id`
* `column_id`
* `task_id`

**Успешный ответ (200):**

```json
{
  "status": "ok"
}
```

**WS-сообщение (дополнительно отправляется в channel команды):**
`type: "delete"` (или `task_deleted` на клиенте) — пример в разделе WebSocket ниже.

**Коды ошибок:**

* 400 — неверные параметры
* 401 — не авторизован
* 403 — нет прав (опционально)
* 404 — задача/колонка не найдены
* 500 — внутренняя ошибка

---

### PUT /api/v1/team/{team_id}/task/{task_id}/title

**Описание:** Обновляет заголовок задачи. Возвращает статус `ok`. После обновления сервер отправит WS-сообщение `title_update`.
**Метод:** `PUT`
**Content-Type:** `application/json`
**Аутентификация:** требуется

**Параметры URL:**

* `team_id`
* `task_id`

**JSON запрос:**

```json
{
  "title": "New concise title"
}
```

**Успешный ответ (200):**

```json
{
  "status": "ok"
}
```

**WS-пейлоуд (пример):**

```json
{
  "type": "title_update",
  "team_id": 5,
  "data": {
    "id": 201,
    "title": "New concise title",
    "column_id": 2,
    "column_position": 0,
    "message":"task updated"
  }
}
```

**Ошибки:**

* 400 — пустой/неправильный `title`
* 401 — не авторизован
* 404 — задача не найдена
* 500 — внутренняя ошибка

---

### PUT /api/v1/team/{team_id}/task/{task_id}/description

**Описание:** Обновляет текст описания задачи. Отправляет `description_update` через WS.
**Метод:** `PUT`
**Content-Type:** `application/json`
**Аутентификация:** требуется

**JSON запрос:**

```json
{
  "description": "Updated detailed description (string)"
}
```

**Успешный ответ (200):**

```json
{
  "status": "ok"
}
```

**WS-пейлоуд:**

```json
{
  "type":"description_update",
  "team_id":5,
  "data": {
    "id":201,
    "description":"Updated detailed description (string)",
    "column_id":2,
    "column_position":0,
    "message":"task updated"
  }
}
```

**Ошибки:** 400/401/404/500 — как обычно.

---

### PUT /api/v1/team/{team_id}/task/{task_id}/deadline

**Описание:** Обновляет дедлайн задачи. Параметр nullable — для удаления дедлайна отправьте `null`. Отправляет `deadline_update` через WS.
**Метод:** `PUT`
**Content-Type:** `application/json`
**Аутентификация:** требуется

**JSON запрос (установить дедлайн):**

```json
{
  "deadline": "2025-12-10T18:00:00Z"
}
```

**JSON запрос (убрать дедлайн):**

```json
{
  "deadline": null
}
```

**Успешный ответ (200):**

```json
{ "status": "ok" }
```

**WS-пейлоуд:**

```json
{
  "type":"deadline_update",
  "team_id":5,
  "data":{
    "id":201,
    "deadline":"2025-12-10T18:00:00Z",
    "column_id":2,
    "column_position":0,
    "message":"task updated"
  }
}
```

**Ошибки:** 400 при неверном формате времени, 401, 404, 500.

---

### PUT /api/v1/team/{team_id}/task/{task_id}/developer

**Описание:** Назначает или убирает (null) разработчика у задачи. Отправляет `developer_update` через WS.
**Метод:** `PUT`
**Content-Type:** `application/json`
**Аутентификация:** требуется

**JSON запрос (назначить):**

```json
{
  "developer_id": 21
}
```

**JSON запрос (снять исполнителя):**

```json
{
  "developer_id": null
}
```

**Успешный ответ (200):**

```json
{ "status": "ok" }
```

**WS-пейлоуд:**

```json
{
  "type":"developer_update",
  "team_id":5,
  "data":{
    "id":201,
    "developer_id":21,
    "column_id":2,
    "column_position":0,
    "message":"task updated"
  }
}
```

**Ошибки:** 400/401/404/500.

---

### PUT /api/v1/team/{team_id}/task/{task_id}/status

**Описание:** Обновляет статус задачи. Поддерживаемые значения enum в БД: `todo`, `in_progress`, `completed`. После смены статуса (если `completed`) может устанавливаться `completed_at`. Отправляет `status_update` через WS.
**Метод:** `PUT`
**Content-Type:** `application/json`
**Аутентификация:** требуется

**JSON запрос:**

```json
{
  "status": "in_progress"  // or "todo" or "completed"
}
```

**Успешный ответ (200):**

```json
{ "status": "ok" }
```

**WS-payload:**

```json
{
  "type":"status_update",
  "team_id":5,
  "data":{
    "id":201,
    "status":"in_progress",
    "column_id":2,
    "column_position":0,
    "message":"task updated"
  }
}
```

**Ошибки:**

* 400 — неверный статус
* 401 — не авторизован
* 404 — задача не найдена
* 500 — внутренняя ошибка

---

### PUT /api/v1/team/{team_id}/task/{task_id}/move

**Описание:** Перемещает задачу в указанную колонку и позицию. Поддерживает:

* перемещение внутри той же колонки (сдвиг промежуточных задач)
* перенос между колонками (сжатие старой колонки и сдвиг в новой)
  DB-логика: у вас уже есть BEFORE INSERT/UPDATE триггеры/функции, которые корректно перераспределяют `position_in_column`. Серверный метод должен вызывать `mgr.UpdateTaskColumn(&req)` или аналог и вернуть `ok`. После успешного изменения отправляется WS-сообщение `column_update`.

**Метод:** `PUT`
**Content-Type:** `application/json`
**Аутентификация:** требуется

**JSON запрос:**

```json
{
  "column_id": 4,          // новая колонка (int)
  "position_in_column": 1  // новая позиция (int)
}
```

> В вашем коде модель запроса называется `TaskColumnUpdateRequest` и возвращает (column_id, position_in_column) из `mgr.UpdateTaskColumn`.

**Успешный ответ (200):**

```json
{ "status": "ok" }
```

**WS-пейлоуд (пример):**

```json
{
  "type":"column_update",
  "team_id":5,
  "data":{
    "id":201,
    "column_id":4,
    "column_position":1,
    "message":"task updated"
  }
}
```

**Ошибки:**

* 400 — неверный JSON
* 401 — не авторизован
* 404 — задача/колонка не найдены
* 409 — конфликт (если позиция занята и механизм не смог корректно сдвинуть) — опционально
* 500 — внутренняя ошибка

---

# WebSocket API (Dashboard)

> Вебсокеты используются для рассыла уведомлений/изменений доски всем подключённым пользователям команды. WS-подключение отправляет сообщения сверху в формате `Envelope` (см. ниже). Сервер использует Redis Pub/Sub с каналами `team:{teamID}` и `team:{teamID}:notifications`.

---

## Подключение

**URL:**

```
ws://host/api/v1/team/{team_id}/dashboard
```

**Query params:** нет обязательных (опционально — можно добавить limit/offset, но реализация пока не требует).
**Заголовки:**

```
Authorization: Bearer {access_token}
```

**Ответ при подключении:** стандартное WebSocket handshake. После успешного апгрейда сервер начнёт слать сообщения про изменения задач.

---

## Формат Envelope (WS)

Все сообщения приходят в виде JSON-объекта:

```json
{
  "team_id": 5,
  "type": "create|delete|title_update|description_update|deadline_update|status_update|developer_update|column_update",
  "data": { /* payload per type */ }
}
```

> Типы (константы в коде):
> `"create"`, `"delete"`, `"title_update"`, `"description_update"`, `"deadline_update"`, `"status_update"`, `"developer_update"`, `"column_update"`.

---

## Типы сообщений и примеры данных

### task created (`type: "create"`)

```json
{
  "team_id": 5,
  "type": "create",
  "data": {
    "id": 201,
    "title": "Implement login",
    "column_id": 2,
    "column_position": 0,
    "deadline": "2025-12-01T18:00:00Z",
    "message": "task created"
  }
}
```

### task deleted (`type: "delete"`)

```json
{
  "team_id": 5,
  "type": "delete",
  "data": {
    "id": 201,
    "column_id": 2,
    "column_position": 3,
    "message": "task deleted"
  }
}
```

### title update (`type: "title_update"`)

```json
{
  "team_id": 5,
  "type": "title_update",
  "data": {
    "id": 201,
    "title": "New concise title",
    "column_id": 2,
    "column_position": 0,
    "message": "task updated"
  }
}
```

### description update (`type: "description_update"`)

```json
{
  "team_id": 5,
  "type": "description_update",
  "data": {
    "id": 201,
    "description": "Full long description text",
    "column_id": 2,
    "column_position": 0,
    "message": "task updated"
  }
}
```

### deadline update (`type: "deadline_update"`)

```json
{
  "team_id": 5,
  "type": "deadline_update",
  "data": {
    "id": 201,
    "deadline": "2025-12-10T18:00:00Z",
    "column_id": 2,
    "column_position": 0,
    "message": "task updated"
  }
}
```

### developer update (`type: "developer_update"`)

```json
{
  "team_id": 5,
  "type": "developer_update",
  "data": {
    "id": 201,
    "developer_id": 21,
    "column_id": 2,
    "column_position": 0,
    "message": "task updated"
  }
}
```

### status update (`type: "status_update"`)

```json
{
  "team_id": 5,
  "type": "status_update",
  "data": {
    "id": 201,
    "status": "in_progress",
    "column_id": 2,
    "column_position": 0,
    "message": "task updated"
  }
}
```

### column / move update (`type: "column_update"`)

```json
{
  "team_id": 5,
  "type": "column_update",
  "data": {
    "id": 201,
    "column_id": 4,
    "column_position": 1,
    "message": "task updated"
  }
}
```

---

# Ошибки и коды ответа — единообразие

Сервис использует одинаковую структуру ошибок в JSON:

**Пример:**

```json
{
  "error": "description of error",
  "details": "optional additional details"
}
```

**Стандартные HTTP-коды:**

* `200 OK` — операции вида GET/PUT/DELETE успешно выполнены
* `201 Created` — ресурс создан (POST)
* `400 Bad Request` — неверные параметры/JSON/типы
* `401 Unauthorized` — отсутствует / недействителен токен
* `403 Forbidden` — нет прав на ресурс (опционально)
* `404 Not Found` — ресурс (team/column/task) не найден
* `409 Conflict` — конфликт состояний (опционально, например, при параллельных изменениях)
* `500 Internal Server Error` — необработанная ошибка на сервере

---

# Примечания для реализации (фронтенд / интеграция)

* Все даты/времена — в ISO 8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`). На фронте — локализовать для показа.
* При создании задачи, если не указана позиция, бэкенд должен ставить `position_in_column = max+1`.
* При перемещениях/вставках сервер держит корректный порядок с помощью триггеров (обновление `position_in_column`) — фронту достаточно принять WS-событие и применить изменения.
* WS-сообщения публикуются через Redis Pub/Sub на канал `team:{teamID}`. Клиентам — подписываться на WS URL `/api/v1/team/{team_id}/dashboard`.
* Любое действие, меняющее задачу, обязательно возвращает/рассылает колонку и позицию задачи (`column_id` и `column_position`) — чтобы фронт мог однозначно разместить карточку.
* Если фронтенд делает optimistic UI (например — мгновенно перемещает карточку), он должен быть готов откатить состояние при приходе WS-сообщения с фактическими данными (в случае конфликтов).

---

# Пример полного рабочего сценария (пошагово)

1. Клиент: `POST /api/v1/team/5/column/2/task` — создаёт задачу в колонке 2 на позицию 0.
2. Сервер: создал задачу с id=201, ответ `201 { "id": 201 }`.
3. Сервер публикует WS-сообщение `create` на канал команды. Все подключённые клиенты получают: задача + column + position.
4. Клиент может затем открыть задачу `GET /api/v1/team/5/task/201` для просмотра полной информации.
5. Другой клиент, поменяв статус `PUT /api/v1/team/5/task/201/status` на `in_progress`, получит `200`, и сервер опубликует `status_update` через WS.

---