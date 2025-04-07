package database

import (
	"database/sql"
	"os"
	"sync"

	"github.com/LeonEstrak/retro-drop/backend/constants"
	"github.com/LeonEstrak/retro-drop/backend/internalTypes"
	"github.com/LeonEstrak/retro-drop/backend/internalUtils"
	_ "github.com/mattn/go-sqlite3"
)

var (
	once sync.Once
	db   *Database
)

var logger = internalUtils.GetLogger()

type Database struct {
	sql *sql.DB
}

// GetDB returns a pointer to the SQLite database instance.
//
// The database is opened on the first call to GetDB and the same
// instance is returned on subsequent calls. The database is opened
// at the path specified by the constant DB_PATH.
func GetDB() *Database {
	once.Do(func() {
		sqlObject, err := sql.Open("sqlite3", constants.DB_PATH)
		if err != nil {
			logger.Fatal("Failed to open database: %v", err)
			os.Exit(1)
		}

		// Initialize the DB Object
		db = &Database{sql: sqlObject}
	})

	return db
}

func (db *Database) Close() {
	db.sql.Close()
}

func (db *Database) CreateTables() error {
	_, err := db.sql.Exec(`
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

func (db *Database) DropTables() error {
	_, err := db.sql.Exec(`
		DROP TABLE IF EXISTS games;
	`)
	if err != nil {
		logger.Error("Failed to drop tables: %v", err)
		return err
	}
	return nil
}

func (db *Database) GetGamesFromDB(system string, limit int) ([]internalTypes.Games, error) {
	query := "SELECT id, game_title, system, download_url FROM games"
	queryArgs := []any{}
	if system != "" {
		query += " WHERE system = ?"
		queryArgs = append(queryArgs, system)
	}
	if limit > 0 {
		query += " LIMIT ?"
		queryArgs = append(queryArgs, limit)
	}
	rows, err := db.sql.Query(query, queryArgs...)

	if err != nil {
		logger.Error("Failed to query database: %v", err)
		return nil, err
	}
	defer rows.Close()

	games := []internalTypes.Games{}
	for rows.Next() {
		var id int
		var gameTitle string
		var system string
		var downloadURL string
		if err := rows.Scan(&id, &gameTitle, &system, &downloadURL); err != nil {
			logger.Error("Failed to scan row: %v", err)
			return nil, err
		}
		games = append(games, internalTypes.Games{
			ID:          id,
			GameTitle:   gameTitle,
			System:      system,
			DownloadURL: downloadURL,
		})
	}
	return games, nil
}

func (db *Database) InsertListOfGamesToDB(games []internalTypes.Games) error {
	for _, game := range games {
		_, err := db.sql.Exec("INSERT INTO games (game_title, system, download_url) VALUES (?, ?, ?)", game.GameTitle, game.System, game.DownloadURL)
		if err != nil {
			logger.Error("Failed to insert game into database: %v", err)
			return err
		}
	}
	return nil
}
