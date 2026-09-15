# fmt.Sprintf — форматирование строк

## Аналогия с JavaScript

```javascript
const name = "John";
const age = 30;
const str = `My name is ${name} and I'm ${age} years old`;
```

В Go:

```go
name := "John"
age := 30
str := fmt.Sprintf("My name is %s and I'm %d years old", name, age)
```

## Плейсхолдеры

| Код  | Тип                      |
|------|--------------------------|
| `%s` | строка (string)          |
| `%d` | целое число (integer)    |
| `%f` | число с плавающей точкой |
| `%v` | любой тип                |
