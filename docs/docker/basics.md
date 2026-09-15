# Docker — основы

## VM vs Docker

```
ВМ (VirtualBox):
┌─────────────────────────────┐
│     Host OS (Linux/Win)     │
├─────────────────────────────┤
│   Hypervisor (VirtualBox)   │
├─────────────────────────────┤
│  Guest OS Kernel (Linux)    │  ← полный Linux внутри
│  Libs, Bins, Apps           │
└─────────────────────────────┘

Docker (Контейнеризация):
┌─────────────────────────────┐
│     Host OS Kernel          │  ← используем одно ядро
├─────────────────────────────┤
│  Container1: Libs + App     │
│  Container2: Libs + App     │
│  Container3: Libs + App     │
└─────────────────────────────┘
```

Docker использует **namespace** и **cgroups** хоста:

### Namespace (изоляция)

- **PID namespace** — процессы контейнера видят только себя
- **Network namespace** — свой virtual network interface
- **Filesystem namespace** — свой корневой `/` (через Union FS)
- **IPC namespace** — своя межпроцессная коммуникация

### Cgroups (ограничение ресурсов)

- Лимиты памяти, CPU, I/O
- Контейнер не может больше взять, чем выделили

**Пример:** Твой процесс Go видит себя как PID 1 внутри контейнера, но на хосте это какой-то другой PID.

## Слои образа и кеширование

Когда ты делаешь `docker build`, это не просто копирование файлов:

```dockerfile
FROM golang:1.21-alpine        # Слой 1: 500MB базовый образ
WORKDIR /app                   # Слой 2: смена директории (мета)
COPY main.go .                 # Слой 3: копия файла
RUN go build -o app main.go    # Слой 4: результат сборки
```

**Каждая инструкция = новый слой** (как коммиты в git). Docker кеширует слои:

```bash
# Первый раз: строит все 4 слоя
docker build -t my-go-app .

# Второй раз (если main.go не изменился): использует кеш
docker build -t my-go-app .
```

Если изменишь `main.go` — пересчитается только слой 3 и 4.
