# hello-go-backend

Минимальный Go-бэкенд с PostgreSQL в Docker.

## Быстрый старт

```bash
docker compose up --build
```

## Веб-панель (`index.html`)

В корне есть `index.html` — самодостаточная страница (чистый JS, без сборки) для
работы с API: CRUD для `users`, `status`, `rules`, `user_status`, `products`, плюс индикатор
`/health`. Откройте файл в браузере (двойным щелчком или через любой статический
сервер) и укажите базовый URL API в поле сверху (по умолчанию
`http://localhost:8080`; значение сохраняется в `localStorage`).

`auth.html` содержит форму входа по email и паролю. После успешной авторизации
сессия сохраняется в `localStorage`, а пользователь переходит в `index.html`.
Операции создания, изменения и удаления в панели доступны только пользователю
с ролью `admin`; это ограничение пока проверяется на стороне клиента.

Запросы из браузера идут с другого origin (`file://` или другой порт), поэтому
бэкенд обёрнут в CORS-middleware (`corsMiddleware` в `main.go`): он отдаёт
`Access-Control-Allow-Origin/-Methods/-Headers` и отвечает `204` на preflight
`OPTIONS`. Если в панели появляется «Сетевая/CORS ошибка» — проверьте, что
бэкенд запущен и доступен по указанному URL.

## Локальная разработка (без Docker)

1. Запустить PostgreSQL и создать БД `myapp`
2. Экспортировать переменные окружения:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
```

Подключение к базе данных
```bash
docker-compose exec postgres psql -U postgres -d myapp
```

Проверка таблиц БД
```bash
docker-compose exec postgres psql -U postgres -d myapp -c "\dt"
```

3. Запустить:

```bash
go run .
```

## Стек

- Go 1.22
- PostgreSQL 15
- Docker / Docker Compose
