package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (h *Handler) decodeEventBodyJSON(r *http.Request) (*eventInput, error) {
	input := &eventInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func getParamsInput(url *url.URL) (int, time.Time, error) {
	userId := url.Query().Get("user_id")
	date := url.Query().Get("date")

	inpDate := Date{}
	err := inpDate.UnmarshalJSON([]byte(date))
	if err != nil {
		return 0, time.Time{}, err
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return 0, time.Time{}, err
	}

	return id, time.Time(inpDate), nil
}
