package handler

import (
	"fmt"
	"gain-v2/features/logging"
	"gain-v2/helper"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type LoggingHandler struct {
	s   logging.LoggingServiceInterface
	jwt helper.JWTInterface
}

func NewHandler(service logging.LoggingServiceInterface, jwt helper.JWTInterface) logging.LoggingHandlerInterface {
	return &LoggingHandler{
		s:   service,
		jwt: jwt,
	}
}

func (lh *LoggingHandler) AddLog() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoggingInput)

		if err := c.Bind(input); err != nil {
			fmt.Println("Errornya disini:", err)
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "invalid log input", nil))
		}

		var serviceInput = new(logging.LogData)
		serviceInput.LogID = input.LogID
		serviceInput.LogType = input.LogType
		serviceInput.Message = input.Message
		serviceInput.CreatedAt = input.CreatedAt
		serviceInput.InitiatedBy = input.InitiatedBy
		serviceInput.InitiatedByUserID = input.InitiatedByUserID

		result, err := lh.s.AddLog(*serviceInput)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		var response = new(LoggingResponse)
		response.LogID = result.LogID
		response.LogType = result.LogType
		response.Message = result.Message
		response.CreatedAt = result.CreatedAt
		response.InitiatedBy = result.InitiatedBy
		response.InitiatedByUserID = result.InitiatedByUserID

		return c.JSON(http.StatusCreated, helper.FormatResponse(true, "log added", response))
	}
}
func (lh *LoggingHandler) ViewLog() echo.HandlerFunc {
	return func(c echo.Context) error {

		result, err := lh.s.ViewLog()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "view all logs", result))
	}
}
func (lh *LoggingHandler) ViewOneLog() echo.HandlerFunc {
	return func(c echo.Context) error {
		param := c.Param("log_id")

		result, err := lh.s.ViewOneLog(param)

		if err != nil {
			if strings.Contains(err.Error(), "[404 Not Found]") {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, fmt.Sprintf("error finding log with id %v : not found", param), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "view log by id", result))
	}
}
func (lh *LoggingHandler) DeleteLog() echo.HandlerFunc {
	return func(c echo.Context) error {
		param := c.Param("log_id")

		_, err := lh.s.DeleteLog(param)

		if err != nil {
			if strings.Contains(err.Error(), "[404 Not Found]") {
				return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, fmt.Sprintf("error deleting log with id %v : not found", param), nil))
			}
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusNoContent, helper.FormatResponse(true, "success delete log", nil))
	}
}
