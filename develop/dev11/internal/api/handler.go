package api

import (
	"WB_LVL2/develop/dev11/internal/models"
	"WB_LVL2/develop/dev11/internal/storage"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	storage.Storage
}

func NewHandler(u *storage.Storage) *Handler {
	return &Handler{*u}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", h.createEvent)
	mux.HandleFunc("/update_event", h.updateEvent)
	mux.HandleFunc("/delete_event", h.deleteEvent)
	mux.HandleFunc("/events_for_day", h.getEventsForDay)
	mux.HandleFunc("/events_for_week", h.getEventsForWeek)
	mux.HandleFunc("/events_for_month", h.getEventsForMonth)
	handler := Log(mux)
	return handler
}

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	if httpMethodErrorCheck(w, http.MethodPost, r.Method) {
		return
	}

	input, err := h.decodeEventBodyJSON(r)
	if err != nil {
		h.httpErrorResponse(w, http.StatusBadRequest, err.Error())
		log.Println(err.Error())
		return
	}

	err = h.Storage.CreateEvent(input.UserID, models.Event{
		Id:          input.EventID,
		Name:        input.Name,
		Description: input.Description,
		Date:        time.Time(input.Date),
	})
	if err != nil {
		h.httpErrorResponse(w, http.StatusServiceUnavailable, err.Error())
		log.Println(err.Error())
		return
	}

	h.httpSuccessEventActionResponse(w, http.StatusOK, "Эвент успешно создан.")
}

func (h *Handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	if httpMethodErrorCheck(w, http.MethodPost, r.Method) {
		return
	}

	input, err := h.decodeEventBodyJSON(r)
	if err != nil {
		h.httpErrorResponse(w, http.StatusBadRequest, err.Error())
		log.Println(err.Error())
		return
	}

	err = h.Storage.UpdateEvent(input.UserID, models.Event{
		Id:          input.EventID,
		Name:        input.Name,
		Description: input.Description,
		Date:        time.Time(input.Date),
	})

	if err != nil {
		h.httpErrorResponse(w, http.StatusServiceUnavailable, err.Error())
		log.Println(err.Error())
		return
	}

	h.httpSuccessEventActionResponse(w, http.StatusOK, "Эвент успешно обновлен.")
}

func (h *Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	if httpMethodErrorCheck(w, http.MethodPost, r.Method) {
		return
	}

	input, err := h.decodeEventBodyJSON(r)
	if err != nil {
		h.httpErrorResponse(w, http.StatusBadRequest, err.Error())
		log.Println(err.Error())
		return
	}

	err = h.Storage.DeleteEvent(input.UserID, models.Event{
		Id:          input.EventID,
		Name:        input.Name,
		Description: input.Description,
		Date:        time.Time(input.Date),
	})

	if err != nil {
		h.httpErrorResponse(w, http.StatusServiceUnavailable, err.Error())
		log.Println(err.Error())
		return
	}

	h.httpSuccessEventActionResponse(w, http.StatusOK, "Эвент успешно удален.")
}

func (h *Handler) getEventsForDay(w http.ResponseWriter, r *http.Request) {
	if httpMethodErrorCheck(w, http.MethodGet, r.Method) {
		return
	}

	userID, date, err := getParamsInput(r.URL)
	if err != nil {
		h.httpErrorResponse(w, http.StatusBadRequest, err.Error())
		log.Println(err.Error())
		return
	}
	events, err := h.Storage.GetEventForDay(userID, date)
	if err != nil {
		h.httpErrorResponse(w, http.StatusServiceUnavailable, err.Error())
		log.Println(err.Error())
		return
	}
	h.httpEventsResponse(w, http.StatusOK, events)
}

func (h *Handler) getEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if httpMethodErrorCheck(w, http.MethodGet, r.Method) {
		return
	}

	userID, date, err := getParamsInput(r.URL)
	if err != nil {
		h.httpErrorResponse(w, http.StatusBadRequest, err.Error())
		log.Println(err.Error())
		return
	}
	events, err := h.Storage.GetEventForWeek(userID, date)
	if err != nil {
		h.httpErrorResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	h.httpEventsResponse(w, http.StatusOK, events)
}

func (h *Handler) getEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if httpMethodErrorCheck(w, http.MethodGet, r.Method) {
		return
	}

	userID, date, err := getParamsInput(r.URL)
	if err != nil {
		h.httpErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	events, err := h.Storage.GetEventForMonth(userID, date)
	if err != nil {
		h.httpErrorResponse(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	h.httpEventsResponse(w, http.StatusOK, events)
}
