package main

import (
	"testing"
)

func TestCreateEvent(t *testing.T) {
	storage := NewStorage()

	event := Event{
		UserID: 1,
		Date:   "2024-01-01",
		Text:   "New Year",
	}

	created, err := storage.Create(event)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}
}

func TestCreateEventValidation(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name  string
		event Event
	}{
		{"empty text", Event{UserID: 1, Date: "2024-01-01", Text: ""}},
		{"invalid user", Event{UserID: 0, Date: "2024-01-01", Text: "test"}},
		{"invalid date", Event{UserID: 1, Date: "invalid", Text: "test"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.Create(tt.event)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	storage := NewStorage()

	event := Event{UserID: 1, Date: "2024-01-01", Text: "Original"}
	created, _ := storage.Create(event)

	created.Text = "Updated"
	err := storage.Update(created)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events, _ := storage.GetEventsForDay(1, "2024-01-01")
	if len(events) != 1 || events[0].Text != "Updated" {
		t.Error("event was not updated correctly")
	}
}

func TestUpdateNonExistentEvent(t *testing.T) {
	storage := NewStorage()

	err := storage.Update(Event{ID: 999, UserID: 1, Date: "2024-01-01", Text: "test"})
	if err == nil {
		t.Error("expected error for non-existent event")
	}
}

func TestDeleteEvent(t *testing.T) {
	storage := NewStorage()

	event := Event{UserID: 1, Date: "2024-01-01", Text: "test"}
	created, _ := storage.Create(event)

	err := storage.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events, _ := storage.GetEventsForDay(1, "2024-01-01")
	if len(events) != 0 {
		t.Error("event was not deleted")
	}
}

func TestDeleteNonExistentEvent(t *testing.T) {
	storage := NewStorage()

	err := storage.Delete(999)
	if err == nil {
		t.Error("expected error for non-existent event")
	}
}

func TestGetEventsForDay(t *testing.T) {
	storage := NewStorage()

	storage.Create(Event{UserID: 1, Date: "2024-01-01", Text: "Event 1"})
	storage.Create(Event{UserID: 1, Date: "2024-01-01", Text: "Event 2"})
	storage.Create(Event{UserID: 1, Date: "2024-01-02", Text: "Event 3"})
	storage.Create(Event{UserID: 2, Date: "2024-01-01", Text: "Other user"})

	events, err := storage.GetEventsForDay(1, "2024-01-01")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestGetEventsForWeek(t *testing.T) {
	storage := NewStorage()

	storage.Create(Event{UserID: 1, Date: "2024-01-01", Text: "Monday"})
	storage.Create(Event{UserID: 1, Date: "2024-01-05", Text: "Friday"})
	storage.Create(Event{UserID: 1, Date: "2024-01-07", Text: "Sunday"})
	storage.Create(Event{UserID: 1, Date: "2024-01-08", Text: "Next Monday"})

	events, err := storage.GetEventsForWeek(1, "2024-01-03")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	storage := NewStorage()

	storage.Create(Event{UserID: 1, Date: "2024-01-01", Text: "Start"})
	storage.Create(Event{UserID: 1, Date: "2024-01-15", Text: "Middle"})
	storage.Create(Event{UserID: 1, Date: "2024-01-31", Text: "End"})
	storage.Create(Event{UserID: 1, Date: "2024-02-01", Text: "Next month"})

	events, err := storage.GetEventsForMonth(1, "2024-01-15")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}
}
