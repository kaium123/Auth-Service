package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/spf13/viper"
)

var db *sql.DB

func InitDB() *sql.DB {
	// Open a PostgreSQL database connection (replace with your own connection string)
	dbUrl := viper.GetString("DB_URL")
	var err error
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create the user table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
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
	return db
}
