package main

import (
	"atlasExample/db"
	_ "embed"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	prodDBConfig := db.DBConfig{
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

	psqlDB, err := db.New(prodDBConfig, atlasDevDBConfig)
	if err != nil {
		panic("Cannot open the databases: " + err.Error())
	}

	err = psqlDB.ReconcileWithAtlasSQLSchema(db.BuiltInSQLSchema)
	if err != nil {
		panic("Database cannot be reconciled with its Atlas Schema: " + err.Error())
	}

	fmt.Println("Database has been initialized and reconciled with its schema")
}
