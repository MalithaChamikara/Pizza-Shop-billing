package database

import (
	"database/sql"
	"log"
	"os"

	 //loads env variables from .env file
	"github.com/joho/godotenv" 
	//allows us to use mysql driver to connect to the database
	_ "github.com/go-sql-driver/mysql"  

)

//variable to hold the database connection pool
var DB *sql.DB 

//function to connect to the database
func Connect()  {
	//Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	//retrieve the environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	//build the Data Souce Name (DSN)
	DSN := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

	//connect to the database
	DB, err = sql.Open("mysql", DSN)
	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)
	}
	//verify the connection to the database
	if err = DB.Ping(); err != nil {
		log.Fatalf("Error verifying database connection %v", err)
	}

	//log a message to  the console to indicate that the connection was successful
	log.Println("Connected to the database")

}