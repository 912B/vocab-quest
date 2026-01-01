package database

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	DB = db
	log.Println("Database connected")

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            role TEXT DEFAULT 'user',
            avatar TEXT DEFAULT '👨‍🚀'
        );`,
		`CREATE TABLE IF NOT EXISTS dictionaries (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            description TEXT,
            is_active BOOLEAN DEFAULT 0
        );`,
		`CREATE TABLE IF NOT EXISTS words (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            dictionary_id INTEGER,
            text TEXT NOT NULL,
            definition TEXT,
            difficulty INTEGER DEFAULT 1,
            pronunciation_url TEXT,
            FOREIGN KEY(dictionary_id) REFERENCES dictionaries(id)
        );`,
		`CREATE TABLE IF NOT EXISTS user_progress (
            user_id INTEGER,
            word_id INTEGER,
            attempts INTEGER DEFAULT 0,
            successes INTEGER DEFAULT 0,
            last_played_at DATETIME,
            next_review_at DATETIME,
            interval REAL DEFAULT 1,
            ease_factor REAL DEFAULT 2.5,
            srs_stage INTEGER DEFAULT 0,
            PRIMARY KEY (user_id, word_id),
            FOREIGN KEY(user_id) REFERENCES users(id),
            FOREIGN KEY(word_id) REFERENCES words(id)
        );`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	// Auto-Migration for existing DBs (Idempotent-ish)
	migrations := []string{
		"ALTER TABLE user_progress ADD COLUMN next_review_at DATETIME;",
		"ALTER TABLE user_progress ADD COLUMN interval REAL DEFAULT 1;",
		"ALTER TABLE user_progress ADD COLUMN ease_factor REAL DEFAULT 2.5;",
		"ALTER TABLE user_progress ADD COLUMN srs_stage INTEGER DEFAULT 0;",
	}
	for _, m := range migrations {
		db.Exec(m) // Ignore errors (duplicate column)
	}

	return nil
}
