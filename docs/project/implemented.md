# Что уже реализовано

Отметки: `[x]` — готово, `[ ]` — не сделано (см. `tasks.md`).

## Инфраструктура

- [x] Go-модуль `hello-go-backend`, Go 1.22
- [x] Подключение к PostgreSQL через `database/sql` + `lib/pq`
- [x] Retry подключения: 5 попыток с паузой 2 с (`internal/database/database.go`)
- [x] Миграции через `golang-migrate` v4: 15 версий, SQL-файлы встроены в
      бинарник через `embed.FS` и применяются при старте
      (`internal/database/migrations.go`)
- [x] Сидирование при старте: роли `admin`/`user` и технический пользователь
      с bcrypt-паролем (`internal/database/seed.go`, env `ADMIN_*`)
- [x] Dockerfile с multi-stage сборкой (builder → alpine)
- [x] `docker-compose.yml`: app + `postgres:15-alpine`, healthcheck, named volume
- [x] `.dockerignore`
- [x] Статика (`web/*.html`) встроена в бинарник через `embed.FS` и раздаётся
      тем же сервером: `/`, `/auth.html`, `/cart.html`, `/admin_panel.html`;
      приветствие перенесено на `/hello` (`main.go`)
- [x] Endpoint `/health` с ping БД
- [x] CORS-middleware и preflight `OPTIONS` (`main.go`), заголовок
      `Authorization` разрешён; `Access-Control-Allow-Origin` отражает `Origin`

## Модель данных

- [x] Таблицы: `users`, `user_status`, `rules`, `status`, `products`, `sessions`,
      `carts`, `cart_items`, `orders`, `order_items`
- [x] Миграции: 15 версий в `internal/database/migrations/` (создание таблиц,
      `users.password`, `sessions`, `sessions.revoked_at` + индекс,
      уникальность `users.email`, `created_at`/`updated_at`, индекс `expires_at`,
      `products.price`, корзина, `products.stock`, заказы)
- [x] `products.price` (`NUMERIC(10,2)`, ≥ 0) и `products.stock` (`INTEGER`, ≥ 0):
      есть в модели товара и принимаются в `/product/create` /
      `/product/update`
- [x] Внешние ключи: `users.status_id`, `users.rule_id`, `products.status_id`,
      `sessions.user_id` и `carts.user_id` (`ON DELETE CASCADE`),
      `cart_items.cart_id` / `cart_items.product_id` (`ON DELETE CASCADE`),
      `orders.user_id`, `order_items.order_id` (`ON DELETE CASCADE`) и
      `order_items.product_id`
- [x] Nullable-поля: `users.surname`, `users.female`, `users.email`,
      `users.password`, `users.status_id`, `users.rule_id`,
      `products.description`, `products.status_id`, `sessions.revoked_at`
- [x] `users.email` уникален (PostgreSQL допускает несколько `NULL`)
- [x] ERD: `docs/ERD.png`, текстовое описание — `docs/database/schema.md`

## API и CRUD

- [x] Полный CRUD для `users` — `/users`, `/user`, `/user/{create,update,delete}`
- [x] Полный CRUD для `status` — `/statuses`, `/status`, `/status/{create,update,delete}`
- [x] Полный CRUD для `rules` — `/rules`, `/rule`, `/rule/{create,update,delete}`
- [x] Полный CRUD для `user_status` — `/user-statuses`, `/user-status`, `/user-status/{create,update,delete}`
- [x] Полный CRUD для `products` — `/products`, `/product`, `/product/{create,update,delete}`;
      ответ содержит `price` (строка) и `stock`, create/update принимают оба поля
- [x] `GET /cart` — просмотр корзины (только для авторизованных): объект
      `{"items":[{product_id,name,price,quantity,total}],"total"}`, цена и сумма
      приходят строками из PostgreSQL `NUMERIC`
- [x] `POST /cart/items` — добавить товар в корзину (только для авторизованных);
      повторное добавление увеличивает `quantity`
- [x] `DELETE /cart/items?product_id=` — удалить позицию из корзины (только для
      авторизованных); `204 No Content` при успехе, `404` если позиции нет
- [x] `POST /order/checkout` — оформить заказ из корзины; обязателен заголовок
      `Idempotency-Key` (≤ 64 символов), ответ `Order{id,status,total,created_at}`
- [x] Регистрация `POST /user/register`, вход `POST /user/login`,
      выход `POST /user/logout`
- [x] Роутинг с методами Go 1.22 (`"GET /user"`)
- [x] `POST .../create` и register возвращают `201 Created`
- [x] `GET` / `PUT` возвращают `200 OK`
- [x] `DELETE` и logout возвращают `204 No Content`
- [x] `GET /cart` возвращает `200 OK` даже для пустой корзины (`items: []`)
- [x] `POST /cart/items` возвращает `200 OK` без тела; `DELETE /cart/items` —
      `204 No Content` (`404`, если позиции нет); `POST /order/checkout` —
      `201 Created` при новом заказе и `200 OK` при повторе с тем же ключом
- [x] `401` при неверных учётных данных и при logout без токена

