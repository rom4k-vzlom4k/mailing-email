package models

import "time"

type SentStatus string

const (
	StatusPending    SentStatus = "pending"
	StatusInProgress SentStatus = "in_progress"
	StatusDone       SentStatus = "done"
	StatusFailed     SentStatus = "failed"
)

type AddEmail struct {
	ID          int64      `json:"id"`
	ToEmail     string     `json:"to_email"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body"`
	ScheduledAt time.Time  `json:"sсheduled_at"`
	SentAt      *time.Time `json:"sent_at"`
	Status      SentStatus `json:"status"`
	Error       *string    `json:"error"`
}
