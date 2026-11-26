# Проект hello-go-backend

## Что это

REST API на Go с PostgreSQL и тремя статическими страницами без сборки:
витрина товаров (`index.html`), вход/регистрация (`auth.html`) и админ-панель
(`admin_panel.html`). Страницы встроены в бинарник через `go:embed web` и
раздаются тем же сервером, что и API. API реализует CRUD для пяти сущностей,
аутентификацию на серверных сессиях и базовые ручки корзины/заказов
(`GET /cart`, `POST /cart/items`, `POST /order/checkout`). Витрина позволяет
добавить товар в корзину и посмотреть её на отдельной странице; изменение
состава и оформление заказа через UI пока не реализованы.

- API и страницы: `http://localhost:8080`
- Статика: `web/index.html`, `web/auth.html`, `web/cart.html`, `web/admin_panel.html`
- Документация: `docs/`

> Публичны витрина и справочники (`/products`, `/statuses`, `/rules`,
> `/user-statuses`) и карточки товаров/справочников. `GET /users` — только для
> `admin`; `GET /user?id=` и `GET /user/me` требуют авторизацию (обычный
> пользователь видит только себя). Все `POST`/`PUT`/`DELETE` сущностей — только
> `admin`. `/user/logout` читает только куку `session`; остальные защищённые
> роуты принимают куку или `Authorization: Bearer`. Долги — в `tasks.md`.

## Стек

| Компонент   | Версия / технология                |
|-------------|------------------------------------|
| Язык        | Go 1.22                            |
| БД          | PostgreSQL 15                      |
| Драйвер БД  | `github.com/lib/pq`                |
| Миграции    | `github.com/golang-migrate/migrate/v4` (SQL в `embed.FS`) |
| Пароли      | `golang.org/x/crypto/bcrypt`       |
| Контейнеры  | Docker, Docker Compose             |
| Фронтенд    | HTML + vanilla JS (без сборки)     |

## Архитектура

```
main.go
   ├─ database.InitDB / RunMigrations / SeedAdmin
   ├─ corsMiddleware        CORS + preflight OPTIONS
   ├─ http.ServeMux         роутинг "GET /user" (Go 1.22)
   └─ HomeHandler / HealthHandler

internal/
   ├─ httpx/                ParseID, WriteJSON, WriteError
   ├─ errs/                 ErrNotFound, ErrUnauthorized
   ├─ auth/                 HashPassword, CheckPassword, HashOptional, ValidatePassword
   ├─ middleware/           Auth.Require, Auth.RequireRole, IdentityFromContext
    ├─ session/              Create, Validate, Renew, Revoke (SHA-256 токенов в БД)
    ├─ database/
    │    ├─ database.go      подключение с retry
    │    ├─ migrations.go    golang-migrate + migrations/*.sql (15 версий)
    │    └─ seed.go          роли admin/user и технический пользователь
    ├─ cart/                добавление товара в корзину
    ├─ order/               checkout: транзакция, остатки, идемпотентность
    └─ <entity>/
        ├─ models.go        структуры запроса/ответа
        ├─ repository.go    SQL (Query/QueryRow/Exec)
        └─ controller/      HTTP-обработчики + Handler
```

Слои: `controller` → `repository` → PostgreSQL. Контроллер отвечает за HTTP и
коды ответов, репозиторий — за SQL и маппинг строк в модели. Аутентификация
живёт в `user/controller` и использует пакеты `auth` и `session`.
`session.Validate` вызывается в `middleware.Auth.Require`: он возвращает id и
код роли, кладёт `Identity{ID, Role}` в `context.Context`, обновляет сессию и
перевыпускает куку. `GET /user/me` берёт id из этого контекста.
`Auth.RequireRole("admin")` сравнивает `Identity.Role` и возвращает `403` без
вызова обработчика; `GET /user?id=` дополнительно разрешает обычному
пользователю только свою карточку (чужие — `404`).

## Структура проекта

```
backend-go/
├─ main.go                 точка входа, маршруты, CORS, раздача статики
├─ web/                    статические страницы (go:embed)
│  ├─ index.html           витрина товаров, кнопка «В корзину»
│  ├─ auth.html            вход и регистрация
│  ├─ cart.html            просмотр корзины
│  └─ admin_panel.html     CRUD-панель (+ /health)
├─ internal/
│  ├─ auth/                bcrypt и валидация пароля
│  ├─ cart/                корзина: AddItem
│  ├─ database/            InitDB, CloseDB, RunMigrations, SeedAdmin, migrations/
│  ├─ order/               оформление заказа: Checkout
│  ├─ errs/                общие ошибки
│  ├─ httpx/               HTTP/JSON-хелперы
│  ├─ middleware/          проверка сессии и роли
│  ├─ product/             товары
│  ├─ rules/               роли
│  ├─ session/             серверные сессии
│  ├─ status/              статусы товаров
│  ├─ user/                пользователи, login/register/logout
│  └─ user_status/         статусы пользователей
├─ docs/                   документация
├─ Dockerfile              multi-stage сборка
└─ docker-compose.yml      app + postgres
```

