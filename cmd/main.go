package main

import (
	"atlasExample/business/data"
	"atlasExample/business/data/dbinterfaces/customer"
	"atlasExample/business/data/dbinterfaces/warehouse"
	"atlasExample/business/data/dbschema"
	"atlasExample/foundation/db"
	"context"
	"database/sql"
	"fmt"
)

func initdb() (*sql.DB, error) {
	appDBConfig := db.DBConfig{
		User:       "app",
		Host:       "localhost",
		Name:       "example",
		Schema:     "shop",
		DisableTLS: true,
	}

	atlasDevDBConfig := db.DBConfig{
		User:       "app",
		Host:       "localhost",
		Name:       "example_atlas_dev",
		DisableTLS: true,
	}

	appDB, err := db.OpenDB(appDBConfig)
	if err != nil {
		panic("cannot open the database: " + err.Error())
	}

	atlasDevDB, err := db.OpenDB(atlasDevDBConfig)
	if err != nil {
		panic("cannot open the database: " + err.Error())
	}

	atlas := data.NewAtlas()
	err = atlas.ReconcileWithDeclaredSQLSchema(dbschema.SQLFiles, appDB, atlasDevDB)
	if err != nil {
		return nil, fmt.Errorf("database cannot be reconciled with its Atlas Schema: %w", err)
	}

	return appDB, nil
}

func main() {

	appDB, err := initdb()
	if err != nil {
		panic(err)
	}
	fmt.Println("database has been initialized and reconciled with its schema")

	ctx := context.Background()

	customers := customer.New(appDB)
	c, err := customers.CreateAccount(ctx, customer.CreateAccountParams{
		FullName: "John The Customer",
		Email:    "john@example.com",
		Username: "john",
		Password: "hash",
	})
	if err != nil {
		fmt.Println("Customer creation issue: ", err.Error())
	}
	fmt.Printf("Customer info: %v", c)

	s, _ := customers.AddState(ctx, "New York")
	customers.AddAddressToAccount(ctx, customer.AddAddressToAccountParams{
		AccountID:     c.ID,
		Address1:      "123 Main Street",
		Address2:      sql.NullString{},
		City:          "New York",
		StateID:       s.ID,
		ZipCode:       10015,
		RecipientName: "",
	})

	customers.UpdateAccount(ctx, customer.UpdateAccountParams{
		ID:       c.ID,
		FullName: "Moreno",
		Email:    "moreno.com",
		Username: "m",
		Password: "p",
	})

	warehouses := warehouse.New(appDB)
	err = warehouses.CreateWarehouse(ctx, "Middleborough")
	if err != nil {
		fmt.Println("Customer creation issue: ", err.Error())
	}

	fmt.Println()
	fmt.Println(customers.GetAllAccounts(ctx))
	fmt.Println(warehouses.GetWarehouses(ctx))

	err = customers.DeleteAllAccounts(ctx)
	if err != nil {
		fmt.Println(err)
	}

	warehouses.DeleteAllWarehouses(ctx)
}
