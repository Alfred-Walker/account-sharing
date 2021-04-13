package constants

// IMPORTANT: DO NOT UPLOAD ID OR PASSWORD TO GITHUB OR OTHER REPOSITORIES.
// Use your environment variables or external file to handle private values.
// e.g. BUCKET := os.Getenv("BUCKET")

// use "YOUR OWN" rds setup constants
const DB_USERNAME = "YOUR_DB_USERNAME"
const DB_PASSWORD = "YOUR_DB_PASSWORD"

// you can use your own rds db name
// (Please be sure to sync the name with AWS MySQL DB name configuration.)
const DB_NAME = "accountSharing"

// use "YOUR OWN" DB Host such as my-practice-db.abcdef.ap-northeast-2.rds.amazonaws.com
// (Please be sure to sync the name with AWS MySQL DB configuration.)
const DB_HOST = "my-practice-db.c7hdfqvga72s.ap-northeast-2.rds.amazonaws.com"
const DB_PORT = "3306"
