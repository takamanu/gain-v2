package configs

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type ProgrammingConfig struct {
	ServerPort            int
	Environment           string
	DBPOSTGRESPort        uint16
	DBPOSTGRESHost        string
	DBPOSTGRESUser        string
	DBPOSTGRESPass        string
	DBPOSTGRESName        string
	DBPOSTGRESModeSSL     string
	DBRedisPort           string
	DBRedisAddress        string
	DBRedisPassword       string
	DBRedisDatabase       int
	DBElasticAddress      []string
	DBElasticUsername     string
	DBElasticPassword     string
	Secret                string
	RefSecret             string
	BaseURL               string
	BucketAccessKeyID     string
	BucketSecretAccessKey string
	BucketRegion          string
	BucketEndpoint        string
	BucketName            string
	SMTPEmail             string
	SMTPPasword           string
	SMTPHost              string
	SMTPPort              int
}

func InitConfig() *ProgrammingConfig {
	var res = new(ProgrammingConfig)
	err := godotenv.Load(".env")

	if err != nil {
		err := godotenv.Load(".env.staging")

		if err != nil {
			return nil
		}

	}

	res, errorRes := loadConfig()

	logrus.Error(errorRes)
	if res == nil {
		logrus.Error("Config: Cannot start program, failed to load configuration")
		return nil
	}

	return res
}

func ReadData() *ProgrammingConfig {
	var data = new(ProgrammingConfig)
	data, _ = loadConfig()

	if data == nil {
		err := godotenv.Load(".env")
		data, errorData := loadConfig()

		fmt.Println(errorData)

		if err != nil || data == nil {
			return nil
		}
	}
	return data
}

func loadConfig() (*ProgrammingConfig, error) {
	var error error
	var res = new(ProgrammingConfig)
	var permit = true

	if val, found := os.LookupEnv("SERVER"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config: Invalid port value,", err.Error())
			permit = false
		}
		res.ServerPort = port
	} else {
		permit = false
		error = errors.New("Port undefined")
	}

	if val, found := os.LookupEnv("DB_POSTGRES_PORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : Invalid port value,", err.Error())
			permit = false
		}
		res.DBPOSTGRESPort = uint16(port)
	} else {
		permit = false
		error = errors.New("DBPOSTGRES Port undefined")
	}

	if val, found := os.LookupEnv("SMTP_HOST"); found {
		res.SMTPHost = val
	} else {
		permit = false
		error = errors.New("SMTP_HOST undefined")
	}

	if val, found := os.LookupEnv("SMTP_PORT"); found {
		port, err := strconv.Atoi(val)
		if err != nil {
			logrus.Error("Config : Invalid smtp port value,", err.Error())
			permit = false
		}
		res.SMTPPort = port
	} else {
		permit = false
		error = errors.New("SMTP_PORT undefined")
	}

	if val, found := os.LookupEnv("SMTP_EMAIL"); found {
		res.SMTPEmail = val
	} else {
		permit = false
		error = errors.New("SMTP_EMAIL undefined")
	}

	if val, found := os.LookupEnv("SMTP_PASSWORD"); found {
		res.SMTPPasword = val
	} else {
		permit = false
		error = errors.New("SMTP_PASSWORD undefined")
	}

	if val, found := os.LookupEnv("DB_POSTGRES_HOST"); found {
		res.DBPOSTGRESHost = val
	} else {
		permit = false
		error = errors.New("DBPOSTGRES Host undefined")
	}

	if val, found := os.LookupEnv("DB_POSTGRES_USER"); found {
		res.DBPOSTGRESUser = val
	} else {
		permit = false
		error = errors.New("DBPOSTGRES User undefined")
	}

	if val, found := os.LookupEnv("DB_POSTGRES_PASS"); found {
		res.DBPOSTGRESPass = val
	} else {
		// permit = false
		// error = errors.New("DBPOSTGRES Pass undefined")
		res.DBPOSTGRESPass = ""
	}

	if val, found := os.LookupEnv("DB_POSTGRES_NAME"); found {
		res.DBPOSTGRESName = val
	} else {
		permit = false
		error = errors.New("DBPOSTGRES Name undefined")
	}

	if val, found := os.LookupEnv("DB_POSTGRES_MODE_SSL"); found {
		res.DBPOSTGRESModeSSL = val
	} else {
		permit = false
		error = errors.New("DBPOSTGRES MODE SSL undefined")
	}

	if val, found := os.LookupEnv("DB_REDIS_PORT"); found {
		res.DBRedisPort = val
	} else {
		permit = false
		error = errors.New("DB_REDIS_PORT undefined")
	}

	if val, found := os.LookupEnv("DB_REDIS_ADDRESS"); found {
		res.DBRedisAddress = val
	} else {
		permit = false
		error = errors.New("DB_REDIS_ADDRESS undefined")
	}

	if val, found := os.LookupEnv("DB_REDIS_PASSWORD"); found {
		res.DBRedisPassword = val
	} else {
		permit = false
		error = errors.New("DB_REDIS_PASSWORD undefined")
	}

	if val, found := os.LookupEnv("DB_REDIS_DATABASE"); found {
		redisDB, err := strconv.Atoi(val)

		if err != nil {
			logrus.Error("Config : Invalid redis db value,", err.Error())
			permit = false
		}
		res.DBRedisDatabase = redisDB
	} else {
		permit = false
		error = errors.New("DB_REDIS_DATABASE undefined")
	}

	if val, found := os.LookupEnv("DB_ELASTIC_USER"); found {
		res.DBElasticUsername = val
	} else {
		permit = false
		error = errors.New("DB_ELASTIC_USER undefined")
	}

	if val, found := os.LookupEnv("DB_ELASTIC_PASS"); found {
		res.DBElasticPassword = val
	} else {
		permit = false
		error = errors.New("DB_ELASTIC_PASS undefined")
	}

	if val, found := os.LookupEnv("DB_ELASTIC_ADDRESS"); found {
		res.DBElasticAddress = append(res.DBElasticAddress, val)
	} else {
		permit = false
		error = errors.New("DB_ELASTIC_ADDRESS undefined")
	}

	if val, found := os.LookupEnv("BASE_URL"); found {
		res.BaseURL = val
	} else {
		res.BaseURL = ""
		// permit = false
		// error = errors.New("BASE_URL undefined")
	}
	if val, found := os.LookupEnv("BUCKET_ACCESS_KEY_ID"); found {
		res.BucketAccessKeyID = val
	} else {
		permit = false
		error = errors.New("Config : Invalid BUCKET ACCESS KEY ID undefined")
	}

	if val, found := os.LookupEnv("BUCKET_SECRET_ACCESS_KEY"); found {
		res.BucketSecretAccessKey = val
	} else {
		permit = false
		error = errors.New("Config : Invalid BUCKET SECRET ACCESS KEY undefined")
	}

	if val, found := os.LookupEnv("BUCKET_REGION"); found {
		res.BucketRegion = val
	} else {
		permit = false
		error = errors.New("Config : Invalid BUCKET REGION undefined")
	}

	if val, found := os.LookupEnv("BUCKET_ENDPOINT"); found {
		res.BucketEndpoint = val
	} else {
		permit = false
		error = errors.New("Config : Invalid BUCKET ENDPOINT undefined")
	}

	if val, found := os.LookupEnv("BUCKET_NAME"); found {
		res.BucketName = val
	} else {
		permit = false
		error = errors.New("Config : Invalid BUCKET NAME undefined")
	}

	if !permit {
		return nil, error
	}

	return res, nil
}
