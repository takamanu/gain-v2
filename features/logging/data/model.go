package data

import (
	"time"

	"github.com/google/uuid"
)

// const TableName = "loggings"

type LogData struct {
	LogID             uuid.UUID `json:"log_id"`
	LogType           string    `json:"log_type"`
	Message           string    `json:"message"`
	CreatedAt         time.Time `json:"created_at"`
	InitiatedBy       string    `json:"created_by"`
	InitiatedByUserID uint      `json:"created_by_user_id"`
}

func (LogData) IndexName() string {
	return "loggings"
}
