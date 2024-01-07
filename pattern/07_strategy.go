package pattern

import (
	"errors"
	"fmt"
)

/*
Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы,
а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Strategy_pattern

Стратегия — это поведенческий паттерн проектирования,
который определяет семейство схожих алгоритмов и
помещает каждый из них в собственный класс,
после чего алгоритмы можно взаимозаменять прямо во
время исполнения программы.

Применяемость:
1. Когда нужно использовать разные вариации какого-то
алгоритма внутри одного объекта.
2. Когда у вас есть множество похожих классов,
отличающихся только некоторым поведением.
3. Когда вы не хотите обнажать детали реализации
алгоритмов для других классов.
4. Когда различные вариации алгоритмов реализованы
в виде развесистого условного оператора.
Каждая ветка такого оператора представляет собой вариацию алгоритма.

Преимущества:
1. Горячая замена алгоритмов на лету.
2. Изолирует код и данные алгоритмов от остальных классов.
3. Уход от наследования к делегированию.
4. Реализует принцип открытости/закрытости.

Недостатки:
1. Усложняет программу за счёт дополнительных классов.
2. Клиент должен знать, в чём состоит разница между стратегиями,
чтобы выбрать подходящую.
*/

//Мой пример:
//Мы хотим предоставить клиенту реализацию in memory
//cache. Но так как задачи у всех клиентов разные,
//один и тотже метод очистки кэша оказалался не эффективен.
//Поэтому с помощью данного паттерна, мы можем
//позволить всем выбирать необходимую стратегию очистки кэша.

// Интерфейс стратегии, определяет набор методов для всех стратегий.
type EvictionAlg interface {
	evict(c *Cache)
}

// Структура кэша.
type Cache struct {
	storage     map[int]int
	curCapacity int
	maxCapacity int
	evictionAlg EvictionAlg
}

// Функция - конструктор кэша.
func NewCache(e EvictionAlg) *Cache {
	s := make(map[int]int, 2)
	return &Cache{
		storage:     s,
		maxCapacity: 2,
		evictionAlg: e,
	}
}

// Метод, позволяющий на лету менять стратегию
// очистки кэша.
func (c *Cache) SetEvictionAlg(e EvictionAlg) {
	c.evictionAlg = e
}

// Метод очистки кэша.
func (c *Cache) evict() {
	c.evictionAlg.evict(c) //Вызываем метод очистки конкретной стратегии.
	c.curCapacity--
}

// Некоторые методы реализации кэша.
func (c *Cache) Add(k, v int) {
	if c.maxCapacity == c.curCapacity {
		c.evict()
	}
	c.curCapacity++
	c.storage[k] = v
}

func (c *Cache) Get(k int) (int, error) {
	v, ok := c.storage[k]
	if !ok {
		return -1, errors.New("key isn't exists")
	}
	return v, nil
}

// Конкретный объект стратегии.
type LRU struct{}

// Метод, имплементирующий интерфейс EvictionAlg.
func (lru *LRU) evict(c *Cache) {
	fmt.Println("Evicting by lru strtegy.")
}

// Конкретный объект стратегии.
type LFU struct{}

// Метод, имплементирующий интерфейс EvictionAlg.
func (lfu *LFU) evict(c *Cache) {
	fmt.Println("Evicting by lfu strtegy.")
}
func testStrategy() {
	lru := &LRU{}
	cache := NewCache(lru)
	cache.Add(1, 1)
	cache.Add(2, 2)
	cache.Add(3, 3) //Evicting by lru strtegy.
	lfu := &LFU{}
	cache.SetEvictionAlg(lfu) //Меняем стратегию очистки кэша.
	cache.Add(4, 4)           //Evicting by lfu strtegy.
}
