package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

// InitMysql initializes and returns a connection to the MySQL database
func InitMysql() (*sql.DB, error) {

	// Get credentials from environment variables
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")
	params := os.Getenv("MYSQL_PARAMS")

	// Construct DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, dbName, params)

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("ERR connect to DB:", err)
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Println("ERR ping DB:", err)
		return nil, err
	}

	log.Println("Database connected successfully!")
	return db, nil
}
