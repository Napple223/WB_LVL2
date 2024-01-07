package user

import "WB_LVL2/develop/dev11/internal/models"

const initialEventsMapSize = 10

type User struct {
	Id     int                  `json:"id"`
	Events map[int]models.Event `json:"events"`
}

func NewUser(id int) *User {
	return &User{
		Id:     id,
		Events: make(map[int]models.Event, initialEventsMapSize),
	}
}
