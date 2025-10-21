# Документация по CALLS API

## Система звонков (WebRTC + WebSocket)
Базовый URL: `/api/v1`

> **Важно:** Все эндпойнты требуют аутентификации через Bearer токен в заголовке Authorization.

---

## WebRTC API

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

---

## WebRTC Сигналинг

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

---

## Технические детали

### WebRTC Конфигурация
- **STUN сервер:** `stun:stun.l.google.com:19302` (будет заменен на наш)
- **Поддерживаемые кодеки:**
  - Аудио: Opus, G.722, G.711
  - Видео: VP8, VP9, H.264
- **Направления трансиверов:** Receive-only (по умолчанию)

### Управление комнатами
- Автоматическое создание комнат при первом подключении
- Автоматическое удаление пустых комнат
- Поддержка множественных участников в одной комнате
- Изоляция комнат по типу чата (приватный/групповой)

### Обработка медиа-потоков
- Автоматическое пересылка RTP пакетов между участниками
- Поддержка множественных треков (аудио + видео)
- Автоматическое переподключение при потере соединения
- Буферизация ICE кандидатов до установления соединения

### Состояния соединения
- `new` - новое соединение
- `connecting` - установление соединения
- `connected` - соединение установлено
- `disconnected` - соединение разорвано
- `failed` - ошибка соединения
- `closed` - соединение закрыто

---

## Модели данных

### PeerState
```json
{
  "sessionId": "unique-session-id",
  "userId": 1,
  "tracks": {
    "track-id": {
      "track": "TrackLocalStaticRTP",
      "sender": "RTPSender"
    }
  },
  "pendingICE": [
    {
      "candidate": "candidate:1 1 UDP 2113667326 192.168.1.100 54400 typ host",
      "sdpMLineIndex": 0,
      "sdpMid": "0"
    }
  ]
}
```

### Room
```json
{
  "peers": {
    "session-id": "PeerState"
  },
  "published": {
    "track-id": {
      "remote": "TrackRemote",
      "codec": "RTPCodecCapability",
      "locals": {
        "session-id": "Track"
      }
    }
  }
}
```

### SignalMessage
```json
{
  "type": "join|offer|answer|ice|leave",
  "payload": "json-encoded-data",
  "roomId": 1,
  "sessionId": "unique-session-id"
}
```

---

## Обработка ошибок

### Типичные ошибки
- **WebSocket закрыт неожиданно** - переподключение клиента
- **Неверный JSON** - игнорирование сообщения
- **Комната не найдена** - создание новой комнаты
- **PeerConnection ошибка** - закрытие соединения
- **ICE кандидат не добавлен** - буферизация до установления соединения

### Логирование
- Все WebRTC события логируются
- Отслеживание состояния соединений
- Мониторинг производительности медиа-потоков
- Ошибки соединения с детальной информацией

---

## Безопасность

### Аутентификация
- JWT токен в заголовке Authorization
- Валидация пользователя перед подключением к комнате
- Проверка прав доступа к чату/группе

### Изоляция
- Комнаты изолированы по ID чата/группы
- Невозможность подключения к чужим комнатам
- Автоматическая очистка ресурсов при отключении

### Ограничения
- Максимальное количество участников в комнате: не ограничено
- Таймаут неактивных соединений: 60 секунд
- Автоматическое удаление пустых комнат

---

## Примеры использования

### Инициализация звонка
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/chat/1/call', {
  headers: {
    'Authorization': 'Bearer your-jwt-token'
  }
});

// Присоединение к комнате
ws.send(JSON.stringify({
  type: 'join',
  roomId: 1,
  sessionId: 'unique-session-id'
}));
```

### Обработка WebRTC событий
```javascript
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  
  switch(message.type) {
    case 'offer':
      // Обработка WebRTC offer
      break;
    case 'answer':
      // Обработка WebRTC answer
      break;
    case 'ice':
      // Обработка ICE candidate
      break;
  }
};
```

### Отправка ICE кандидата
```javascript
peerConnection.onicecandidate = (event) => {
  if (event.candidate) {
    ws.send(JSON.stringify({
      type: 'ice',
      payload: event.candidate,
      roomId: 1,
      sessionId: 'unique-session-id'
    }));
  }
};
```

---

## Примечания

### Производительность
- Использование Redis для масштабирования
- Оптимизированная пересылка RTP пакетов
- Минимальная задержка медиа-потоков
- Эффективное управление памятью

### Совместимость
- Поддержка современных браузеров с WebRTC
- Автоматическое определение доступных кодеков
- Fallback для старых браузеров
- Кроссплатформенная совместимость

### Мониторинг
- Метрики качества соединения
- Статистика использования полосы пропускания
- Отслеживание ошибок соединения
- Профилирование производительности
