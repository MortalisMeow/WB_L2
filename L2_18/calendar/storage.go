package main

import (
	"errors"
	"sync"
	"time"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int]Event
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		events: make(map[int]Event),
		nextID: 1,
	}
}

func (s *Storage) Create(event Event) (Event, error) {
	if err := event.Validate(); err != nil {
		return Event{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.nextID
	s.nextID++
	s.events[event.ID] = event

	return event, nil
}

func (s *Storage) Update(event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; !exists {
		return errors.New("event not found")
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return errors.New("event not found")
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) GetEventsForDay(userID int, date string) ([]Event, error) {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Event
	for _, event := range s.events {
		if event.UserID == userID && event.Date == date {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForWeek(userID int, date string) ([]Event, error) {
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	weekday := targetDate.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	startOfWeek := targetDate.AddDate(0, 0, -int(weekday)+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Event
	for _, event := range s.events {
		if event.UserID != userID {
			continue
		}

		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			continue
		}

		if (eventDate.Equal(startOfWeek) || eventDate.After(startOfWeek)) &&
			eventDate.Before(endOfWeek) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsForMonth(userID int, date string) ([]Event, error) {
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	year, month, _ := targetDate.Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Event
	for _, event := range s.events {
		if event.UserID != userID {
			continue
		}

		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil {
			continue
		}

		if (eventDate.Equal(startOfMonth) || eventDate.After(startOfMonth)) &&
			eventDate.Before(endOfMonth) {
			result = append(result, event)
		}
	}
	return result, nil
}
