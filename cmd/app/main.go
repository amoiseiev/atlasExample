package main

import (
	"atlasExample/business/data"
	"atlasExample/business/data/dbschema"
	"atlasExample/foundation/db"
	"fmt"
)

func main() {
	dstDBConfig := db.DBConfig{
		User:       "app",
		Host:       "localhost",
		Name:       "atlas_example",
		DisableTLS: true,
	}

	atlasDevDBConfig := db.DBConfig{
		User:       "app",
		Host:       "localhost",
		Name:       "atlas_example_test",
		DisableTLS: true,
	}

	dstDB, err := db.OpenDB(dstDBConfig)
	if err != nil {
		panic("cannot open the database: " + err.Error())
	}

	atlasDevDB, err := db.OpenDB(atlasDevDBConfig)
	if err != nil {
		panic("cannot open the database: " + err.Error())
	}

	atlas := data.NewAtlas()
	err = atlas.ReconcileWithAtlasSQLSchema(dbschema.SQLFiles, dstDB, atlasDevDB)
	if err != nil {
		panic("database cannot be reconciled with its Atlas Schema: " + err.Error())
	}

	fmt.Println("database has been initialized and reconciled with its schema")
}
