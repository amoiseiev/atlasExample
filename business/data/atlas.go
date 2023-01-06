package data

import (
	atlasmigrate "ariga.io/atlas/sql/migrate"
	atlaspostgres "ariga.io/atlas/sql/postgres"
	atlasschema "ariga.io/atlas/sql/schema"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/fs"
)

type Atlas struct{}

func NewAtlas() *Atlas {
	return &Atlas{}
}

// getDBDesiredStateFromAtlasSQLDirectory Returns a StateReader for the desired state based on files in the migration dir
func (r *Atlas) getDBDesiredStateFromAtlasSQLDirectory(atlasMigrationDir fs.FS, devDBAtlasDriver atlasmigrate.Driver) (
	atlasmigrate.StateReader, error) {
	// scratchDir is used as not all directories are local and may not allow "write" operations used for check sum
	// calculation.
	scratchDir := &atlasmigrate.MemDir{}

	// getting all sql files from the supplied FS and aborting when no files are found as a failsafe.
	files, err := fs.Glob(atlasMigrationDir, "*.sql")
	if err != nil {
		return nil, fmt.Errorf("db/atlas: can't get the list of migration files: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("db/atlas: no schema files found, aborting")
	}

	// copying schema files from the supplied FS to the scratch Dir.
	for _, file := range files {
		fileContent, err := fs.ReadFile(atlasMigrationDir, file)
		if err != nil {
			return nil, fmt.Errorf("db/atlas: can't read schema files: %w", err)
		}

		err = scratchDir.WriteFile(file, fileContent)
		if err != nil {
			return nil, fmt.Errorf("db/atlas: can't write schema files into MemDir: %w", err)
		}
	}

	// initializing Atlas executor, files in the scratchDir will also be checked against Atlas Dev Database.
	ex, err := atlasmigrate.NewExecutor(devDBAtlasDriver, scratchDir, atlasmigrate.NopRevisionReadWriter{})
	if err != nil {
		return nil, fmt.Errorf("db/atlas: cannot get new migrate executor: %w", err)
	}

	// using background context here for illustration purposes only.
	ctx := context.Background()

	// getting the inspection results from replaying files in the migration directory.
	sr, err := ex.Replay(ctx, func() atlasmigrate.StateReader {
		return atlasmigrate.RealmConn(devDBAtlasDriver, &atlasschema.InspectRealmOption{})
	}())
	if err != nil {
		return nil, fmt.Errorf("db/atlas: cannot execute replay: %w", err)
	}

	return atlasmigrate.Realm(sr), nil
}

// ReconcileWithAtlasSQLSchema Executes unattended schema reconciliation against destination database.
func (r *Atlas) ReconcileWithAtlasSQLSchema(atlasMigrationDir fs.FS, dstDB *sqlx.DB, atlasDevDB *sqlx.DB) error {
	// initializing Atlas Drivers for the destination and Atlas Dev databases.
	dstDBAtlasDriver, err := atlaspostgres.Open(dstDB)
	if err != nil {
		return fmt.Errorf("db/atlas: error opening source connection driver: %w", err)
	}

	atlasDevDBAtlasDriver, err := atlaspostgres.Open(atlasDevDB)
	if err != nil {
		return fmt.Errorf("db/atlas: error opening dev connection driver: %w", err)
	}

	// using background context here for illustration purposes only.
	ctx := context.Background()

	// getting state for our desired and current configs, desired config is validated against Atlas Dev DB.
	desiredStateReader, err := r.getDBDesiredStateFromAtlasSQLDirectory(atlasMigrationDir, atlasDevDBAtlasDriver)
	if err != nil {
		return fmt.Errorf("db/atlas: cannot read desired state reader: %w", err)
	}

	desiredState, err := desiredStateReader.ReadState(ctx)
	if err != nil {
		return fmt.Errorf("db/atlas: cannot read desired state: %w", err)
	}

	currentStateReader := atlasmigrate.RealmConn(dstDBAtlasDriver, &atlasschema.InspectRealmOption{})
	currentState, err := currentStateReader.ReadState(ctx)
	if err != nil {
		return fmt.Errorf("db/atlas: cannot read current state: %w", err)
	}

	// comparing the states and generating a list of changes.
	changes, err := dstDBAtlasDriver.RealmDiff(currentState, desiredState)
	if err != nil {
		return fmt.Errorf("db/atlas: cannot get state diff: %w", err)
	}

	if len(changes) > 0 {
		err = dstDBAtlasDriver.ApplyChanges(ctx, changes)
		if err != nil {
			return fmt.Errorf("db/atlas: failed to apply changes: %w", err)
		}
	}

	return nil
}
