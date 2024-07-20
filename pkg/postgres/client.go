package postgres

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	// MaxIdleConns is the maximum number of connections in the idle connection
	// pool.
	// Default: 10 * GOMAXPROCS.
	MaxIdleConns int `yaml:"maxIdleConns" env:"MAX_IDLE_CONNS"`

	// MaxOpenConns is the maximum number of open connections to the database.
	// Default: 10 * GOMAXPROCS.
	MaxOpenConns int `yaml:"maxOpenConns" env:"MAX_OPEN_CONNS"`

	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Default: 1 minute.
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime" env:"CONN_MAX_IDLE_TIME"`

	// ConnMaxLifeTime is the maximum amount of time a connection may be reused.
	// Default: 5 minutes.
	ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime" env:"CONN_MAX_LIFE_TIME"`

	Database    string `yaml:"database" env:"DATABASE"`
	Username    string `yaml:"username" env:"USERNAME"`
	Password    string `yaml:"password" env:"PASSWORD"`
	HostPrimary string `yaml:"hostPrimary" env:"HOST_PRIMARY"`
	Port        string `yaml:"port" env:"PORT"`
	SSLMode     string `yaml:"sslmode" env:"SSLMODE"`
}

func NewPgxPool(config Config) (*sqlx.DB, error) {
	conf, err := pgxpool.ParseConfig(config.createDSN())
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	setPgxPoolOptions(conf, config)

	pool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	return sqlx.NewDb(stdlib.OpenDBFromPool(pool), "pgx").Unsafe(), nil
}

func (c Config) createDSN() string {
	host := c.HostPrimary
	kvs := []string{
		"host=" + host,
		"port=" + c.Port,
		"user=" + c.Username,
		"password=" + c.Password,
		"dbname=" + c.Database,
		"sslmode=" + c.SSLMode,
	}

	return strings.Join(kvs, " ")
}

//nolint:gomnd // default values
func setPgxPoolOptions(poolConf *pgxpool.Config, config Config) {
	gomaxprocs := runtime.GOMAXPROCS(0)

	maxOpenConns := 10 * gomaxprocs
	if config.MaxOpenConns != 0 {
		maxOpenConns = config.MaxOpenConns
	}
	connMaxLifetime := 5 * time.Minute
	if config.ConnMaxLifeTime != 0 {
		connMaxLifetime = config.ConnMaxLifeTime
	}
	connMaxIdleTime := 1 * time.Minute
	if config.ConnMaxIdleTime != 0 {
		connMaxIdleTime = config.ConnMaxIdleTime
	}

	poolConf.MaxConns = int32(maxOpenConns)
	poolConf.MaxConnLifetime = connMaxLifetime
	poolConf.MaxConnIdleTime = connMaxIdleTime
}
