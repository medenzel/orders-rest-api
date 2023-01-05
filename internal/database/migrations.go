package database

import (
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// Migrate - runs all migrations in the migrations folder
func (db *Database) Migrate(migrationsFS fs.FS, dir string) error {
	if dir == "" {
		dir = "."
	}

	goose.SetBaseFS(migrationsFS)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate setting dialect: %w", err)
	}

	err = goose.Up(db.DB, dir)
	if err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}
