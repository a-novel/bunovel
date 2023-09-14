package bunovel

import (
	"crypto/tls"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

// PGDriver for usage with "github.com/uptrace/bun/driver/pgdriver". You can use this object directly, or call
// NewPGDriver or NewPGDriverWithDSN for shorter declarations.
//
// All options match the ones available under pgdriver:
// https://pkg.go.dev/github.com/uptrace/bun/driver/pgdriver#Option.
type PGDriver struct {
	// DSN is a string used to connect to a psql instance.
	DSN string `json:"dsn" yaml:"dsn"`

	// Network to use, either tcp or unix. Defaults to tcp.
	Network string `json:"network" yaml:"network"`
	// Address of the psql instance, in the format host:port. This parameter is overridden by DSN, if set.
	Addr string `json:"addr" yaml:"addr"`
	// ConnParams are optional parameters for db connection. They are overridden by DSN, if set.
	ConnParams map[string]interface{} `json:"connParams" yaml:"connParams"`
	// Database to connect to. This parameter is overridden by DSN, if set.
	Database string `json:"database" yaml:"database"`
	// User used to connect database. This parameter is overridden by DSN, if set.
	User string `json:"user" yaml:"user"`
	// TLS configuration for connecting to database. This parameter is overridden by DSN, if set wth sslmode or
	// sslrootcert query parameters.
	TLS *tls.Config `json:"-" yaml:"-"`
	// ReadTimeout is the timeout to apply to read requests to the database. This parameter is overridden by DSN,
	// if set with read_timeout query parameter.
	ReadTimeout time.Duration `json:"readTimeout" yaml:"readTimeout"`
	// Timeout is the timeout to apply to all requests to the database. This parameter is overridden by DSN,
	// if set with timeout query parameter.
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// WriteTimeout is the timeout to apply to write requests to the database. This parameter is overridden by DSN,
	// if set with write_timeout query parameter.
	WriteTimeout time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	// DialTimeout is the timeout to apply to dial requests to the database. This parameter is overridden by DSN,
	// if set with dial_timeout or connect_timeout query parameters.
	DialTimeout time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	// AppName sets the application name to be reported in statistics and logs. This parameter is overridden by DSN,
	// if set with application_name query parameter.
	AppName  string `json:"appName" yaml:"appName"`
	Insecure *bool  `json:"insecure" yaml:"insecure"`
	Password string `json:"password" yaml:"password"`

	// Options is a fallback/security, to still allow to pass options in a conventional way. Also it
	// allows Driver to accept new options that have not or cannot (for any reason) be configured within
	// the object.
	Options []pgdriver.Option `json:"-" yaml:"-"`
}

// NewPGDriver generates a new PGDriver object from options declared the legacy way.
// See https://bun.uptrace.dev/postgres/#pgdriver.
//
// For a more granular configuration, use the PGDriver object directly.
func NewPGDriver(options ...pgdriver.Option) *PGDriver {
	return &PGDriver{Options: options}
}

// NewPGDriverWithDSN generates a new PGDriver from a dsn string. You can also pass additional options,
// the legacy way. See https://bun.uptrace.dev/postgres/#pgdriver.
//
// For a more granular configuration, use the Driver object directly.
func NewPGDriverWithDSN(dsn string, options ...pgdriver.Option) *PGDriver {
	return &PGDriver{DSN: dsn, Options: options}
}
