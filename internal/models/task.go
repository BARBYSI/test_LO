package models

import (
	"errors"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t *Task) Validate() error {
	if len(t.Title) == 0 {
		return errors.New("title is required")
	}
	return nil
}
