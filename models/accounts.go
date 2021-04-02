package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)


// * Create rds table before starting lambda handler at main.go
// below one is sample table and data creation queries for this project
//
// sample query for MYSQL
// CREATE table `accounts` (
// 	id int primary key auto_increment,
//     account_name varchar(255),
//     end_time varchar(255),
//     occupier varchar(255)
// );
//
// insert into accounts(account_name, end_time, occupier)
// values ('DEFAULT_SHARING_ACCOUNT_NAME','','');

// (This tutorial does not assume multiple account sharing yet.)

// Account model for RDS (MySQL)
// (Struct fields must be public in order for MarshalMap to do its magic)
// https://stackoverflow.com/questions/56827932/one-or-more-parameter-values-were-invalid-missing-the-key-id-in-the-item-status
type Account struct {
	Id          int    `json:"id" gorm:"primaryKey"`
	AccountName string `json:"account" gorm:"column:account_name"`
	EndTime     string `json:"endtime" gorm:"column:end_time"`
	Occuppier   string `json:"occupier" gorm:"column:occupier"`
}

// Custom error struct for account request format errors.
type AccountRequestFormatError struct {
	err string
}

// Custom Error interface implementation for AccountRequestFormatError
func (e *AccountRequestFormatError) Error() string {
	return fmt.Sprintf("syntax error - %s", e.err)
}

// GetFirstAccount retrieves first record from the GORM db.
// Receives GORM db instance.
// It returns first account record.
func GetFirstAccount(db *gorm.DB) *Account {
	var account Account
	db.Model(&Account{}).First(&account)

	return &account
}

// OccupyAccount registers new occupier info to the rds.
// It receives GORM db instance, occupier's name, and endtime string.
// It returns new occupier account record or returns parse error.
func OccupyAccount(db *gorm.DB, occupier string, endtime string) (*Account, error) {
	var account Account

	if occupier == "" {
		return nil, &AccountRequestFormatError{err: "accounts: Empty string name is not allowed"}
	}

	db.Model(&Account{}).First(&account)

	// need to change your own location if necessary
	loc, err := time.LoadLocation("Asia/Seoul")
	localTime := time.Now().In(loc)

	if err != nil {
		log.Printf("accounts: Failed to load location")
		return nil, err
	}

	// validate and apply new endtime (manual)
	// assume that only hour:minute info is transferred from clients
	// 15:04 is custom format for time
	// https://gobyexample.com/time-formatting-parsing
	parsed, err := time.Parse("15:04", endtime)

	if err != nil {
		log.Printf("accounts: Failed to parse request endtime")
		return nil, err
	}

	// set year-month-date info and location
	// (clients may send hour-minuts info only)
	parsed = time.Date(
		localTime.Year(), 
		localTime.Month(), 
		localTime.Day(), 
		parsed.Hour(), 
		parsed.Minute(), 
		parsed.Second(), 
		0, 
		loc,
	)

	// set EndTime
	// assumes that parsed data is local time (Asia/Seoul) already. === no need to call In(loc) for parsed
	if time.Now().In(loc).After(parsed) {
		account.EndTime = parsed.AddDate(0, 0, 1).Format("2006-01-02-15:04")
	} else {
		account.EndTime = parsed.Format("2006-01-02-15:04")
	}

	// change occupier
	account.Occuppier = occupier

	// save changes
	db.Save(&account)

	return &account, nil
}

// ReleaseAccount is called when sharing account is freed from an occupier.
// It receives GORM db instance and returns initialized account record.
func ReleaseAccount(db *gorm.DB) *Account {
	account := GetFirstAccount(db)
	account.EndTime = ""
	account.Occuppier = ""

	// save changes
	db.Save(&account)

	return account
}

// CheckEndTime checks current occupier's endtime is passed or not,
// and if current occupier's endtime is passed, it initializes endtime and occupier name.
// It receives GORM db instance and returns current occupier account record or returns parse error.
func (account Account) CheckEndTime(db *gorm.DB) *Account {
	if account.EndTime == "" {
		// no occupier, so do nothing
		return &account
	}

	parsedTime, err := time.Parse("2006-01-02-15:04", account.EndTime)

	if err != nil {
		log.Printf("accounts: Failed to parse rds endtime record")
		return &account
	}

	// timezone for korean
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now().Local().In(loc)

	// set year-month-date info and location
	parsedTime = time.Date(
		parsedTime.Year(), 
		parsedTime.Month(), 
		parsedTime.Day(), 
		parsedTime.Hour(), 
		parsedTime.Minute(), 
		parsedTime.Second(), 
		0, 
		loc,
	)

	// initialize occupation at endtime
	if now.After(parsedTime) {
		account.EndTime = ""
		account.Occuppier = ""

		db.Save(account)
	}
	return &account
}
