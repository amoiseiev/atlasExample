package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/url"
)

// DBConfig is the configuration for Postgres connection
type DBConfig struct {
	User       string
	Password   string
	Host       string
	Name       string
	Schema     string
	DisableTLS bool
}

// OpenDB Opens SQLx connection to Postgres
func OpenDB(cfg DBConfig) (*sql.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")
	if cfg.Schema != "" {
		q.Set("search_path", cfg.Schema)
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sql.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}

	return db, nil
}
