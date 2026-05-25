package main

import (
	"errors"
	"time"
)

type Event struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Text   string `json:"text"`
}

func (e *Event) Validate() error {
	if e.UserID <= 0 {
		return errors.New("user_id must be positive")
	}
	if e.Text == "" {
		return errors.New("event text cannot be empty")
	}

	_, err := time.Parse("2006-01-02", e.Date)
	if err != nil {
		return errors.New("invalid date format, use YYYY-MM-DD")
	}

	return nil
}
