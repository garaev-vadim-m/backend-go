# Go Pointers — Шпаргалка

## Основное
```
& = "дай адрес"       (address of)
* = "дай значение"    (dereference)
```

---

## Три вещи
```go
age := 25              // переменная (значение)
ptr := &age            // указатель (адрес)
fmt.Println(*ptr)      // разыменование (вытаскиваем значение)
```

**Результат:**
```
age    = 25                    ← значение
&age   = 0xc0000180d0          ← адрес в памяти
*ptr   = 25                    ← значение по адресу
```

---

## Изменение через указатель
```go
*ptr = 30
fmt.Println(age)  // 30 ← age изменилась!
```

---

## nil — "пусто"
```go
var ptr *string = nil   // nil указатель (как NULL в БД)

if ptr == nil {
	fmt.Println("пусто")
}
```

---

## В структурах
```go
type User struct {
	ID       int     // обязательно
	Name     string  // обязательно
	Female   *string // может быть nil (опционально)
	StatusID *int    // может быть nil (опционально)
}

// Создание
female := "Female"
user := User{
	ID:     1,
	Name:   "John",
	Female: &female,   // адрес переменной
	StatusID: nil,     // пусто
}

// Использование
if user.Female != nil {
	fmt.Println(*user.Female)  // Female
}
```

---

## Аналогия с почтой
```
age = 25              ← письмо (значение)
&age = адрес ящика   ← где лежит письмо
*ptr = 25             ← открыли ящик, прочитали
```

---

## JS аналогия
```javascript
// JS: всё автоматически по ссылке
let user = { name: "John" };
let copy = user;
copy.name = "Jane";
console.log(user.name);  // "Jane" (один объект)
```

```go
// Go: явно указываем
user := User{Name: "John"}
ptr := &user          // адрес
(*ptr).Name = "Jane"
fmt.Println(user.Name)  // "Jane" (один объект)
```

---

## Чек-лист
- [ ] `&variable` — получить адрес
- [ ] `*pointer` — разыменовать (вытащить значение)
- [ ] `*Type` в структуре — поле может быть nil
- [ ] `if ptr != nil` — проверка перед использованием
