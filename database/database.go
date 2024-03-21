package database

import (
	"lda/config"
	"lda/logging"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection.
var DB *sqlx.DB

// Setup initializes the database connection.
func Setup() {

	dbPath := filepath.Join(config.LdaDir, "lda.db")

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to setup database")
	}
	DB = db
}
