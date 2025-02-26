package main

import (
	"gain-v2/configs"
	"gain-v2/utils/database"

	"github.com/sirupsen/logrus"
)

func main() {

	var config = configs.InitConfig()

	db, err := database.InitDBPostgres(*config)
	if err != nil {
		logrus.Fatal("Cannot run database: ", err.Error())
	}

	database.MigrateWithDrop(db)
}
