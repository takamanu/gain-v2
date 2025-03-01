package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gain-v2/features/upload_csv"
	"gain-v2/helper"
	email "gain-v2/helper/email"
	encrypt "gain-v2/helper/encrypt"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

type UploadCSVService struct {
	d     upload_csv.UploadCSVDataInterface
	j     helper.JWTInterface
	e     encrypt.HashInterface
	email email.EmailInterface
}

func NewService(data upload_csv.UploadCSVDataInterface, jwt helper.JWTInterface, email email.EmailInterface, encrypt encrypt.HashInterface) upload_csv.UploadCSVServiceInterface {
	return &UploadCSVService{
		d:     data,
		j:     jwt,
		email: email,
		e:     encrypt,
	}
}

func (us *UploadCSVService) CreateToken(c echo.Context) (*upload_csv.ResCreateToken, error) {

	jwtToken := us.j.GenerateUploadJWTToken(c.Request().Header.Get("Signature"))

	return &upload_csv.ResCreateToken{
		Token:      jwtToken,
		Expiration: time.Now().Add(1 * time.Hour),
	}, nil
}
func (us *UploadCSVService) UploadCSV(csvFile *multipart.FileHeader) (*upload_csv.ResGainUpload, error) {
	// Preload data

	timeStart := time.Now()
	fmt.Println("time start: ", timeStart)

	regions, err := us.d.GetRegions()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}
	sectors, err := us.d.GetSectors()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}
	// Read CSV file
	dataObjects, err := us.ReadCSVFile(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	// Insert data into database
	result, err := us.InsertDataToDatabase(dataObjects, regions, sectors)
	if err != nil {
		return nil, fmt.Errorf("failed to insert data: %w", err)
	}

	fmt.Println("Upload complete!")
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	fmt.Println("Duration: ", duration)
	return result, nil
}

// Reads CSV file and converts it into an array of GainCSVData objects
func (us *UploadCSVService) ReadCSVFile(csvFile *multipart.FileHeader) ([]upload_csv.GainCSVData, error) {
	file, err := csvFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	var (
		dataObjects []upload_csv.GainCSVData
		errorCount  int32
		wg          sync.WaitGroup
		mu          sync.Mutex
	)

	// Process rows in parallel
	for i, row := range rows {
		if i == 0 || i >= 150000 {
			continue
		}
		wg.Add(1)
		go func(row []string) {
			defer wg.Done()
			data, err := us.ExtractDataToGainObject(row, &errorCount)
			if err == nil {
				mu.Lock()
				dataObjects = append(dataObjects, data)
				mu.Unlock()
			}
		}(row)
	}

	wg.Wait()
	return dataObjects, nil
}

func (us *UploadCSVService) ExtractDataToGainObject(row []string, errorCount *int32) (upload_csv.GainCSVData, error) {
	if len(row) < 12 || row[11] == "" {
		return upload_csv.GainCSVData{}, errors.New("invalid row")
	}

	year, err := strconv.Atoi(row[2])
	if err != nil {
		atomic.AddInt32(errorCount, 1)
		return upload_csv.GainCSVData{}, err
	}

	value := strings.TrimSpace(row[11])
	if value == "" || value == "-" {
		return upload_csv.GainCSVData{}, errors.New("empty value")
	}

	// Convert negative values in parentheses
	if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
		value = "-" + value[1:len(value)-1]
	}

	value = strings.ReplaceAll(value, ",", "")
	dv, err := strconv.ParseFloat(value, 64)
	if err != nil {
		atomic.AddInt32(errorCount, 1)
		return upload_csv.GainCSVData{}, err
	}

	return upload_csv.GainCSVData{
		Zone:              row[0],
		AreaName:          row[1],
		TimePeriod:        row[2],
		Source:            row[3],
		Nb:                row[4],
		Sector:            helper.RemoveNumberAndPeriods(row[5]),
		Subsector:         helper.RemoveNumberAndPeriods(row[6]),
		IndicatorMetadata: row[7],
		Unit:              row[8],
		StartYear:         year - 1,
		EndYear:           year,
		DataValue:         dv,
	}, nil
}

func (us *UploadCSVService) InsertDataToDatabase(dataObjects []upload_csv.GainCSVData, regions []upload_csv.Region, sectors []upload_csv.Sector) (*upload_csv.ResGainUpload, error) {
	var (
		dataChan   = make(chan upload_csv.GainCSVData, len(dataObjects)/2)
		errorCount int32
		wg         sync.WaitGroup
	)

	// Send data to channel
	go func() {
		for _, data := range dataObjects {
			dataChan <- data
		}
		close(dataChan)
	}()

	// Start workers
	// numWorkers := runtime.NumCPU() * 2
	wg.Add(100) // Add to wait group before starting workers

	for i := 0; i < 100; i++ {
		go us.processCSVWorker(dataChan, &errorCount, regions, sectors, &wg)
	}

	wg.Wait() // Wait for all workers to finish

	return &upload_csv.ResGainUpload{
		StatusUpload:        true,
		SuccessInsertedData: uint(len(dataObjects) - int(errorCount)),
		FailedInsertData:    uint(errorCount),
		Errors:              nil,
	}, nil
}

func (us *UploadCSVService) processCSVWorker(dataChan chan upload_csv.GainCSVData, errorCount *int32, regions []upload_csv.Region, sectors []upload_csv.Sector, wg *sync.WaitGroup) {
	defer wg.Done()

	i := 0
	for data := range dataChan {
		regionID := findRegionID(data.AreaName, regions)
		sectorID, subsectorID, indicatorID := findSectorDetails(data, sectors)
		if regionID == 0 || sectorID == 0 || subsectorID == 0 || indicatorID == 0 {
			atomic.AddInt32(errorCount, 1)
			continue
		}

		data.RegionID = regionID
		data.SectorID = sectorID
		data.SubsectorID = subsectorID
		data.IndicatorID = indicatorID
		i++

		fmt.Println("[Debug] Ini data ke: ", i)

		if err := us.d.InsertGainData(data); err != nil {
			atomic.AddInt32(errorCount, 1)
		}
	}
}

func findRegionID(name string, regions []upload_csv.Region) int {
	for _, r := range regions {
		if r.Name == name {
			return r.ID
		}
	}
	return 0
}

func findSectorDetails(data upload_csv.GainCSVData, sectors []upload_csv.Sector) (int, int, int) {
	for _, s := range sectors {
		if s.Name == data.Sector {
			for _, sub := range s.Subsectors {
				if sub.Name == data.Subsector {
					for _, ind := range sub.Indicators {
						if ind.Name == data.IndicatorMetadata {
							return s.ID, sub.ID, ind.ID
						}
					}
				}
			}
		}
	}
	return 0, 0, 0
}