Каждая сущность устроена одинаково: `models.go`, `repository.go`,
`controller/{handler,create,get,get_all,update,delete}.go`. У `user` есть
дополнительные файлы `login.go`, `register.go`, `logout.go`.

## Сущности

| Сущность            | Таблица       | Связи / назначение                            |
|---------------------|---------------|-----------------------------------------------|
| Пользователь        | `users`       | `status_id → user_status`, `rule_id → rules`  |
| Статус пользователя | `user_status` | справочник                                    |
| Роль                | `rules`       | справочник (`admin`, `user`)                  |
| Статус товара       | `status`      | справочник                                    |
| Товар               | `products`    | `status_id → status`; есть `price`, `stock`   |
| Сессия              | `sessions`    | `user_id → users` (не CRUD-сущность)          |
| Корзина             | `carts`       | `user_id → users` (одна на пользователя)      |
| Позиция корзины     | `cart_items`  | `cart_id → carts`, `product_id → products`    |
| Заказ               | `orders`      | `user_id → users`, `idempotency_key`          |
| Позиция заказа      | `order_items` | `order_id → orders`, `product_id → products`  |

Подробная схема — `docs/database/schema.md`, картинка — `docs/ERD.png`.
Актуальный перечень колонок и типов — в SQL-миграциях
`internal/database/migrations/`.

## API

`id` передаётся query-параметром, например `GET /user?id=1`.

| Метод  | Путь                      | Описание                                    |
|--------|---------------------------|---------------------------------------------|
| GET    | `/`                       | витрина (`web/index.html`)                  |
| GET    | `/auth.html`              | вход/регистрация                            |
| GET    | `/admin_panel.html`       | админ-панель                                |
| GET    | `/cart.html`              | просмотр корзины                            |
| GET    | `/hello`                  | приветствие (старый `/`)                    |
| GET    | `/health`                 | ping БД                                     |
| POST   | `/user/register`          | регистрация, роль `user`                    |
| POST   | `/user/login`             | вход, ответ `{"ok": true}`, кука `session`  |
| POST   | `/user/logout`            | выход по куке `session`                     |
| GET    | `/user/me`                | текущий пользователь по куке `session`      |
| GET    | `/users`                  | список пользователей                        |
| GET    | `/user?id=`               | пользователь по id                          |
| POST   | `/user/create`            | создать пользователя                        |
| PUT    | `/user/update?id=`        | обновить пользователя                       |
| DELETE | `/user/delete?id=`        | удалить пользователя                        |
| GET    | `/statuses`               | список статусов товаров                     |
| GET    | `/status?id=`             | статус товара                               |
| POST   | `/status/create`          | создать статус                              |
| PUT    | `/status/update?id=`      | обновить статус                             |
| DELETE | `/status/delete?id=`      | удалить статус                              |
| GET    | `/rules`                  | список ролей                                |
| GET    | `/rule?id=`               | роль по id                                  |
| POST   | `/rule/create`            | создать роль                                |
| PUT    | `/rule/update?id=`        | обновить роль                               |
| DELETE | `/rule/delete?id=`        | удалить роль                                |
| GET    | `/user-statuses`          | список статусов пользователей               |
| GET    | `/user-status?id=`        | статус пользователя                         |
| POST   | `/user-status/create`     | создать статус                              |
| PUT    | `/user-status/update?id=` | обновить статус                             |
| DELETE | `/user-status/delete?id=` | удалить статус                              |
| GET    | `/products`               | список товаров                              |
| GET    | `/product?id=`            | товар по id                                 |
| POST   | `/product/create`         | создать товар                               |
| PUT    | `/product/update?id=`     | обновить товар                              |
| DELETE | `/product/delete?id=`     | удалить товар                               |
| GET    | `/cart`                   | текущая корзина (авторизованным)             |
| POST   | `/cart/items`             | добавить товар в корзину (авторизованным)    |
| POST   | `/order/checkout`         | оформить заказ, нужен `Idempotency-Key`      |

Публичны `/products`, `/product?id=`, `/statuses`, `/status?id=`, `/rules`,
`/rule?id=`, `/user-statuses`, `/user-status?id=`. `GET /users` требует роли
`admin`; `GET /user?id=` и `GET /user/me` требуют авторизацию, при этом
не-админ видит только себя. Создание/обновление/удаление любой сущности — только
`admin`. `GET /cart`, `POST /cart/items` и `POST /order/checkout` требуют
авторизацию.
Токен берётся из куки `session` или заголовка
`Authorization: Bearer <токен>`; `/user/logout` читает только куку.

