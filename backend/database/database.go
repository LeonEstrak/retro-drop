package database

import (
	"database/sql"

	"github.com/LeonEstrak/retro-drop/backend/utils"
	_ "github.com/mattn/go-sqlite3"
)

var logger = utils.GetLogger()

// OpenDB opens a SQLite database at the provided dbPath.
// It logs a fatal error and exits if the dbPath is empty
// or if the database cannot be opened. Returns a pointer
// to the sql.DB instance on success.
func OpenDB(dbPath string) *sql.DB {
	if dbPath == "" {
		logger.Fatal("Database path is empty")
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Fatal("Failed to open database: %v", err)
	}
	return db
}

type Games struct {
	ID          int
	GameTitle   string
	System      string
	DownloadURL string
}

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_title TEXT,
			system TEXT,
			download_url TEXT
		);
	`)
	if err != nil {
		logger.Error("Failed to create tables: %v", err)
		return err
	}
	return nil
}

func DropTables(db *sql.DB) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS games;
	`)
	if err != nil {
		logger.Error("Failed to drop tables: %v", err)
		return err
	}
	return nil
}
