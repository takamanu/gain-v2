package upload_csv

import (
	"mime/multipart"
	"time"

	"github.com/labstack/echo/v4"
)

type GainCSVData struct {
	Zone              string
	AreaName          string
	TimePeriod        string
	Source            string
	Nb                string
	Sector            string
	Subsector         string
	IndicatorMetadata string
	Unit              string
	StartYear         int
	EndYear           int
	DataValue         float64
	RegionID          int
	SectorID          int
	SubsectorID       int
	IndicatorID       int
}

type ResGainUpload struct {
	StatusUpload        bool        `json:"status_upload"`
	SuccessInsertedData uint        `json:"success_inserted_data"`
	FailedInsertData    uint        `json:"failed_insert_data"`
	Errors              interface{} `json:"errors"`
}

type ResCreateToken struct {
	Token      interface{} `json:"token"`
	Expiration time.Time   `json:"expiration"`
}

type Region struct {
	ID   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

type Indicator struct {
	ID          int    `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name"`
	SubsectorID int    `gorm:"column:subsectorId"`
}

type Subsector struct {
	ID         int         `gorm:"column:id;primaryKey"`
	Name       string      `gorm:"column:name"`
	SectorID   int         `gorm:"column:sectorId"`
	Indicators []Indicator `gorm:"foreignKey:SubsectorID;references:ID"`
}

type Sector struct {
	ID         int         `gorm:"column:id;primaryKey"`
	Name       string      `gorm:"column:name"`
	Subsectors []Subsector `gorm:"foreignKey:SectorID;references:ID"`
}

func (Subsector) TableName() string {
	return "Subsector"
}

func (Indicator) TableName() string {
	return "Indicator"
}

type UploadCSVHandlerInterface interface {
	CreateToken() echo.HandlerFunc
	UploadCSV() echo.HandlerFunc
}

type UploadCSVServiceInterface interface {
	CreateToken(c echo.Context) (*ResCreateToken, error)
	UploadCSV(csvFile *multipart.FileHeader) (*ResGainUpload, error)
}

type UploadCSVDataInterface interface {
	GetRegionID(areaName string) (int, error)
	GetSectorID(sectorName string) (int, error)
	GetSubsectorID(subsectorName string, sectorID int) (int, error)
	GetIndicatorID(indicatorName string, subsectorID int) (int, error)
	InsertOrUpdateDatum(dataValue float64, startYear, endYear, regionID, indicatorID int) error
	PreloadIDs() (map[string]int, map[string]int, map[string]map[string]int, map[string]map[string]int, error)
	BulkInsertGainData(data []GainCSVData) error
	GetRegions() ([]Region, error)
	GetSectors() ([]Sector, error)
	InsertGainData(data GainCSVData) error
}
