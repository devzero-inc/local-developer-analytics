package database

import "lda/logging"

func RunMigrations() {
	ensureMigrationTableExists()
	createProcessesTable()
	createCommandsTable()
}

func ensureMigrationTableExists() {
	createMigrationTableSQL := `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        migration_name TEXT NOT NULL UNIQUE
    );`

	_, err := DB.Exec(createMigrationTableSQL)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to create schema_migrations table")
	}
}

func createProcessesTable() {
	migrationName := "create_processes_table"
	if !migrationApplied(migrationName) {
		createProcessesTableSQL := `
		CREATE TABLE IF NOT EXISTS processes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pid INTEGER NOT NULL,
			name TEXT NOT NULL,
			status TEXT,
			start_time INTEGER,
			end_time INTEGER,
			execution_time INTEGER,
			os TEXT,
			platform TEXT,
			platform_family TEXT,
			cpu_usage REAL,
			used_memory REAL
		);`

		_, err := DB.Exec(createProcessesTableSQL)
		if err != nil {
			logging.Log.Fatal().Err(err).Msg("Failed to create processes table")
		}
		recordMigration(migrationName)
	}
}

func createCommandsTable() {
	migrationName := "create_commands_table"
	if !migrationApplied(migrationName) {
		createCommandsTableSQL := `
		CREATE TABLE IF NOT EXISTS commands (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT NOT NULL,
			command TEXT NOT NULL,
			user TEXT,
			directory TEXT,
			execution_time INTEGER,
			start_time INTEGER,
			end_time INTEGER 
		);`

		_, err := DB.Exec(createCommandsTableSQL)
		if err != nil {
			logging.Log.Fatal().Err(err).Msg("Failed to create commands table")
		}
		recordMigration(migrationName)
	}
}

func migrationApplied(migrationName string) bool {
	var count int
	err := DB.Get(&count, "SELECT COUNT(*) FROM schema_migrations WHERE migration_name = ?", migrationName)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to query schema_migrations table")
	}
	return count > 0
}

func recordMigration(migrationName string) {
	_, err := DB.Exec("INSERT INTO schema_migrations (migration_name) VALUES (?)", migrationName)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Failed to record migration")
	}
}
