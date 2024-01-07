package cache

//Пакет, реализующий in memory cache.

import (
	"WB_LVL2/develop/dev11/internal/models"
	"WB_LVL2/develop/dev11/internal/storage/cache/user"
	"fmt"
	"net/http"
	"sync"
)

const initialMapSize = 10

// Структура кэша.
type Cache struct {
	Mutex sync.RWMutex
	Data  map[int]*user.User
}

// Функция - конструктор для создания кэша.
func NewCache() *Cache {
	var cache Cache
	cache.Data = make(map[int]*user.User, initialMapSize)
	return &cache
}

// Метод для создания нового юзера.
func (c *Cache) AddUser(id int) {
	c.Mutex.Lock()
	c.Data[id] = user.NewUser(id)
	c.Mutex.Unlock()
}

// Метод для получения юзера.
func (c *Cache) GetUser(id int) (*user.User, error) {
	c.Mutex.RLock()
	user, ok := c.Data[id]
	c.Mutex.RUnlock()
	if ok {
		return user, nil
	}
	return nil, NewErrorHandler(
		fmt.Errorf("пользователя с id %d не существует", id),
		http.StatusBadRequest,
	)
}

// Метод для создания нового эвента.
func (c *Cache) AddUserEvent(id int, event models.Event) error {
	c.Mutex.Lock()
	user, err := c.GetUser(id)
	if err != nil {
		return err
	}
	user.Events[event.Id] = event
	return nil
}
