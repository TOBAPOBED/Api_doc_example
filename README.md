# Todo API Service

REST API для управления задачами (Todo) и пользователями, написанный на Go. Проект демонстрирует использование `swaggo` для автоматической генерации Swagger-документации из аннотаций в коде.

## 📋 Описание

API предоставляет полный набор CRUD-операций для управления пользователями и их задачами. Данные хранятся в оперативной памяти (in-memory map) для демонстрационных целей. Проект разделён на обработчики HTTP-запросов (`handlers`) и модели данных (`models`).

**Особенности:**
- Валидация входных данных (обязательные поля, форматы, ограничения длины)
- Единый формат ответов API (JSON)
- Потокобезопасное хранение данных (`sync.RWMutex`)
- Интерактивная Swagger-документация

## 🚀 Установка и запуск

### Предварительные требования

- Go версии 1.21 или выше
- Установленный CLI-инструмент `swag`:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

### Инструкции по запуску

1. Клонируйте репозиторий:
   ```bash
   git clone <repository-url>
   cd api_doc_example
   ```

2. Установите зависимости:
   ```bash
   go mod tidy
   ```

3. Сгенерируйте Swagger-документацию:
   ```bash
   swag init -g main.go -o docs
   ```

4. Запустите сервер:
   ```bash
   go run main.go
   ```

5. Откройте браузер:
   - Swagger UI: http://localhost:8080/swagger/index.html
   - API JSON: http://localhost:8080/docs/swagger.json

## 📦 Зависимости

Проект использует следующие основные зависимости:

- `github.com/go-chi/chi/v5` v5.0.10 — лёгкий и производительный HTTP-роутер
- `github.com/swaggo/swag` v1.16.2 — генератор Swagger-документации из Go-комментариев
- `github.com/swaggo/http-swagger/v2` v2.0.2 — middleware для отображения Swagger UI

Все зависимости управляются через Go modules (`go.mod` и `go.sum`).

## 📖 API Документация

После запуска сервера интерактивная документация доступна по адресу:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

### Основные эндпоинты

#### Users (Пользователи)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/users` | Получить список всех пользователей |
| GET | `/api/v1/users/{id}` | Получить пользователя по ID |
| POST | `/api/v1/users` | Создать нового пользователя |
| PUT | `/api/v1/users/{id}` | Обновить пользователя |
| DELETE | `/api/v1/users/{id}` | Удалить пользователя |

#### Todos (Задачи)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/todos` | Получить список всех задач |
| GET | `/api/v1/todos/{id}` | Получить задачу по ID |
| POST | `/api/v1/todos` | Создать новую задачу |
| PUT | `/api/v1/todos/{id}` | Обновить задачу |
| DELETE | `/api/v1/todos/{id}` | Удалить задачу |

## 📝 Формат ответов API

Все ответы API следуют единому формату.

### Успешный ответ

```json
{
  "success": true,
  "data": { ... }
}
```

Поле `data` может содержать:
- Один объект (например, при GET по ID)
- Массив объектов (например, при GET списка)

### Ответ с ошибкой

```json
{
  "error": "error_code",
  "message": "Человекочитаемое описание ошибки"
}
```

### Возможные коды ошибок

| Код | HTTP статус | Описание |
|-----|-------------|----------|
| `bad_request` | 400 | Неверный формат запроса (например, некорректный ID) |
| `validation_error` | 400 | Ошибка валидации входных данных |
| `not_found` | 404 | Запрашиваемый ресурс не найден |

## ✅ Валидация входных данных

API проверяет входные данные согласно контракту, описанному в Swagger.

### Пользователь (`CreateUserRequest`)

| Поле | Тип | Обязательность | Ограничения |
|------|-----|----------------|-------------|
| `email` | string | Да | Должен содержать `@` |
| `name` | string | Да | Максимум 100 символов |

### Задача (`CreateTodoRequest`)

