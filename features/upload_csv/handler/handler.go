package handler

import (
	"fmt"
	"net/http"

	"gain-v2/features/upload_csv"
	"gain-v2/helper"

	"github.com/labstack/echo/v4"
)

type UploadCSVHandler struct {
	s upload_csv.UploadCSVServiceInterface
}

func NewHandler(service upload_csv.UploadCSVServiceInterface) upload_csv.UploadCSVHandlerInterface {
	return &UploadCSVHandler{
		s: service,
	}
}

func (uh *UploadCSVHandler) CreateToken() echo.HandlerFunc {
	return func(c echo.Context) error {

		var input = new(ReqCreateToken)

		if err := c.Bind(input); err != nil {
			fmt.Println("Errornya disini:", err)
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "invalid log input", nil))
		}

		if input.Action != "apply" {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "cannot generate token", nil))
		}

		resToken, err := uh.s.CreateToken(c)

		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, err.Error(), nil))

		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "success generate token", resToken))

	}
}

func (uh *UploadCSVHandler) UploadCSV() echo.HandlerFunc {
	return func(c echo.Context) error {

		const maxFileSize = 100 * 1024 * 1024

		csvFile, err := c.FormFile("file")
		if err != nil {

			fmt.Println("err: ", err.Error())

			return c.JSON(http.StatusBadRequest, helper.FormatResponse(false, "invalid file upload", nil))
		}

		if csvFile.Size > maxFileSize {
			return c.JSON(http.StatusRequestEntityTooLarge, helper.FormatResponse(false, "file size exceeds 100MB", nil))
		}

		resData, err := uh.s.UploadCSV(csvFile)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.FormatResponse(false, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, helper.FormatResponse(true, "success", resData))
	}
}
