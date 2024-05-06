package models

import (
	"log"
	"os"
	//"path"

	"github.com/dchote/survey-engine/config"

	"github.com/jinzhu/gorm"

	// the blank import is the way GORM wants you to do it
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//import _ "github.com/jinzhu/gorm/dialects/mysql"
	// import _ "github.com/jinzhu/gorm/dialects/postgres"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"
)

var (
	database *gorm.DB
)

// OpenDatabase opens the SQLite db
func OpenDatabase(cfg config.SurveyConfig) *gorm.DB {
	var err error

	log.Println("opening database: " + cfg.Database.Path)

	database, err = gorm.Open(cfg.Database.Type, cfg.Database.Path)
	//defer database.Close()

	if err != nil {
		log.Fatalf("ERROR openening database: %v\n\n", err.Error())
		os.Exit(-1)
	}

	return database
}

// CloseDatabase closes the DB connection
func CloseDatabase() {
	database.Close()
}
