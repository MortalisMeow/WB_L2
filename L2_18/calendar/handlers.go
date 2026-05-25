package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Server struct {
	storage *Storage
}

func NewServer(storage *Storage) *Server {
	return &Server{storage: storage}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusBadRequest, "only POST method allowed")
		return
	}

	var event Event

	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusBadRequest, "invalid form data")
			return
		}

		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		event = Event{
			UserID: userID,
			Date:   r.FormValue("date"),
			Text:   r.FormValue("text"),
		}
	}

	created, err := s.storage.Create(event)
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"result": created,
	})
}

func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusBadRequest, "only POST method allowed")
		return
	}

	var event Event

	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusBadRequest, "invalid form data")
			return
		}

		id, _ := strconv.Atoi(r.FormValue("id"))
		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		event = Event{
			ID:     id,
			UserID: userID,
			Date:   r.FormValue("date"),
			Text:   r.FormValue("text"),
		}
	}

	if err := s.storage.Update(event); err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"result": "event updated",
	})
}

func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusBadRequest, "only POST method allowed")
		return
	}

	var id int

	if r.Header.Get("Content-Type") == "application/json" {
		var requestData struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			respondError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		id = requestData.ID
	} else {
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusBadRequest, "invalid form data")
			return
		}
		id, _ = strconv.Atoi(r.FormValue("id"))
	}

	if err := s.storage.Delete(id); err != nil {
		respondError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"result": "event deleted",
	})
}

func (s *Server) EventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusBadRequest, "only GET method allowed")
		return
	}

	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	date := r.URL.Query().Get("date")

	events, err := s.storage.GetEventsForDay(userID, date)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"result": events,
	})
}

func (s *Server) EventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusBadRequest, "only GET method allowed")
		return
	}

	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	date := r.URL.Query().Get("date")

	events, err := s.storage.GetEventsForWeek(userID, date)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"result": events,
	})
}

func (s *Server) EventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusBadRequest, "only GET method allowed")
		return
	}

	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	date := r.URL.Query().Get("date")

	events, err := s.storage.GetEventsForMonth(userID, date)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"result": events,
	})
}

func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", s.CreateEventHandler)
	mux.HandleFunc("/update_event", s.UpdateEventHandler)
	mux.HandleFunc("/delete_event", s.DeleteEventHandler)
	mux.HandleFunc("/events_for_day", s.EventsForDayHandler)
	mux.HandleFunc("/events_for_week", s.EventsForWeekHandler)
	mux.HandleFunc("/events_for_month", s.EventsForMonthHandler)

	return mux
}
