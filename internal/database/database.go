package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &DB{db}
	if err := database.createTables(); err != nil {
		return nil, err
	}

	if err := database.updateTables(); err != nil {
		return nil, err
	}

	return database, nil
}

func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL UNIQUE,
		provider TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		ingredients TEXT NOT NULL,
		number_of_ingredients INTEGER NOT NULL,
		directions TEXT NOT NULL,
		notes TEXT NOT NULL,
		prep_time TEXT NOT NULL,
		cook_time TEXT NOT NULL,
		servings INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	db.MustExec(query)

	return nil
}

func (db *DB) updateTables() error {
	migrations := []string{
		`ALTER TABLE recipes ADD COLUMN notes TEXT DEFAULT '' NOT NULL`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			// SQLite will return an error if column already exists
			// This is expected behavior for migrations
			continue
		}
	}

	return nil
}
