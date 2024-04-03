package database

import (
	"fmt"
	"lda/config"
	"os"
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
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to setup database: %s\n", err)
		os.Exit(1)
	}
	DB = db
}
