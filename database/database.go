package database

import (
	"lda/config"
	"lda/logging"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func Setup() {

	dbPath := filepath.Join(config.LdaDir, "lda.db")

	logging.Log.Debug().Msgf("Setting up database at %s", dbPath)

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to setup database")
	}
	DB = db
}
