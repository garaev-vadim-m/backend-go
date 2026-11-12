# Docker — сборка и запуск

## Сборка образа

```bash
docker build -t my-go-app .
```

- `-t my-go-app` — даем имя образу
- `.` — Dockerfile в текущей директории

## Запуск контейнера

```bash
docker run -p 8080:8080 my-go-app
```

- `-p 8080:8080` — маппим порт контейнера (8080) на порт хоста (8080)

### Синхронный vs фоновый режим

Без флагов терминал блокируется и выводит логи контейнера. Ctrl+C останавливает контейнер.

С флагом `-d` (detached):

```bash
docker run -d -p 8080:8080 my-go-app
```

Контейнер работает в фоне. Закрытие терминала не останавливает его.

## Управление контейнерами

```bash
docker ps                    # список работающих
docker ps -a                 # все (включая остановленные)
docker logs -f my-backend    # логи в реальном времени
docker stop my-backend       # остановить
docker rm my-backend         # удалить
docker start my-backend      # запустить остановленный
```

## Жизненный цикл

```
docker build      docker run         docker stop       docker rm
    ↓                ↓                  ↓                ↓
[Образ]  ───→  [Работающий]  ───→  [Остановлен]  ───→ [Удален]
               контейнер            (на диске)       (с диска)
```

Остановленный контейнер остается на диске (`docker ps -a` его покажет).

## Практический тест

```bash
# Запуск в фоне
docker run -d -p 8080:8080 --name my-backend my-go-app

# Проверка
docker ps

# Логи
docker logs -f my-backend

# Остановка
docker stop my-backend

# Проверяем
docker ps         # пусто
docker ps -a      # есть (остановлен)

# Перезапуск
docker start my-backend
```

## Многоэтапная сборка (multi-stage)

В проекте Dockerfile состоит из двух этапов:

```dockerfile
# Этап 1: сборка
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

# Этап 2: финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
```

Схема:

```
builder (golang:1.22-alpine)          final (alpine:latest)
┌────────────────────────┐            ┌──────────────────────┐
│ go.mod, go.sum         │            │ только бинарник app  │
│ исходники              │  ────────► │ EXPOSE 8080          │
│ go build -o app        │  COPY app  │ CMD ["./app"]        │
└────────────────────────┘            └──────────────────────┘
```

Зачем так:

- в финальный образ не попадают Go-компилятор и исходники — образ меньше;
- `COPY go.mod go.sum` и `go mod download` вынесены в отдельные слои, поэтому зависимости кешируются и не перекачиваются при изменении кода;
- `CGO_ENABLED=0` даёт статический бинарник, который работает в `alpine` без дополнительных библиотек.

## .dockerignore

```
.git
*.md
docker-compose.yml
Dockerfile
```

Файлы из `.dockerignore` не попадают в build context: сборка быстрее, а лишнее не копируется в образ.

> Версия Go в Dockerfile (`golang:1.22-alpine`) — это версия сборки. Она может отличаться от версии в `go.mod`; при обновлении Go меняйте оба места.
