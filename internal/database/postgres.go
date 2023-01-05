package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

type Database struct {
	DB *sql.DB
}

// NewDatabase - returns a pointer to a new database object
func NewDatabase() (*Database, error) {
	log.Info("Setting up database connection")

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_DBNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return &Database{}, fmt.Errorf("new database: %w", err)
	}
	return &Database{
		DB: db,
	}, nil
}

// Ping - checks database connection
func (db *Database) Ping(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}