Коды ответов: `200` — `GET`/`PUT` и login, `201` — create и register, `204` —
delete и logout, `400` — неверный JSON/`id`/обязательные поля/слабый пароль,
`401` — неверные учётные данные, отсутствующая/истёкшая сессия, `403` — у
пользователя нет роли `admin`, `404` — сущность не найдена, `409` — email уже
зарегистрирован при register, `500` — внутренняя ошибка. `POST /order/checkout` отдаёт `201` при создании
заказа и `200` при повторе с тем же `Idempotency-Key`; `400` — пустая корзина,
`409` — товар закончился. `GET /cart` отдаёт `200` с полями `items` (пустой
массив для пустой корзины) и `total`; `POST /cart/items` отвечает `200` без тела. Ошибки
отдаются в формате `{"error": "..."}` (кроме текстовых ответов `/health` и
списков при ошибке БД).

Админ-панель проверяет длину пароля (8–72 байта, как и бэкенд), обязательность
`name`/`code` у справочников и лимиты строковых полей из схемы БД до отправки
запроса. Для не-админа она не загружает список пользователей и блокирует
карточки. При `401` очищает локальную сессию и возвращает пользователя на
`auth.html`; при `403` показывает ошибку недостатка прав, оставаясь в сессии.

### Готовность витрины к корзине и заказам

В `web/index.html` у каждой карточки есть кнопка «В корзину» (`POST /cart/items`),
а в шапке — переход на `web/cart.html`. Страница корзины получает `GET /cart`
и показывает товары, цену, количество, сумму по позиции и итог;
неавторизованный пользователь направляется на `auth.html`. Изменить количество или удалить товар
пока нельзя: бэкенд не реализовал `PUT` / `DELETE` для корзины. Ответы
`/products` и `/product?id=` по-прежнему не содержат `price` / `stock`, поэтому
в карточках витрины нет цены и остатка; цена в корзине приходит из `GET /cart`.
Оформление заказа в UI не подключено: для `POST /order/checkout` нужен заголовок
`Idempotency-Key`. Оставшиеся пункты — в `docs/project/tasks.md` (E2, E3, E6) и
`docs/project/implemented.md`.

Формат ответа пользователя (статус и роль развёрнуты, `null` если не назначены;
`created_at` / `updated_at` в RFC3339 и видны в таблице панели; `password` в
ответ не попадает):

```json
{
  "id": 1,
  "name": "Vadim",
  "surname": "Garaev",
  "female": null,
  "email": "user@example.com",
  "created_at": "2026-11-19T12:34:56.789Z",
  "updated_at": "2026-11-19T12:34:56.789Z",
  "status": { "id": 1, "name": "Активен", "code": "active" },
  "rule": { "id": 1, "name": "Admin", "code": "admin" }
}
```

## Переменные окружения

| Переменная       | Значение по умолчанию | Назначение                     |
|------------------|-----------------------|--------------------------------|
| `DB_HOST`        | —                     | хост PostgreSQL                |
| `DB_PORT`        | —                     | порт PostgreSQL                |
| `DB_USER`        | —                     | пользователь БД                |
| `DB_PASSWORD`    | —                     | пароль БД                      |
| `DB_NAME`        | —                     | имя БД                         |
| `ADMIN_EMAIL`    | `email@email`         | email технического админа      |
| `ADMIN_PASSWORD` | `1235`                | пароль технического админа     |
| `ADMIN_NAME`     | `tech`                | имя технического админа        |
| `COOKIE_SECURE`  | `false`               | флаг `Secure` для куки `session` |
| `CORS_ORIGINS`   | пусто                 | разрешённые origin через запятую |

При старте приложение применяет миграции и через `SeedAdmin` создаёт роли
`admin`/`user` и технического пользователя с bcrypt-хешем пароля. Если пароль
в БД отличается от `ADMIN_PASSWORD`, он перезаписывается.

## Запуск

```bash
docker compose up --build
```

Приложение и страницы: `http://localhost:8080`, PostgreSQL: `localhost:5432`, БД `myapp`. Открывать HTML через `file://` не нужно — cookie-сессия и same-origin API-запросы работают только через HTTP.

Локально без Docker:

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=myapp
go run .
```

## Документация
 
| Файл                          | О чём                         |
|-------------------------------|-------------------------------|
| `docs/project/implemented.md` | что уже реализовано           |
| `docs/project/tasks.md`       | задачи на будущее             |
| `docs/database/schema.md`     | схема БД и связи              |
| `docs/go/`                    | шпаргалки по Go               |
| `docs/docker/`                | шпаргалки по Docker           |
