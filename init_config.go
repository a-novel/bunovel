package bunovel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"io/fs"
	"os"
)

// DriverConfig provides a generic representation of a driver implementation.
type DriverConfig interface {
	// Connect to the database using the driver, and return an active bun.DB object.
	Connect(options ...bun.DBOption) (*bun.DB, *sql.DB, error)
}

// GoMigration represents a set of migration functions for bun.
type GoMigration struct {
	// Up applies a migration.
	Up migrate.MigrationFunc
	// Down rollbacks the changes made by a migration, if applied.
	Down migrate.MigrationFunc
}

// MigrateConfig configures migrations for a given *bun.DB instance.
// See https://bun.uptrace.dev/guide/migrations.html.
type MigrateConfig struct {
	// Files is a list of filesystems to lookup for SQL migration files.
	// See https://bun.uptrace.dev/guide/migrations.html#sql-based-migrations.
	Files []fs.FS
	// Go represents a bunch of GoMigration to execute when loading the driver.
	// See https://bun.uptrace.dev/guide/migrations.html#go-based-migrations.
	Go []GoMigration

	// Cache migration results.
	migrations *migrate.MigrationGroup
}

// Config is the main configuration object for a *bun.DB instance.
type Config struct {
	// Driver used to communicate with the instance.
	Driver DriverConfig
	// Migrations is an optional value to run migrations automatically when loading the driver.
	// See https://bun.uptrace.dev/guide/migrations.html.
	Migrations *MigrateConfig

	// Production optimization for bun. It is recommended to set this to true for production builds.
	// See https://bun.uptrace.dev/guide/running-bun-in-production.html#bun-withdiscardunknowncolumns.
	DiscardUnknownColumns bool
	// ResetOnConn resets the whole database content when opening a new connection. Only use this under test
	// environments.
	ResetOnConn bool

	// Options is a fallback/security, to still allow to pass options in a conventional way. Also, it
	// allows Config to accept new options that have not or cannot (for any reason) be configured with
	// a config object.
	Options []bun.DBOption
}

// NewClient creates a new *bun.DB instance from a Config object.
// The returned bun.DB and sql.DB objects need to be closed manually by the application, in a defer clause.
func NewClient(ctx context.Context, config Config) (*bun.DB, *sql.DB, error) {
	database, sqlDB, err := config.Driver.Connect(config.compileOptions()...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect the database client: %w", err)
	}

	if config.ResetOnConn {
		// Prevents accidental settings.
		if os.Getenv("ENV") != "test" {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, ErrResetOnConnOutsideTests
		}

		// Run all rollbacks.
		if err := config.Migrations.RollbackAll(ctx, database); err != nil {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, fmt.Errorf("failed to reset the database: %w", err)
		}
	}

	// Apply migrations.
	// We do the guard check thing, because bun throws an error if we try to apply migrations when there are none.
	if config.Migrations != nil && (len(config.Migrations.Files) > 0 || len(config.Migrations.Go) > 0) {
		if err := config.Migrations.Execute(ctx, database); err != nil {
			_ = database.Close()
			_ = sqlDB.Close()
			return nil, nil, err
		}
	}

	return database, sqlDB, nil
}

// NewClientWithDriver is an abstraction of NewClient, that allows to completely skip the config object, and
// only pass the driver, without extra options.
func NewClientWithDriver(driver DriverConfig) (*bun.DB, *sql.DB, error) {
	return driver.Connect()
}
