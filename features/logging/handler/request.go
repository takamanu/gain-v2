package handler

import "time"

type LoggingInput struct {
	LogID             string    `json:"log_id"`
	LogType           string    `json:"log_type"`
	Message           string    `json:"message"`
	CreatedAt         time.Time `json:"created_at"`
	InitiatedBy       string    `json:"created_by"`
	InitiatedByUserID uint      `json:"created_by_user_id"`
}
