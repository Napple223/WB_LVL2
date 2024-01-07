Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Вывод: error

Аналогично listing03.
Переменная err объявлена как интерфейс ошибки err.
Во время выполнения выполнения функции test, которая
имплементирует интерфейс error, мы присвоили интерфейсу
тип *customError. Т.е. в данной ситуации мы сравниваем
nil не для указателя, а для интерфейса, которому мы
присвоили тип *customError.

 Интерфейс считается nil только тогда, когда его значение
 и тип равны nil, что в нашем случае не выполняется.

```