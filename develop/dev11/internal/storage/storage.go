package storage

import (
	"WB_LVL2/develop/dev11/internal/models"
	"WB_LVL2/develop/dev11/internal/storage/cache"
	"fmt"
	"net/http"
	"time"
)

const (
	hoursDay         = 24 * time.Hour
	daysWeek         = 7
	averageDaysMonth = 30
)

type User interface {
	GetEventForDay(userID int, date time.Time) ([]models.Event, error)
	GetEventForWeek(userID int, weekDate time.Time) ([]models.Event, error)
	GetEventForMonth(userID int, monthDate time.Time) ([]models.Event, error)
	CreateEvent(userID int, event models.Event) error
	UpdateEvent(userID int, event models.Event) error
	DeleteEvent(userID int, event models.Event) error
}

// Структура хранилища.
type Storage struct {
	*cache.Cache
}

// Функция - конструктор хранилища.
func NewStorage() *Storage {
	return &Storage{cache.NewCache()}
}

// Метод для создания эвента.
func (s *Storage) CreateEvent(userID int, event models.Event) error {
	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}
	eventId := len(user.Events) + 1
	event.Id = eventId
	user.Events[eventId] = event
	return nil
}

// Метод для обновления эвента.
func (s *Storage) UpdateEvent(userID int, event models.Event) error {
	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}
	_, ok := user.Events[event.Id]
	if !ok {
		return cache.NewErrorHandler(
			fmt.Errorf("обновляемый эвент не существует"),
			http.StatusBadRequest,
		)
	}
	user.Events[event.Id] = event
	return nil
}

// Метод для удаления эвента.
func (s *Storage) DeleteEvent(userID int, event models.Event) error {
	user, err := s.GetUser(userID)
	if err != nil {
		return err
	}
	_, ok := user.Events[event.Id]
	if !ok {
		return cache.NewErrorHandler(
			fmt.Errorf("удаляемый эвент не существует"),
			http.StatusBadRequest,
		)
	}
	delete(user.Events, event.Id)
	return nil
}

// Метод для получения всех эвентов за день.
func (s *Storage) GetEventForDay(userID int, date time.Time) ([]models.Event, error) {
	events := []models.Event{}
	user, err := s.GetUser(userID)
	if err != nil {
		return events, err
	}

	for _, v := range user.Events {
		if date.Truncate(hoursDay).Equal(v.Date.Truncate(hoursDay)) {
			events = append(events, v)
		}
	}
	return events, nil
}

// Метод для получения всех эвентов за неделю.
func (s *Storage) GetEventForWeek(userID int, weekDate time.Time) ([]models.Event, error) {
	events := []models.Event{}
	user, err := s.GetUser(userID)
	if err != nil {
		return events, err
	}
	var difTime time.Duration
	for _, v := range user.Events {
		difTime = v.Date.Sub(weekDate)
		if difTime > 0 && difTime < hoursDay*daysWeek {
			events = append(events, v)
		}
	}
	return events, nil
}

// Метод получения всех эвентов за месяц.
func (s *Storage) GetEventForMonth(userID int, monthDate time.Time) ([]models.Event, error) {
	events := []models.Event{}
	user, err := s.GetUser(userID)
	if err != nil {
		return events, err
	}
	var difTime time.Duration
	for _, v := range user.Events {
		difTime = v.Date.Sub(monthDate)
		if difTime > 0 && difTime < hoursDay*averageDaysMonth {
			events = append(events, v)
		}
	}
	return events, nil
}
