package bunovel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"io/fs"
	"os"
	"testing"
)

// GetTestPostgres returns an instance for usage in a test suite.
func GetTestPostgres(t *testing.T, migrations []fs.FS) (*bun.DB, *sql.DB) {
	db, sqlDB, err := NewClient(context.TODO(), Config{
		Driver:      NewPGDriverWithDSN(os.Getenv("POSTGRES_URL")),
		Migrations:  &MigrateConfig{Files: migrations},
		ResetOnConn: true,
	})
	require.NoError(t, err)

	return db, sqlDB
}

func RunTransactionalTest[Fixtures any](db bun.IDB, fixtures []Fixtures, call func(ctx context.Context, tx bun.Tx)) error {
	tx, err := db.BeginTx(context.TODO(), new(sql.TxOptions))
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, data := range fixtures {
		if _, err := tx.NewInsert().Model(data).Exec(context.TODO()); err != nil {
			return fmt.Errorf("failed to insert data: %w", err)
		}
	}

	call(context.TODO(), tx)
	return nil
}
