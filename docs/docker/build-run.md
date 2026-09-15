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
