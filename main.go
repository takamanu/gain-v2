package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataRow struct {
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
}

const DSN = "host=46.250.239.194 user=postgres password=B4r4kadut1234 dbname=gain port=5432 sslmode=disable TimeZone=Asia/Jakarta"

func cleanData(input string) string {
	// Define a regex to remove leading numbers and periods (e.g., "5. " or "10. ")
	re := regexp.MustCompile(`^\d+\.\s*`)
	return re.ReplaceAllString(input, "")
}

func main() {

	dsn := DSN
	// Open connection to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Err", err)
		return
	}
	fmt.Println("instance: ", db)

	// Configure connection pooling
	sqlDB, err := db.DB() // Get the underlying sql.DB instance
	if err != nil {
		fmt.Println("Error getting underlying DB:", err)
		return
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(100)          // Set the max number of open connections
	sqlDB.SetMaxIdleConns(10)           // Set the max number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set the maximum connection lifetime

	// Open CSV file
	filePath := "tes150000data.csv"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read CSV file
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	// Initialize a slice to hold the data
	var data []DataRow

	// Loop through the rows and append the data into the slice
	for i, row := range rows {
		if i == 0 {
			// Skip header row
			continue
		}

		if i == 145000 {
			break
		}

		year, err := strconv.Atoi(row[2])
		if err != nil {
			fmt.Println("err: ", err.Error())
			return
		}

		if row[11] == "" {
			continue
		}

		value := strings.TrimSpace(row[11])
		if value == "" {
			continue
		}

		if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
			// Remove the parentheses and add a minus sign
			value = "-" + value[1:len(value)-1]
		}

		value = strings.ReplaceAll(value, ",", "")

		if value == "-" {
			continue
		}

		dv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Println("err: ", err.Error())
			return
		}

		// Append data to the slice
		data = append(data, DataRow{
			Zone:              row[0],
			AreaName:          row[1],
			TimePeriod:        row[2],
			Source:            row[3],
			Nb:                row[4],
			Sector:            cleanData(row[5]),
			Subsector:         cleanData(row[6]),
			IndicatorMetadata: row[7],
			Unit:              row[8],
			StartYear:         year - 1,
			EndYear:           year,
			DataValue:         dv,
		})
	}

	fmt.Println("Length: ", len(data))

	// Initialize time tracking
	timeStart := time.Now()
	fmt.Println("time start: ", timeStart)

	// Semaphore channel to limit the number of concurrent goroutines
	const maxConcurrentGoroutines = 10000
	sem := make(chan struct{}, maxConcurrentGoroutines) // Semaphore channel with a buffer size of maxConcurrentGoroutines

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish

	// Iterate through the data and launch goroutines
	for _, row := range data {
		wg.Add(1)

		// Launch a goroutine for each row, with a limit on concurrency
		go func(row DataRow) {
			defer wg.Done()

			// Acquire a token from the semaphore (this will block if maxConcurrency is reached)
			sem <- struct{}{}

			// Variables to hold IDs
			var regionID, sectorID, subsectorID, indicatorID int

			// Raw SQL Query 1: Get Region ID
			rqOne := `SELECT id FROM "Region" WHERE name = ?`
			if err := db.Raw(rqOne, row.AreaName).Scan(&regionID).Error; err != nil {
				// If error occurs, release the semaphore and return without doing anything
				<-sem
				return
			}
			fmt.Printf("Region ID: %d\n", regionID)

			// Raw SQL Query 2: Get Sector ID
			rqTwo := `SELECT id FROM "Sector" WHERE name = ?`
			if err := db.Raw(rqTwo, row.Sector).Scan(&sectorID).Error; err != nil {
				// If error occurs, release the semaphore and return without doing anything
				<-sem
				return
			}
			fmt.Printf("Sector ID: %d\n", sectorID)

			// Raw SQL Query 3: Get Subsector ID
			rqThree := `SELECT id FROM "Subsector" WHERE name = ? AND "sectorId" = ?`
			if err := db.Raw(rqThree, row.Subsector, sectorID).Scan(&subsectorID).Error; err != nil {
				// If error occurs, release the semaphore and return without doing anything
				<-sem
				return
			}
			fmt.Printf("Subsector ID: %d\n", subsectorID)

			// Raw SQL Query 4: Get Indicator ID
			rqFour := `SELECT id FROM "Indicator" WHERE name = ? AND "subsectorId" = ?`
			if err := db.Raw(rqFour, row.IndicatorMetadata, subsectorID).Scan(&indicatorID).Error; err != nil {
				// If error occurs, release the semaphore and return without doing anything
				<-sem
				return
			}
			fmt.Printf("Indicator ID: %d\n", indicatorID)

			// Insert or update the data
			rqFive := `INSERT INTO "Datum" ("value", "startYear", "endYear", "regionId", "indicatorId")
					   VALUES (?, ?, ?, ?, ?)
					   ON CONFLICT ("indicatorId", "regionId", "startYear", "endYear")
					   DO UPDATE SET value = EXCLUDED.value`

			// Execute the raw query
			if err := db.Exec(rqFive, row.DataValue, row.StartYear, row.EndYear, regionID, indicatorID).Error; err != nil {
				// If error occurs, release the semaphore and return without doing anything
				<-sem
				return
			}

			// Release the semaphore token after the goroutine finishes
			<-sem
		}(row) // Pass the row to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Track and print the total execution time
	timeEnd := time.Now()
	duration := timeEnd.Sub(timeStart)
	fmt.Println("Duration: ", duration)

	// Additional message after all goroutines have finished
	fmt.Println("hello")
}
