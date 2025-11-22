package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) getMigrate() (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(m.db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create mysql driver: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", m.migrationsDir),
		"mysql",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create migration instance: %w", err)
	}

	return migration, nil
}

func (m *Migrator) Up() error {
	migration, err := m.getMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

func (m *Migrator) Down() error {
	migration, err := m.getMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration down failed: %w", err)
	}

	fmt.Println("Migrations rolled back successfully")
	return nil
}

func (m *Migrator) Steps(n int) error {
	migration, err := m.getMigrate()
	if err != nil {
		return err
	}
	defer migration.Close()

	if err := migration.Steps(n); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration steps failed: %w", err)
	}

	fmt.Printf("Migrated %d steps successfully\n", n)
	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	migration, err := m.getMigrate()
	if err != nil {
		return 0, false, err
	}
	defer migration.Close()

	version, dirty, err := migration.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("could not get version: %w", err)
	}

	return version, dirty, nil
}