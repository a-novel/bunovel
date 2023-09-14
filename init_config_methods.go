package bunovel

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"io/fs"
)

func (config Config) compileOptions() []bun.DBOption {
	options := config.Options

	if config.DiscardUnknownColumns {
		options = append(options, bun.WithDiscardUnknownColumns())
	}

	return options
}

// Convert the configuration migrations to a bun migrate.Migrations object.
func (config *MigrateConfig) compileMigrations() (*migrate.Migrations, error) {
	migrations := migrate.NewMigrations()

	for i, migration := range config.Files {
		if err := migrations.Discover(migration); err != nil {
			return nil, fmt.Errorf("failed to discover filesystem migrations at index %v: %w", i, err)
		}
	}

	for i, migration := range config.Go {
		if err := migrations.Register(migration.Up, migration.Down); err != nil {
			return nil, fmt.Errorf("failed to register go migration at index %v: %w", i, err)
		}
	}

	return migrations, nil
}

// RegisterSQLMigrations adds new sql-based migrations to the current configuration.
// See https://bun.uptrace.dev/guide/migrations.html#sql-based-migrations.
func (config *MigrateConfig) RegisterSQLMigrations(migrations ...fs.FS) {
	config.Files = append(config.Files, migrations...)
}

// RegisterGoMigrations adds new Go-based migrations to the current configuration.
// See https://bun.uptrace.dev/guide/migrations.html#go-based-migrations.
func (config *MigrateConfig) RegisterGoMigrations(migrations ...GoMigration) {
	config.Go = append(config.Go, migrations...)
}

// Execute runs every registered migrations. You can call this method multiple times, as bun knows
// to skip already executed migrations.
func (config *MigrateConfig) Execute(ctx context.Context, db *bun.DB) error {
	migrations, err := config.compileMigrations()
	if err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	groups, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	config.migrations = groups
	return nil
}

// Rollback previously executed migrations. It rollbacks a group at a time.
// https://bun.uptrace.dev/guide/migrations.html#migration-groups-and-rollbacks.
func (config *MigrateConfig) Rollback(ctx context.Context, db *bun.DB, opts ...migrate.MigrationOption) error {
	migrations, err := config.compileMigrations()
	if err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	group, err := migrator.Rollback(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	config.migrations = group
	return nil
}

// RollbackAll rollbacks every registered migration group.
// https://bun.uptrace.dev/guide/migrations.html#migration-groups-and-rollbacks.
func (config *MigrateConfig) RollbackAll(ctx context.Context, db *bun.DB, opts ...migrate.MigrationOption) error {
	migrations, err := config.compileMigrations()
	if err != nil {
		return err
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	group, err := migrator.Rollback(ctx, opts...)
	for len(group.Migrations) > 0 && err == nil {
		group, err = migrator.Rollback(ctx, opts...)
	}
	if err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	config.migrations = group
	return nil
}

// Report returns the status of migrations after an Execute statement. It returns nil if Execute has not been called
// or has failed.
func (config *MigrateConfig) Report() *migrate.MigrationGroup {
	return config.migrations
}
