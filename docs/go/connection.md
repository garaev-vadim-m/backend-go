# Подключение к PostgreSQL

## Строка подключения

```go
connStr := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
)
```

Аналогия в JavaScript:

```javascript
const connStr = `host=localhost port=5432 user=postgres password=postgres dbname=myapp sslmode=disable`
```

## Параметры

| Параметр      | Значение                                                              |
|---------------|-----------------------------------------------------------------------|
| `host`        | Адрес хоста (в docker-compose это имя сервиса — `postgres`)          |
| `port`        | Порт PostgreSQL (стандартно 5432)                                     |
| `user`        | Пользователь БД (`postgres`)                                          |
| `password`    | Пароль (`postgres`)                                                   |
| `dbname`      | Имя базы данных (`myapp`)                                             |
| `sslmode=disable` | SSL не используем (только для локальной разработки!)             |

## `%s` — плейсхолдер

После `fmt.Sprintf` переменная `connStr` будет содержать:

```
"host=postgres port=5432 user=postgres password=postgres dbname=myapp sslmode=disable"
```
