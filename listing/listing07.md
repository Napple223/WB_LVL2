Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:
```
Вывод: в случайном порядке выводятся значения, переданные в качестве аргументов функции asChan. Потом бесконечные нули.

Функция asChan пишет переданные аргументы в канал, а потом закрывает его.
Функция merge бесконечно читает переданные каналы.
В функции merge не предусмотрена проверка булевого значения, показывающего открыт канал или нет. Поэтому после того как канал закрывается, он бесконечно передает нулевое значение типа канала (в нашем случае 0).

Важно отметить, что читать из закрытого канала можно, а если писать в закрытый канал, то приложение запаникует.

```