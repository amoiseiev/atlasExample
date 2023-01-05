package main

import (
	"atlasExample/db"
	"atlasExample/db/sql"
	_ "embed"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
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
		panic("Cannot open the database: " + err.Error())
	}

	atlasDevDB, err := db.OpenDB(atlasDevDBConfig)
	if err != nil {
		panic("cannot open the database: " + err.Error())
	}

	err = db.ReconcileWithAtlasSQLSchema(sql.SchemaFiles, dstDB, atlasDevDB)
	if err != nil {
		panic("database cannot be reconciled with its Atlas Schema: " + err.Error())
	}

	fmt.Println("database has been initialized and reconciled with its schema")
}
