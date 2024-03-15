package database

import (
	"lda/logging"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const DBPath = "/tmp/lda.db"

var DB *sqlx.DB

func Setup() {

	db, err := sqlx.Connect("sqlite3", DBPath)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to setup database")
	}

	DB = db
}
