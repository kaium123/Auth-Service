package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/spf13/viper"
)

var db *sql.DB

// const (
// 	HOST     = "localhost"
// 	PORT     = 54322
// 	USER     = "sm_user3"
// 	PASSWORD = "12345678"
// 	DBNAME   = "testdb"
// )

func InitDB() *sql.DB {

	HOST := viper.GetString("HOST")
	PORT := viper.GetString("PORT")
	USER := viper.GetString("DB_USER")
	fmt.Println(USER)
	PASSWORD := viper.GetString("PASSWORD")
	DBNAME := viper.GetString("DBNAME")
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME,
	)

	db, err := sql.Open("postgres", connString)
	
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(connString)

	fmt.Println("here")

	// Create the user table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sm_users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE,
			password TEXT,
			name TEXT,
			user_name TEXT,
			phone TEXT,
			website TEXT,
			bio TEXT,
			gender TEXT,
			profile_pic TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS websites (
			id SERIAL PRIMARY KEY,
			url TEXT ,
			user_id INTEGER
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sent_requests (
			id SERIAL PRIMARY KEY,
			"from" INTEGER,
			"to" INTEGER
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS friends (
			id SERIAL PRIMARY KEY,
			"from" INTEGER,
			"to" INTEGER
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("here")

	return db
}
