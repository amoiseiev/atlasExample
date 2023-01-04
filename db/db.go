package db

import (
	"context"
	_ "embed"
	"errors"
	"net/url"

	"github.com/jmoiron/sqlx"

	atlasmigrate "ariga.io/atlas/sql/migrate"
	atlaspostgres "ariga.io/atlas/sql/postgres"
	atlasschema "ariga.io/atlas/sql/schema"
)

//go:embed sql/schema.sql
var BuiltInSQLSchema []byte

// DBConfig is the configuration for Postgres connection
type DBConfig struct {
	User       string
	Password   string
	Host       string
	Name       string
	Schema     string
	DisableTLS bool
}

type PostgresDB struct {
	prod       *sqlx.DB
	atlasDevDB *sqlx.DB
}

func (r *PostgresDB) getDBDesiredStateFromAtlas(sqlSchema []byte, devDBAtlasDriver atlasmigrate.Driver) (atlasmigrate.StateReader, error) {
	ctx := context.Background()

	dir := &atlasmigrate.MemDir{}

	if err := dir.WriteFile("schemaAtlas.sql", sqlSchema); err != nil {
		return nil, errors.New("Cannot write into MemDir " + err.Error())
	}

	ex, err := atlasmigrate.NewExecutor(devDBAtlasDriver, dir, atlasmigrate.NopRevisionReadWriter{})
	if err != nil {
		return nil, errors.New("Cannot get new migrate executor " + err.Error())
	}

	sr, err := ex.Replay(ctx, func() atlasmigrate.StateReader {
		return atlasmigrate.RealmConn(devDBAtlasDriver, &atlasschema.InspectRealmOption{})
	}())
	if err != nil {
		return nil, errors.New("Cannot execute replay: " + err.Error())
	}

	return atlasmigrate.Realm(sr), nil
}

func (r *PostgresDB) ReconcileWithAtlasSQLSchema(schemaSQL []byte) error {
	prodDBAtlasDriver, err := atlaspostgres.Open(r.prod.DB)
	if err != nil {
		return errors.New("Error opening source connection driver: " + err.Error())
	}

	devDBAtlasDriver, err := atlaspostgres.Open(r.atlasDevDB.DB)
	if err != nil {
		return errors.New("Error opening dev connection driver: " + err.Error())
	}

	desiredStateReader, err := r.getDBDesiredStateFromAtlas(schemaSQL, devDBAtlasDriver)
	if err != nil {
		return err
	}

	currentStateReader := atlasmigrate.RealmConn(prodDBAtlasDriver, &atlasschema.InspectRealmOption{})

	ctx := context.Background()

	desiredState, err := desiredStateReader.ReadState(ctx)
	if err != nil {
		return errors.New("Cannot read desired state: " + err.Error())
	}

	currentState, err := currentStateReader.ReadState(ctx)
	if err != nil {
		return errors.New("Cannot read current state: " + err.Error())
	}

	changes, err := prodDBAtlasDriver.RealmDiff(currentState, desiredState)
	if err != nil {
		return errors.New("Cannot get state diff: " + err.Error())
	}

	if len(changes) > 0 {
		err = prodDBAtlasDriver.ApplyChanges(ctx, changes)
		if err != nil {
			return errors.New("Failed to apply changes: " + err.Error())
		}
	}

	return nil
}

func New(prodDBConfig, atlasDevDBConfig DBConfig) (*PostgresDB, error) {
	prodDB, err := openDB(prodDBConfig)
	if err != nil {
		return nil, err
	}

	atlasDevDB, err := openDB(atlasDevDBConfig)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{prod: prodDB, atlasDevDB: atlasDevDB}, nil
}

// OpenDB Opens SQLx connection to Postgres
func openDB(cfg DBConfig) (*sqlx.DB, error) {
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

	db, err := sqlx.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}

	return db, nil
}
