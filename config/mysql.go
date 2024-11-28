package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Importing the MySQL driver package
)

// InitMysql initializes and returns a connection to the MySQL database
func InitMysql() (*sql.DB, error) {
	// Data Source Name (DSN) used to connect to the MySQL database
	dsn := "ralfi:elyanapunya@tcp(127.0.0.1:3306)/movies?charset=utf8mb4&parseTime=True&loc=Local"

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// If there is an error while opening the database connection, log the error and return nil
		log.Println("ERR connect to DB : ", err)
		return nil, err
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		// If pinging the database fails, log the error and return nil
		log.Println("ERR connect to DB : ", err)
		return nil, err
	}

	// Log that the database connection was successful
	log.Println("database connected")

	// Return the established database connection
	return db, nil
}