## Аутентификация

- [x] Пароли хранятся как bcrypt-хеши (`internal/auth`), длина пароля 8–72 байт
- [x] При входе поддерживается старый пароль в открытом виде (сравнение
      `subtle.ConstantTimeCompare`), после входа выдаётся новая сессия
- [x] Регистрация не принимает `rule_id`/`status_id`: роль `user` назначается
      через `rules.GetIDByCode`
- [x] Дублирующийся email даёт `409 Conflict` (`pq.Error` код `23505`)
- [x] Сессия: случайный токен 32 байта (hex), в БД хранится только его
      SHA-256-хеш (`internal/session`)
- [x] TTL сессии 24 часа; `Validate` учитывает `expires_at`, `revoked_at` и
      возвращает id и код роли; вызывается в `middleware.Auth.Require`
- [x] Продление сессии: `Renew` сдвигает `expires_at`, если до конца меньше 12 ч;
      кука перевыпускается при куки-авторизации; `/user/me` тоже проходит через
      middleware, поэтому работает и с `Bearer`
- [x] Logout отзывает сессию через `revoked_at`, строка остаётся в истории
- [x] `errs.ErrUnauthorized` для невалидной сессии
- [x] `GET /user/me`, `Auth.Require`, `Auth.RequireRole("admin")` —
      защищённые роуты, в контексте `Identity{ID, Role}`
- [x] `GET /cart`, `POST /cart/items`, `DELETE /cart/items` и
      `POST /order/checkout` защищены `Auth.Require`; `user_id` берётся из
      контекста, а не из тела запроса
- [x] `GET /users` — только `admin`; `GET /user?id=` — авторизованным, не-админ
      видит только свою карточку (чужая отдаётся как `404`)
- [x] Кука `session`: `HttpOnly`, `SameSite=Lax`, `Secure` при `COOKIE_SECURE=true`

## Паттерны и качество кода

- [x] Единый пакет ошибок `internal/errs` с `ErrNotFound` и `ErrUnauthorized`
- [x] Обёртка ошибок через `%w` и проверка `errors.Is` (не везде: 4 репозитория
      UPDATE не оборачивают `ErrNotFound`)
- [x] HTTP/JSON-хелперы `internal/httpx`: `ParseID`, `ParseInt`, `WriteJSON`,
      `WriteError`
- [x] Формат ошибки `{"error": "..."}` (кроме текстовых `/health` и списков
      при ошибке БД)
- [x] Checkout использует транзакцию, `SELECT ... FOR UPDATE` для строк корзины,
      списание остатка с условием `stock >= quantity`, снимок цены и идемпотентность
      по `(user_id, idempotency_key)`; `400` — пустая корзина, `409` — нет остатка
- [x] Handler-структура с `*sql.DB` и конструктором `NewHandler`
- [x] Общий `scanner`-интерфейс для `*sql.Row` и `*sql.Rows`
- [x] JOIN статусов/ролей в выборках (развёрнутые объекты вместо `*_id`)
- [x] `users.created_at` / `users.updated_at` (миграции 000010/000011) попадают
      в ответ; `updated_at` пересчитывается в `UPDATE`
- [x] `INSERT` / `UPDATE` через CTE + `RETURNING` + `JOIN`
- [x] `sql.Null*` для nullable-полей из JOIN
- [x] Проверка `id` в query (`httpx.ParseID`) и `product_id` в
      `DELETE /cart/items` (`httpx.ParseInt`)
- [x] Проверка обязательного `name` при создании/обновлении
- [x] Проверка `price` (число с ≤ 2 знаками, максимум `99999999.99`) и
      `stock` (≥ 0) в create/update товара — `400` при неверном значении
- [x] Проверка пароля при создании/обновлении пользователя; пустой пароль
      при `PUT` оставляет старый (`COALESCE(NULLIF($5, ''), password)`)
- [x] Проверка `RowsAffected` при удалении
- [x] Логирование внутренних ошибок через `log.Printf`

## Веб-страницы

- [x] `web/index.html`: витрина товаров из `/products`, кнопка «В корзину»,
      ссылка на `cart.html` со счётчиком товаров, ссылка на админ-панель для
      роли `admin`, индикатор сессии
- [x] `web/cart.html`: просмотр корзины из `GET /cart` — товары, цена, количество,
      сумма по позиции и итог; кнопка «Удалить» вызывает `DELETE /cart/items`;
      пустая корзина и `401` → `auth.html`
- [x] `web/auth.html`: вход и регистрация; серверная сессия живёт в
      `HttpOnly`-куке `session`, в `localStorage` (`hello-go-backend.auth`)
      сохраняется только профиль для UI, админ перенаправляется в
      `/admin_panel.html`
- [x] `web/admin_panel.html`: CRUD для `users`, `status`, `rules`, `user_status`,
      `products`, проверка `/health`, выход из сессии; для не-админа список users
      не загружается, а карточки пользователей заблокированы
- [x] Форма и таблица товаров в админ-панели показывают `price` / `stock`;
      поля редактируются, значение цены отправляется строкой и проверяется
      по правилам бэкенда (до двух знаков, максимум `99999999.99`)
