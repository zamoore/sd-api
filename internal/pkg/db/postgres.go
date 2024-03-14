package database

import (
	"database/sql"
	"fmt"
	"snipdrop-rest-api/internal/pkg/config"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func ConnectDatabase(cfg config.Config) (*sql.DB, error) {
	// Construct the Postgres connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)

	// Attempt to open a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Verify the connection to the database
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Assign the database connection to the global variable
	return db, nil
}
