# System Overview — Messenger Architecture

## 📡 Communication Overview

## 🧱 Подход к межсервисному взаимодействию

В рамках микросервисной архитектуры мессенджера используются следующие методы коммуникации:

- **gRPC** — для синхронных команд между сервисами
- **REST** — для публичных данных, которые удобно кэшировать (например, профили)
- **Очередь сообщений (Kafka или NATS)** — для событийной коммуникации между сервисами, где требуется масштабируемость, надёжность или fan-out

---

## 🔧 Используемые протоколы по сервисам

| Сервис             | Протоколы       | Обоснование                                                                 |
|--------------------|------------------|------------------------------------------------------------------------------|
| `gateway`          | REST (HTTP/1.1)  | Универсальный интерфейс для фронта, токен в заголовке Authorization         |
| `auth-service`     | gRPC             | Простые команды (login, register, verify), нет необходимости в REST        |
| `session-service`  | gRPC             | Валидация, refresh и удаление сессий — командный стиль                     |
| `profile-service`  | REST             | Публичная информация, можно кэшировать, хорошо ложится в REST              |
| `contact-service`  | gRPC             | Социальные действия — чистые команды                                        |
| `messaging-service`| gRPC + очередь   | gRPC для отправки/чтения, очередь для доставки событий и интеграций        |

---

## 🔁 Событийная модель (Event-driven)

События публикуются через очередь (Kafka или NATS), чтобы обеспечить:

- надёжную доставку и масштабирование
- асинхронную обработку
- независимость между сервисами

| Источник          | Событие               | Подписчики             | Назначение                                  |
|-------------------|-----------------------|-------------------------|--------------------------------------------|
| `auth-service`    | `user.created`        | profile, contact        | Создание базовых записей                   |
|                   | `user.email_verified` | —                       | Статус email                               |
| `profile-service` | `profile.updated`     | (опционально)           | Кэширование или индексация                 |
| `contact-service` | `friendship.created`  | (опционально)           | Уведомления                                |
|                   | `friendship.accepted` | messaging-service       | Разрешение чата                            |
|                   | `friendship.removed`  | messaging-service       | Прекращение доступа к сообщениям           |
| `messaging-service` | `message.sent`      | (уведомления, аналитика)| Событие отправки сообщения                 |
|                   | `message.seen`        | (аналитика)             | Обновление статуса                         |
| `session-service` | `session.created`     | (опционально)           | Аудит входов                               |
|                   | `session.revoked`     | (опционально)           | Принудительный выход                       |

---

## ⚙️ gRPC команды (синхронные вызовы)

| Откуда           | Куда              | Метод                    | Назначение                          |
|------------------|-------------------|---------------------------|-------------------------------------|
| `gateway`        | `auth-service`    | `Login`, `Register`, `VerifyEmail` | Аутентификация             |
| `gateway`        | `session-service` | `RefreshSession`, `RevokeSession`  | Обновление и выход из сессии        |
| `gateway`        | `profile-service` | `GetProfile`, `UpdateProfile`      | Публичный профиль                   |
| `gateway`        | `contact-service` | `AddFriend`, `RemoveFriend`, `GetFriends` | Управление связями         |
| `gateway`        | `messaging-service` | `SendMessage`, `GetMessages`, `MarkSeen` | Работа с перепиской     |
| `messaging-service` | `contact-service` | `IsFriend(userA, userB)`        | Проверка прав перед отправкой       |

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
