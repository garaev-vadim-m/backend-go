# hello-go-backend

Небольшой Go-бэкенд с PostgreSQL и Docker.

## Быстрый старт

```bash
docker compose up --build
```

## Веб-страницы

Статика встроена в бинарник (`go:embed web`) и раздаётся тем же сервером на `http://localhost:8080`: `/` → `web/index.html`, также доступны `/auth.html` и `/admin_panel.html`. Приветствие из старого `/` перенесено на `/hello`.

- `web/index.html` — простая витрина интернет-магазина с карточками товаров (данные из `/products`) и кнопкой «В корзину».
- `web/cart.html` — просмотр корзины: состав, цена, количество, сумма по позиции и итог из `GET /cart`.
- `web/admin_panel.html` — самодостаточная страница (ванильный JS) для работы с API: CRUD для `users`, `status`, `rules`, `user_status`, `products` (включая `price` / `stock` в таблице и форме), проверка `/health`.
- `web/auth.html` — форма входа по email и паролю (`/user/login`).

После успешной авторизации сервер устанавливает `HttpOnly`-куку `session` (`SameSite=Lax`, `Secure` при `COOKIE_SECURE=true`). Токен в `localStorage` не сохраняется и заголовок `Authorization` не используется: все запросы к API идут с `credentials: "include"`. Роль и отображаемого пользователя страницы берут из `GET /user/me`; пользователь с ролью `admin` перенаправляется в `/admin_panel.html`. Операции создания/изменения/удаления доступны только роли `admin`: панель скрывает кнопки, а сервер проверяет роль в middleware на всех CRUD-роутах (кука `session` или `Authorization: Bearer`). `GET /users` тоже доступен только админу; `GET /user?id=` требует авторизацию, и обычный пользователь видит только свою карточку. Панель для не-админа не загружает список пользователей и блокирует карточки. Панель проверяет длину пароля (8–72 байта) и лимиты полей из схемы БД; при `401` она очищает локальные данные сессии и возвращает на `/auth.html`, при `403` показывает ошибку прав без выхода из аккаунта. Выход выполняется через `POST /user/logout`, который отзывает куку на сервере.

Бэкенд предоставляет `GET /cart` (просмотр корзины), `POST /cart/items`
(добавление товара, только для авторизованных) и `POST /order/checkout`
(оформление заказа с обязательным заголовком `Idempotency-Key`). На витрине
есть кнопка «В корзину» и переход на `cart.html`, где показываются состав,
цена, количество, сумма по позиции и итог. Неавторизованного пользователя
корзина отправляет на `auth.html`. Изменять количество и удалять товары пока
нельзя, оформление заказа через UI тоже не подключено; в карточках витрины
цена и остаток пока не показываются, хотя API товаров их отдаёт. Детали и
оставшиеся пункты — в
`docs/project/tasks.md` (E2, E3, E6) и `docs/project/implemented.md`.

Страницы обращаются к API по `API base URL`; по умолчанию используется пустое значение, то есть текущий origin (браузер ходит на тот же `http://localhost:8080`, cookie отправляются автоматически). Явное значение хранится в `localStorage` (`hello-go-backend.apiBase`) и нужно только для отдельного фронтенд-сервера; в этом случае его origin должен быть перечислен в `CORS_ORIGINS`, иначе браузер заблокирует ответ. Сервер отдаёт CORS-заголовки (`corsMiddleware` в `main.go`) и отвечает `204` на preflight `OPTIONS`. Открывать HTML через `file://` не рекомендуется: origin `null` не разрешён и cookie-сессия не работает надёжно.

Переменные окружения помимо `DB_*`: `ADMIN_EMAIL`, `ADMIN_PASSWORD`, `ADMIN_NAME` (сид технического админа), `COOKIE_SECURE` (по умолчанию `false`; включите за HTTPS) и `CORS_ORIGINS` (список origin через запятую для отдельного фронтенда).

## Локальная разработка (без Docker)

1. Поднимите PostgreSQL и создайте БД `myapp`
2. Экспортируйте переменные окружения:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
```

Подключение к базе в контейнере:
```bash
docker compose exec postgres psql -U postgres -d myapp
```

Проверить таблицы:
```bash
docker compose exec postgres psql -U postgres -d myapp -c "\dt"
```

3. Запуск:

```bash
go run .
```

Создание миграции
```bash
./scripts/new-migration.sh add_something
```

## Стек

- Go 1.22
- PostgreSQL 15
- Docker / Docker Compose
