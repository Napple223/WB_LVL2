package pattern

import "errors"

/*
Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,
а также реальные примеры использования данного примера на практике.
https://en.wikipedia.org/wiki/Facade_pattern

Фасад — это структурный паттерн проектирования,
который предоставляет простой интерфейс к сложной
системе классов, библиотеке или фреймворку.

Фасад — это простой интерфейс для работы со сложной подсистемой,
содержащей множество классов. Фасад может иметь урезанный интерфейс,
не имеющий 100% функциональности, которой можно достичь,
используя сложную подсистему напрямую.
Но он предоставляет именно те фичи, которые нужны клиенту,
и скрывает все остальные.

Применяемость:
1. Когда вам нужно представить простой или урезанный интерфейс к сложной подсистеме.
2. Когда вы хотите разложить подсистему на отдельные слои.

Преимущества:
Изолирует клиентов от компонентов сложной подсистемы.

Недостатки:
Фасад рискует стать божественным объектом, привязанным ко всем классам программы.

*/

//Мой пример:
//Допустим мы осуществялем in memory кэширование данных.
//И хотим обеспечить доступ клиента к нашему кэшу, но
//хотим спрятать от него логику его работы (проверку валидности
//данных, проверку существования, методы добавления в БД и т.д.)

//Структура валидатора.
type validator struct{}

//Функция для проверки валидности данных.
func (v *validator) dataIsValid(data string) bool {
	return true
}

//Структура хранилища данных.
type storage struct{}

//Функция проверки наличия данных в хранилище.
func (s *storage) DataAlreadyExists(data string) bool {
	return false
}

//Структура in memory cache.
type cache struct{}

//Функция для добавления данных в кэш.
func (c *cache) push(data string) error {
	return nil
}

//Функция для получения данных из кэша.
func (c *cache) pull(id int) error {
	return nil
}

//Наш фасад, предоставляющий упрощенный доступ к кэшу.
type userCache struct {
	validator *validator
	storage   *storage
	cache     *cache
}

//Функция-конструктор.
func newUserCache() *userCache {
	u := userCache{
		validator: &validator{},
		storage:   &storage{},
		cache:     &cache{},
	}
	return &u
}

//Функция для добавления данных в кэш.
func (u *userCache) push(data string) error {
	ok := u.validator.dataIsValid(data)
	if !ok {
		return errors.New("data is not valid")
	}
	ok = u.storage.DataAlreadyExists(data)
	if ok {
		return errors.New("data already exists in cache")
	}
	err := u.cache.push(data)
	if err != nil {
		return err
	}
	return nil
}

//Функция получения данных из кэша.
func (u *userCache) pull(id int) error {
	err := u.cache.pull(id)
	if err != nil {
		return err
	}
	return nil
}
