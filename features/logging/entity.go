package logging

import (
	"time"

	"github.com/labstack/echo/v4"
)

type LogData struct {
	LogID             string    `json:"log_id"`
	LogType           string    `json:"log_type"`
	Message           string    `json:"message"`
	CreatedAt         time.Time `json:"created_at"`
	InitiatedBy       string    `json:"created_by"`
	InitiatedByUserID uint      `json:"created_by_user_id"`
}

type LoggingHandlerInterface interface {
	AddLog() echo.HandlerFunc
	ViewLog() echo.HandlerFunc
	ViewOneLog() echo.HandlerFunc
	// UpdateLog() echo.HandlerFunc
	DeleteLog() echo.HandlerFunc
}

type LoggingServiceInterface interface {
	AddLog(newData LogData) (*LogData, error)
	ViewLog() ([]LogData, error)
	ViewOneLog(logID string) (*LogData, error)
	// UpdateLog(logID uint, newData LogData) (bool, error)
	DeleteLog(logID string) (bool, error)
}

type LoggingDataInterface interface {
	AddLog(newData LogData) (*LogData, error)
	ViewLog() ([]LogData, error)
	ViewOneLog(logID string) (*LogData, error)
	// UpdateLog(logID uint, newData LogData) (bool, error)
	DeleteLog(logID string) (bool, error)
}
