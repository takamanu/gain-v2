package database

import (
	"crypto/tls"
	"fmt"
	"gain-v2/configs"
	"net/http"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/labstack/gommon/log"
)

func InitElasticSearch(c configs.ProgrammingConfig) (*es.Client, error) {
	cfg := es.Config{
		Addresses: c.DBElasticAddress,
		Username:  c.DBElasticUsername,
		Password:  c.DBElasticPassword,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	es, err := es.NewClient(cfg)

	if err != nil {
		return nil, err
	}

	res, err := es.Cluster.Health()
	if err != nil {
		log.Error("Terjadi kesalahan saat memeriksa kesehatan cluster Elasticsearch, error:", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Error("Elasticsearch returned an error response:", res.String())
		return nil, fmt.Errorf("elasticsearch returned an error response: %s", res.String())
	}

	fmt.Println("Elasticsearch cluster health:", res.String())

	return es, nil
}
