# Интерфейсы в Go — scanner

## Неявная реализация

В Go не нужно писать `implements`. Если у типа есть нужный набор методов, он автоматически подходит под интерфейс.

```go
type scanner interface {
    Scan(dest ...any) error
}
```

Такой интерфейс объявлен в `internal/product/repository.go` и `internal/user/repository.go`.

## Зачем он нужен

И `*sql.Row`, и `*sql.Rows` умеют `Scan`, поэтому обе структуры подходят под `scanner`. Функция разбора строки пишется один раз:

```go
func scanProduct(s scanner) (*Product, error) {
    p := &Product{}
    // ...
    if err := s.Scan(&p.ID, &p.Name, &p.Description); err != nil {
        return nil, err
    }
    return p, nil
}
```

Использование:

```go
// одна строка
p, err := scanProduct(db.QueryRow(query, id))

// много строк
for rows.Next() {
    p, err := scanProduct(rows)
    // ...
}
```

## Почему `*sql.Row`, а не `sql.Row`

Метод `Scan` объявлен с указательным получателем:

```go
func (r *Row) Scan(dest ...any) error
```

Поэтому интерфейс удовлетворяет именно указатель `*sql.Row`. Значение `sql.Row` в интерфейс не попадёт.

## `...any` — вариативный аргумент

```go
Scan(dest ...any) error
```

- `any` — псевдоним `interface{}` («любой тип»);
- `...any` — можно передать сколько угодно аргументов;
- внутри `Scan` они рассматриваются как `[]any`.

Аналогия в JavaScript: `function scan(...dest) {}` — rest-параметры.

## Go vs TypeScript

| Go                           | TypeScript                                 |
|------------------------------|--------------------------------------------|
| интерфейс описывает методы   | интерфейс описывает поля и методы          |
| реализация неявная           | реализация через `implements`              |
| проверка на этапе компиляции | типы стираются, проверка в рантайме        |

---

## Чек-лист

- [ ] интерфейс содержит только нужные методы
- [ ] `*sql.Row` и `*sql.Rows` подходят под `scanner`
- [ ] для указательных методов в интерфейс кладём указатель (`*T`, а не `T`)
- [ ] `...any` = «любое количество аргументов любого типа»
