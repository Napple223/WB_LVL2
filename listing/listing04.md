Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}
```

Ответ:
```
Вывод: значения от 0 до 9, а потом deadlock.

Дедлок происходит потому что горутина не закрыла канал после отправки всех значений, а main горутина заблокирована циклом.

```