// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: warehouse.sql

package datawarehouse

import (
	"context"
)

const CreateWarehouse = `-- name: CreateWarehouse :exec
INSERT INTO warehouses (
    name
) VALUES (
             $1
         )
`

func (q *Queries) CreateWarehouse(ctx context.Context, name string) error {
	_, err := q.db.Exec(ctx, CreateWarehouse, name)
	return err
}

const DeleteAllWarehouses = `-- name: DeleteAllWarehouses :exec
DELETE FROM warehouses
`

func (q *Queries) DeleteAllWarehouses(ctx context.Context) error {
	_, err := q.db.Exec(ctx, DeleteAllWarehouses)
	return err
}

const GetWarehouses = `-- name: GetWarehouses :many
SELECT id, name FROM warehouses
`

func (q *Queries) GetWarehouses(ctx context.Context) ([]Warehouse, error) {
	rows, err := q.db.Query(ctx, GetWarehouses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Warehouse
	for rows.Next() {
		var i Warehouse
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
