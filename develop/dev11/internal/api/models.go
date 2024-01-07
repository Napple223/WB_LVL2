package api

import (
	"WB_LVL2/develop/dev11/internal/models"
	"encoding/json"
	"strings"
	"time"
)

type Date time.Time

type eventInput struct {
	UserID      int    `json:"user_id"`
	EventID     int    `json:"event_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        Date   `json:"date"`
}

type successEventOutput struct {
	Result string `json:"result"`
}

type eventsOutput struct {
	Result []models.Event `json:"result"`
}

type errorOutput struct {
	Error string `json:"error"`
}

func newSuccessEventOutput(result string) successEventOutput {
	return successEventOutput{Result: result}
}

func newEventsOutput(result []models.Event) eventsOutput {
	return eventsOutput{Result: result}
}

func newErrorOutput(message string) errorOutput {
	return errorOutput{Error: message}
}

func (i *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*i = Date(t)
	return nil
}

func (i Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(i))
}
