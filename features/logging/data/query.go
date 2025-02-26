package data

import (
	"encoding/json"
	"fmt"
	"gain-v2/features/logging"
	"log"
	"strings"
	"time"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type LoggingData struct {
	rds *redis.Client
	es  *es.Client
}

func NewData(redis *redis.Client, es *es.Client) logging.LoggingDataInterface {
	return &LoggingData{
		rds: redis,
		es:  es,
	}
}

func (ld *LoggingData) AddLog(newData logging.LogData) (*logging.LogData, error) {
	uniqueID := uuid.Must(uuid.NewRandom()).String()

	newData.CreatedAt = time.Now()
	newData.LogID = uniqueID

	reqData, err := json.Marshal(newData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling document: %w", err)
	}

	res, err := ld.es.Index(
		LogData{}.IndexName(),
		strings.NewReader(string(reqData)),
		ld.es.Index.WithDocumentID(uniqueID),
		ld.es.Index.WithRefresh("true"),
	)

	if err != nil {
		return nil, fmt.Errorf("error indexing document: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error indexing document ID=%s: %s", uniqueID, res.String())
	}

	log.Printf("[%s] Document ID=%s indexed successfully\n", res.Status(), uniqueID)
	return &newData, nil
}

func (ld *LoggingData) ViewOneLog(logID string) (*logging.LogData, error) {
	var logData logging.LogData
	var resBody map[string]interface{}

	res, err := ld.es.Get(
		LogData{}.IndexName(),
		logID,
		ld.es.Get.WithPretty(),
	)

	if err != nil {
		return nil, fmt.Errorf("error getting document: %w", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting document ID=%s: %s", logID, res.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	log.Printf("Raw response: %+v", resBody)

	source, found := resBody["_source"].(map[string]interface{})
	if !found {
		return nil, fmt.Errorf("document not found for ID=%s", logID)
	}

	sourceBytes, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("error marshalling source: %w", err)
	}

	if err := json.Unmarshal(sourceBytes, &logData); err != nil {
		return nil, fmt.Errorf("error unmarshalling source: %w", err)
	}

	return &logData, nil
}

func (ld *LoggingData) ViewLog() ([]logging.LogData, error) {
	var logData logging.LogData

	query := `{
        "query": {
            "match_all": {}
        }
    }`

	res, err := ld.es.Search(
		ld.es.Search.WithIndex(LogData{}.IndexName()),
		// ld.es.Search.WithIndex("tes"),
		ld.es.Search.WithBody(strings.NewReader(query)),
		ld.es.Search.WithSize(1000),
		ld.es.Search.WithPretty(),
	)

	if err != nil {
		return nil, fmt.Errorf("error searching documents: %w", err)
	}

	defer res.Body.Close()

	var result map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	if errorInfo, ok := result["error"].(map[string]interface{}); ok {
		if status, ok := result["status"].(float64); ok {

			return nil, fmt.Errorf("index error %v: %v", status, LogData{}.IndexName())
		}
		return nil, fmt.Errorf("error in response: %v", errorInfo)
	}

	fmt.Println("Res: ", result)

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	logs := make([]logging.LogData, 0, len(hits))

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})

		sourceBytes, err := json.Marshal(source)

		if err != nil {
			return nil, fmt.Errorf("error marshalling source: %w", err)
		}

		if err := json.Unmarshal(sourceBytes, &logData); err != nil {
			return nil, fmt.Errorf("error unmarshalling source: %w", err)
		}

		logs = append(logs, logData)
	}

	return logs, nil
}

func (ld *LoggingData) DeleteLog(logID string) (bool, error) {
	res, err := ld.es.Delete(
		LogData{}.IndexName(),
		logID,
		ld.es.Delete.WithRefresh("true"),
	)
	if err != nil {
		return false, fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, fmt.Errorf("error deleting document ID=%s: %s", logID, res.String())
	}

	log.Printf("[%s] Document ID=%s deleted successfully\n", res.Status(), logID)

	return true, nil
}
