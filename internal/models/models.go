package models

import (
	"time"
)

type Task struct {
	Id          int64
	Title       string
	Description string
	Status      bool
	DueDate     time.Time
	UserId      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
