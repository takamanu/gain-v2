package data

import (
	"context"
	"fmt"
	"gain-v2/features/upload_csv"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UploadCSVData struct {
	db  *gorm.DB
	rds *redis.Client
	ctx context.Context
}

func NewData(redis *redis.Client, db *gorm.DB, ctx context.Context) upload_csv.UploadCSVDataInterface {
	return &UploadCSVData{
		db:  db,
		rds: redis,
		ctx: ctx,
	}
}

func (ud *UploadCSVData) GetRegionID(areaName string) (int, error) {
	var regionID int
	query := `SELECT id FROM "Region" WHERE name = ?`
	err := ud.db.Raw(query, areaName).Scan(&regionID).Error
	if err != nil {
		return 0, err
	}

	return regionID, nil
}

// GetSectorID retrieves the ID from the "Sector" table
func (ud *UploadCSVData) GetSectorID(sectorName string) (int, error) {
	var sectorID int
	query := `SELECT id FROM "Sector" WHERE name = ?`
	err := ud.db.Raw(query, sectorName).Scan(&sectorID).Error
	if err != nil {
		return 0, err
	}
	fmt.Printf("Sector ID: %d\n", sectorID)

	return sectorID, nil
}

// GetSubsectorID retrieves the ID from the "Subsector" table
func (ud *UploadCSVData) GetSubsectorID(subsectorName string, sectorID int) (int, error) {
	var subsectorID int
	query := `SELECT id FROM "Subsector" WHERE name = ? AND "sectorId" = ?`
	err := ud.db.Raw(query, subsectorName, sectorID).Scan(&subsectorID).Error
	if err != nil {
		return 0, err
	}
	fmt.Printf("Subsector ID: %d\n", subsectorID)

	return subsectorID, nil
}

// GetIndicatorID retrieves the ID from the "Indicator" table
func (ud *UploadCSVData) GetIndicatorID(indicatorName string, subsectorID int) (int, error) {
	var indicatorID int
	query := `SELECT id FROM "Indicator" WHERE name = ? AND "subsectorId" = ?`
	err := ud.db.Raw(query, indicatorName, subsectorID).Scan(&indicatorID).Error
	if err != nil {
		return 0, err
	}

	fmt.Printf("Indicator ID: %d\n", indicatorID)
	return indicatorID, nil
}

// InsertOrUpdateDatum inserts or updates the "Datum" table
func (ud *UploadCSVData) InsertOrUpdateDatum(dataValue float64, startYear, endYear, regionID, indicatorID int) error {
	query := `INSERT INTO "Datum" ("value", "startYear", "endYear", "regionId", "indicatorId")
	          VALUES (?, ?, ?, ?, ?)
	          ON CONFLICT ("indicatorId", "regionId", "startYear", "endYear")
	          DO UPDATE SET value = EXCLUDED.value`
	err := ud.db.Exec(query, dataValue, startYear, endYear, regionID, indicatorID).Error
	if err != nil {
		return err
	}

	fmt.Println("Data inserted/updated successfully!")
	return nil
}

func (us *UploadCSVData) PreloadIDs() (map[string]int, map[string]int, map[string]map[string]int, map[string]map[string]int, error) {
	regionMap := make(map[string]int)
	sectorMap := make(map[string]int)
	subsectorMap := make(map[string]map[string]int) // sector -> subsector
	indicatorMap := make(map[string]map[string]int) // subsector -> indicator

	// Load Region IDs
	var regions []struct {
		ID   int
		Name string
	}
	err := us.db.Raw(`SELECT id, name FROM "Region"`).Scan(&regions).Error
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, r := range regions {
		regionMap[r.Name] = r.ID
	}

	// Load Sector IDs
	var sectors []struct {
		ID   int
		Name string
	}
	err = us.db.Raw(`SELECT id, name FROM "Sector"`).Scan(&sectors).Error
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, s := range sectors {
		sectorMap[s.Name] = s.ID
	}

	// Load Subsector IDs
	var subsectors []struct {
		ID       int
		Name     string
		SectorID int
	}
	err = us.db.Raw(`SELECT id, name, "sectorId" FROM "Subsector"`).Scan(&subsectors).Error
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, ss := range subsectors {
		if _, exists := subsectorMap[ss.Name]; !exists {
			subsectorMap[ss.Name] = make(map[string]int)
		}
		subsectorMap[ss.Name][fmt.Sprintf("%d", ss.SectorID)] = ss.ID
	}

	// Load Indicator IDs
	var indicators []struct {
		ID          int
		Name        string
		SubsectorID int
	}
	err = us.db.Raw(`SELECT id, name, "subsectorId" FROM "Indicator"`).Scan(&indicators).Error
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, ind := range indicators {
		if _, exists := indicatorMap[ind.Name]; !exists {
			indicatorMap[ind.Name] = make(map[string]int)
		}
		indicatorMap[ind.Name][fmt.Sprintf("%d", ind.SubsectorID)] = ind.ID
	}

	return regionMap, sectorMap, subsectorMap, indicatorMap, nil
}

func (ud *UploadCSVData) BulkInsertGainData(data []upload_csv.GainCSVData) error {
	if len(data) == 0 {
		return nil
	}

	// Create a slice of values for bulk insert
	values := []interface{}{}
	for _, d := range data {
		values = append(values, d.DataValue, d.StartYear, d.EndYear, d.RegionID, d.IndicatorID)
	}

	fmt.Println(data)

	// Prepare SQL query for bulk insert
	query := `INSERT INTO "Datum" ("value", "startYear", "endYear", "regionId", "indicatorId")
	          VALUES ` + strings.Repeat("(?, ?, ?, ?, ?),", len(data)-1) + "(?, ?, ?, ?, ?) " +
		`ON CONFLICT ("indicatorId", "regionId", "startYear", "endYear") 
	          DO UPDATE SET value = EXCLUDED.value`

	// Execute bulk insert
	err := ud.db.Exec(query, values...).Error
	if err != nil {
		return err
	}

	fmt.Printf("✅ Bulk Inserted %d rows successfully!\n", len(data))
	return nil
}

func (ud *UploadCSVData) InsertGainData(data upload_csv.GainCSVData) error {
	query := `INSERT INTO "Datum" ("value", "startYear", "endYear", "regionId", "indicatorId")
	          VALUES (?, ?, ?, ?, ?)
	          ON CONFLICT ("indicatorId", "regionId", "startYear", "endYear") 
	          DO UPDATE SET value = EXCLUDED.value`

	err := ud.db.Exec(query, data.DataValue, data.StartYear, data.EndYear, data.RegionID, data.IndicatorID).Error
	if err != nil {
		return err
	}

	return nil
}

func (us *UploadCSVData) GetRegions() ([]upload_csv.Region, error) {
	var regions []upload_csv.Region
	err := us.db.Table("Region").Find(&regions).Error
	if err != nil {
		return nil, err
	}
	return regions, nil
}

func (us *UploadCSVData) GetSectors() ([]upload_csv.Sector, error) {
	var sectors []upload_csv.Sector
	err := us.db.Table("Sector").Preload("Subsectors.Indicators").Find(&sectors).Error
	if err != nil {
		return nil, err
	}
	return sectors, nil
}
