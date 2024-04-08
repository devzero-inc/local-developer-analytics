package database

import (
	"fmt"
	"lda/config"
	"lda/util"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB is the database connection.
var DB *sqlx.DB

// Setup initializes the database connection.
func Setup(ldaDir string, user *user.User) {

	dbPath := filepath.Join(ldaDir, "lda.db")

	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Failed to setup database: %s\n", err)
		os.Exit(1)
	}

	if err := util.ChangeFileOwnership(dbPath, user); err != nil {
		fmt.Fprintf(config.SysConfig.ErrOut, "Failed to change ownership of database: %s\n", err)
		os.Exit(1)
	}

	DB = db
}