- [x] Текущий пользователь и роль во всех страницах берутся из `GET /user/me`
- [x] Базовый URL API по умолчанию пустой (текущий origin), явное значение
      сохраняется в `localStorage`; запросы идут с `credentials: "include"`,
      `Authorization` не используется
- [x] Кнопки создания/изменения/удаления в панели доступны только роли `admin`;
      сервер дублирует проверку через `Auth.RequireRole("admin")`
- [x] Пароль в форме пользователей проверяется по правилам бэкенда: 8–72 байта
      (`minlength`/`maxlength` + проверка `TextEncoder` на отправке); пустой
      пароль при редактировании оставляет прежний
- [x] Строковые поля ограничены длинами из схемы БД: 255 символов, `female` — 10;
      `name` и `code` в справочниках обязательны (в БД `code` — `NOT NULL`)
- [x] Таблица пользователей показывает `created_at` / `updated_at` в локальном
      формате; карточкой по ID эти поля тоже видны в JSON
- [x] При `401` обработчик `api()` очищает локальную сессию и возвращает на
      `auth.html`; при `403` показывает ошибку прав без выхода; logout также
      идёт через этот обработчик
- [x] Заголовок панели показывает фактическую роль из сессии, а не всегда `admin`;
       атрибут `hidden` гарантированно скрывает элементы

### Готовность фронтенда к корзине и заказам

- [x] Все страницы продолжают работать с текущими публичными ручками и
      `/user/me`: новые бэкенд-роуты не сломали существующий UI
- [x] `web/index.html`: кнопка «В корзину» вызывает `POST /cart/items`,
      счётчик в шапке обновляется из `GET /cart`, есть переход на `cart.html`
- [x] `web/cart.html`: базовый просмотр корзины — состав, цена, количество,
      сумма по позиции и итог; пустая корзина показывает отдельное сообщение
- [x] `web/cart.html`: удаление позиции через кнопку «Удалить» и
      `DELETE /cart/items?product_id=`, список перезагружается после ответа `204`
- [x] API товаров отдаёт `price` / `stock`, а админ-панель позволяет их
      создавать и редактировать; корзина берёт цену из `GET /cart`
- [ ] Карточки витрины `index.html` пока не показывают цену и остаток,
      хотя поля уже приходят в `GET /products`
- [ ] Нельзя изменить количество и очистить корзину: бэкенд пока не даёт
      `PUT /cart/items` и `DELETE /cart`
- [ ] Нет кнопки оформления: `POST /order/checkout` с `Idempotency-Key`
      не вызывается, а `api()` в админ-панели не принимает свои заголовки

Итог: добавление, просмотр и удаление позиций корзины готовы; изменение
количества, очистка корзины и checkout в UI ещё не реализованы. Оставшиеся
пункты — в `tasks.md` (E2/E3/E6).

## Документация

- [x] `README.md`
- [x] `docs/go/` — connection, fmt, pointers, syntax, errors, interfaces, sql, json, http
- [x] `docs/docker/` — basics, build-run, compose
- [x] `docs/database/schema.md`, `docs/ERD.png`
- [x] `docs/project/` — overview, implemented, tasks

## Чего пока нет

- [ ] Перевода всех пользователей с открытых паролей на bcrypt (fallback в
      `login.go` ещё есть)
- [ ] JWT и refresh-токенов
- [ ] Redis, анонимной корзины и merge при входе
- [ ] Изменения количества и очистки корзины (`PUT /cart/items`, `DELETE /cart`)
- [ ] Кнопки оформления заказа в UI и отображения `price` / `stock`
      в карточках витрины
- [ ] Чтения оформленного заказа по id/списком
- [ ] Тестов (unit/integration)
- [ ] Пагинации, фильтрации, сортировки
- [ ] OpenAPI/Swagger
- [ ] Валидации внешних ключей (`status_id`, `rule_id`) с кодом 400/404
- [ ] Graceful shutdown и структурированного логирования

## Известные ограничения

- [ ] `PUT` для несуществующего `product`/`status`/`rule`/`user_status`
      возвращает `500` вместо `404`: репозитории не оборачивают `ErrNotFound`
- [ ] Дублирующийся email в `/user/create` даёт `500`; `409` обрабатывается только
      в `/user/register`
- [ ] Списки при ошибке БД отдают текст `Error: ...`, а не `{"error": "..."}`
- [ ] `Secure` куки зависит от `COOKIE_SECURE` (по умолчанию `false`)
- [ ] Корзина и checkout доступны только авторизованным; анонимная корзина
      не реализована
- [ ] `GET /cart` показывает состав и позволяет удалить позицию через
      `DELETE /cart/items`, но не даёт изменить количество или очистить корзину —
      для этого нет `PUT /cart/items` и `DELETE /cart`
- [ ] `POST /order/checkout` требует `Idempotency-Key`; без него — `400`
- [ ] При нехватке остатка checkout отдаёт `409` с текстом
      `product N: out of stock`, остальные ошибки БД — `500`
