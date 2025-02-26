package database

import (
	"fmt"
	"time"

	"gain-v2/configs"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDBPostgres(c configs.ProgrammingConfig) (*gorm.DB, error) {

	// dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", c.DBPOSTGRESUser, c.DBPOSTGRESPass, c.DBPOSTGRESHost, c.DBPOSTGRESPort, c.DBPOSTGRESName)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Jakarta",
		c.DBPOSTGRESHost,
		c.DBPOSTGRESUser,
		c.DBPOSTGRESPass,
		c.DBPOSTGRESName,
		c.DBPOSTGRESPort,
		c.DBPOSTGRESModeSSL,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NowFunc: func() time.Time {
			utc, _ := time.LoadLocation("Asia/Jakarta")
			return time.Now().In(utc)
		},
	})

	if err != nil {
		log.Error("Terjadi kesalahan pada database, error:", err.Error())
		return nil, err
	}

	return db, nil
}
