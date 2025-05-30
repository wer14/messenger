# Auth Service

## Назначение

Auth Service отвечает за идентификацию пользователей:

- регистрация по email и паролю
- вход по email и паролю
- OAuth-авторизация через сторонние провайдеры (Google, GitHub и т.п.)

Сервис не хранит сессии и не валидирует refresh-токены. Его задача — удостовериться, что пользователь ввёл корректные учётные данные и получить доступ к системе.

## Хранимые данные

- Пользователи:
  - `id`, `email`, `password_hash`, `is_verified`, `created_at`
- OAuth-связи:
  - `id`, `user_id`, `provider`, `external_id`, `created_at`
- Email verification tokens (временные, в Postgres)

## Взаимодействие с другими сервисами

| Сервис            | Назначение                                       |
|-------------------|--------------------------------------------------|
| `session-service` | Делегирование создания и обновления сессий       |
| `gateway`         | Все публичные вызовы идут через него             |
| `profile-service` | (опционально) создание профиля после регистрации |

## Методы / API

| Метод                     | Описание                         |
|---------------------------|----------------------------------|
| `POST /auth/register`     | Регистрация по email и паролю    |
| `POST /auth/login`        | Вход по email и паролю           |
| `POST /auth/oauth`        | Вход через OAuth2                |
| `POST /auth/verify-email` | Подтверждение email              |

## Публикуемые события (NATS)

| Событие                   | Назначение                      |
|---------------------------|---------------------------------|
| `user.created`            | Создание пользователя           |
| `user.email_verified`     | Email подтверждён               |

## Потребляемые события

_Отсутствуют_

## 🔒 Примечания

- Пароли хранятся как `bcrypt`/`argon2` хэши
- Валидация токенов доступа (`access_token`) осуществляется downstream-сервисами локально
