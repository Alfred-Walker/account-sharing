package rds

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Alfred-Walker/AccountSharing/constants"
)

var RdsGorm *gorm.DB

// InitRDS initialize GORM instance for RDS.
// It constructs dsn using DB related constants described in constants package.
// It returns error when MYSQL connection is not established.
// (See belows if you need more guide for GORM.)
// https://gorm.io/docs/connecting_to_the_database.html
func InitRDS() error {
	var err error

	dsn := constants.DB_USERNAME +
		":" +
		constants.DB_PASSWORD +
		"@tcp" +
		"(" +
		constants.DB_HOST +
		":" +
		constants.DB_PORT +
		")/" +
		constants.DB_NAME +
		"?" +
		"charset=utf8mb4&parseTime=true&loc=Local"

	RdsGorm, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("RDS connection failed to open!")
	} else {
		log.Printf("RDS setup initialized...")
	}

	return err
}
