package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./veterinary.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 2. เรียกใช้ฟังก์ชันสร้างตาราง
	createTables()
	log.Println("Database initialized successfully!")
}

func createTables() {

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		name TEXT NOT NULL,
		role TEXT NOT NULL
	);`

	createSlotsTable := `
	CREATE TABLE IF NOT EXISTS slots (
    	id TEXT PRIMARY KEY,
    	vet_id TEXT NOT NULL,
    	date TEXT NOT NULL,           
    	time_period TEXT NOT NULL,   
    	slot_limit INTEGER DEFAULT 1,
    	FOREIGN KEY(vet_id) REFERENCES users(id)
	);`

	createAppointmentsTable := `
	CREATE TABLE IF NOT EXISTS appointments (
		id TEXT PRIMARY KEY,
		slot_id TEXT NOT NULL,
		pet_name TEXT NOT NULL,
		pet_type TEXT NOT NULL,
		client_name TEXT NOT NULL,
		reason TEXT,
		status TEXT DEFAULT 'pending',
		FOREIGN KEY(slot_id) REFERENCES slots(id)
	);`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	_, err = DB.Exec(createSlotsTable)
	if err != nil {
		log.Fatal("Error creating slots table:", err)
	}

	_, err = DB.Exec(createAppointmentsTable)
	if err != nil {
		log.Fatal("Error creating appointments table:", err)
	}
}
