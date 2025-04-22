# System Overview — Messenger Architecture

## Микросервисы

| Сервис             | Назначение                                                      |
|--------------------|-----------------------------------------------------------------|
| `gateway`          | Единая точка входа, прокси + авторизация                        |
| `auth-service`     | Регистрация, вход, верификация, OAuth                           |
| `session-service`  | Хранение refresh-токенов, обновление и инвалидация сессий      |
| `profile-service`  | Управление публичной информацией о пользователях               |
| `contact-service`  | Социальные связи (друзья, поиск, заявки)                       |
| `messaging-service`| Отправка и получение сообщений, управление чатами              |

---

## События (Event Flow)

| Событие               | Источник         | Подписчики                       | Назначение                          |
|-----------------------|------------------|----------------------------------|-------------------------------------|
| `user.created`        | auth-service     | profile-service, contact-service | Создание профиля, social-записи     |
| `user.email_verified` | auth-service     | —                                | (опционально: уведомления)          |
| `session.created`     | session-service  | —                                | Аудит входа                         |
| `session.revoked`     | session-service  | —                                | Логаут                              |
| `profile.updated`     | profile-service  | —                                | (опционально: кэш, индексация)      |
| `friendship.created`  | contact-service  | —                                | (опционально: уведомления)          |
| `friendship.accepted` | contact-service  | messaging-service                | Разрешение общения в чате           |
| `friendship.removed`  | contact-service  | messaging-service                | Блокировка/отключение сообщений     |
| `message.sent`        | messaging-service| —                                | Уведомления, аналитика              |
| `message.seen`        | messaging-service| —                                | Обновление статуса сообщений        |

---

## Взаимодействие (Sync API)

| Источник           | Вызывает            | Метод                             | Назначение                          |
|--------------------|---------------------|-----------------------------------|-------------------------------------|
| `gateway`          | все сервисы         | proxy HTTP/gRPC                   | Внешний API                         |
| `auth-service`     | session-service     | `POST /session/refresh`           | Выдача токенов                      |
| `gateway`          | session-service     | `DELETE /session/*`               | Логаут                              |
| `messaging-service`| contact-service     | `GET /contacts/:id`               | Проверка дружбы перед отправкой     |

---

## Блок-схема взаимодействий (Mermaid)

```mermaid
flowchart TD
    subgraph External
        Client[Пользователь]
    end

    subgraph EntryPoint
        GW[Gateway Service]
    end

    subgraph AuthLayer
        AUTH[Auth Service]
        SESSION[Session Service]
    end

    subgraph UserLayer
        PROFILE[Profile Service]
        CONTACT[Contact Service]
    end

    subgraph Messaging
        MSG[Messaging Service]
    end

    Client --> GW

    %% Auth
    GW --> AUTH
    AUTH --> SESSION
    GW --> SESSION

    %% Events
    AUTH -- user.created --> PROFILE
    AUTH -- user.created --> CONTACT
    AUTH -- user.email_verified --> PROFILE

    %% Profile
    GW --> PROFILE
    PROFILE -- profile.updated --> /*

    %% Contact
    GW --> CONTACT
    CONTACT -- friendship.accepted --> MSG
    CONTACT -- friendship.removed --> MSG

    %% Messaging
    GW --> MSG
    MSG --> CONTACT

    %% Session
    SESSION -- session.revoked --> Logout
