package bunovel

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

// Connect opens a new connection with pgdriver connector.
func (config *PGDriver) Connect(options ...bun.DBOption) (*bun.DB, *sql.DB, error) {
	pgconn := pgdriver.NewConnector(config.parseOptions()...)
	sqlDB := sql.OpenDB(pgconn)
	database := bun.NewDB(sqlDB, pgdialect.New(), options...)
	return database, sqlDB, nil
}

func parseStringOption(option string, callback func(string) pgdriver.Option, options *[]pgdriver.Option) {
	if option == "" {
		return
	}

	*options = append(*options, callback(option))
}

func parseDurationOption(option time.Duration, callback func(time.Duration) pgdriver.Option, options *[]pgdriver.Option) {
	if option == 0 {
		return
	}

	*options = append(*options, callback(option))
}

func parseMapOption[K comparable, V any](option map[K]V, callback func(map[K]V) pgdriver.Option, options *[]pgdriver.Option) {
	if option == nil {
		return
	}

	*options = append(*options, callback(option))
}

func parsePTROption[O any](option *O, callback func(O) pgdriver.Option, options *[]pgdriver.Option) {
	if option == nil {
		return
	}

	*options = append(*options, callback(*option))
}

func (config *PGDriver) parseOptions() []pgdriver.Option {
	options := config.Options

	parseStringOption(config.DSN, pgdriver.WithDSN, &options)

	parseStringOption(config.Addr, pgdriver.WithAddr, &options)
	parseStringOption(config.AppName, pgdriver.WithApplicationName, &options)
	parseMapOption(config.ConnParams, pgdriver.WithConnParams, &options)
	parseStringOption(config.Database, pgdriver.WithDatabase, &options)
	parseDurationOption(config.DialTimeout, pgdriver.WithDialTimeout, &options)
	parsePTROption(config.Insecure, pgdriver.WithInsecure, &options)
	parseStringOption(config.Network, pgdriver.WithNetwork, &options)
	parseStringOption(config.Password, pgdriver.WithPassword, &options)
	parseDurationOption(config.ReadTimeout, pgdriver.WithReadTimeout, &options)
	parseDurationOption(config.Timeout, pgdriver.WithTimeout, &options)
	parseStringOption(config.User, pgdriver.WithUser, &options)
	parseDurationOption(config.WriteTimeout, pgdriver.WithWriteTimeout, &options)

	if config.TLS != nil {
		options = append(options, pgdriver.WithTLSConfig(config.TLS))
	}

	return options
}