| Поле | Тип | Обязательность | Ограничения |
|------|-----|----------------|-------------|
| `user_id` | int | Да | — |
| `title` | string | Да | Максимум 200 символов |
| `description` | string | Нет | Максимум 1000 символов |
| `done` | bool | Нет | — |

## 🔌 Примеры запросов и ответов

### Создать пользователя

**Запрос:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "Иван Иванов"
  }'
```

**Ответ (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "Иван Иванов",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Ответ при ошибке валидации (400 Bad Request):**
```json
{
  "error": "validation_error",
  "message": "Email обязателен"
}
```

### Создать задачу

**Запрос:**
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "title": "Изучить Swagger",
    "description": "Добавить аннотации в код",
    "done": false
  }'
```

**Ответ (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Изучить Swagger",
    "description": "Добавить аннотации в код",
    "done": false,
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  }
}
```

**Ответ при ошибке валидации (400 Bad Request):**
```json
{
  "error": "validation_error",
  "message": "Title не должен превышать 200 символов"
}
```

### Получить все задачи

**Запрос:**
```bash
curl -X GET http://localhost:8080/api/v1/todos
```

**Ответ (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "title": "Изучить Swagger",
      "description": "Добавить аннотации в код",
      "done": false,
      "created_at": "2024-01-15T10:35:00Z",
      "updated_at": "2024-01-15T10:35:00Z"
    }
  ]
}
```

### Получить ресурс по несуществующему ID

**Запрос:**
```bash
curl -X GET http://localhost:8080/api/v1/todos/999
```

**Ответ (404 Not Found):**
```json
{
  "error": "not_found",
  "message": "Задача не найдена"
}
```

## 📂 Структура проекта

```
api-doc-example/
├── docs/                          # Сгенерированная Swagger-документация
│   ├── docs.go                    # Автогенерируемый код для Swagger UI
│   ├── swagger.json               # JSON-спецификация API
│   └── swagger.yaml               # YAML-спецификация API
├── internal/
│   ├── handlers/                  # HTTP-обработчики
│   │   ├── todo.go                # Обработчики для задач
│   │   └── user.go                # Обработчики для пользователей
│   └── models/                    # Модели данных (DTO)
│       ├── todo.go                # Модель Todo и запросы/ответы
│       └── user.go                # Модель User и запросы/ответы
├── main.go                        # Точка входа, роутинг, глобальные Swagger-аннотации
├── go.mod                         # Описание Go-модуля и зависимостей
├── go.sum                         # Контрольные суммы зависимостей
└── README.md                      # Документация проекта
```

### Описание компонентов

- **`main.go`** — точка входа приложения, инициализация роутера Chi, глобальные Swagger-аннотации (`@title`, `@version`, `@BasePath`, `@schemes`)
- **`internal/handlers/`** — HTTP-обработчики с Swagger-аннотациями (`@Summary`, `@Description`, `@Param`, `@Success`, `@Failure`, `@Router`)
- **`internal/models/`** — структуры данных (DTO) с примерами (`example`) для генерации схем в Swagger
- **`docs/`** — автоматически генерируемая папка (не редактировать вручную)

## 🔧 Разработка

### Генерация Swagger-документации

После изменения аннотаций в коде перегенерируйте документацию:

```bash
swag init -g main.go -o docs
```

### Проверка валидации

Для проверки работы валидации можно использовать curl с некорректными данными:

```bash
# Пустой email
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email": "", "name": "Тест"}'

# Слишком длинное имя (более 100 символов)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "name": "AAAA...A (101 символ)"}'

# Пустой title задачи
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "title": "", "description": "Тест"}'
```

Все эти запросы вернут JSON-ответ с кодом 400 и описанием ошибки.

##  Лицензия

Этот проект распространяется под лицензией MIT.

##  Контакты

- Репозиторий: https://github.com/TOBAPOBED/Api_doc_example
- Поддержка API: api-support@example.com