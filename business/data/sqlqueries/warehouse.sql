-- name: GetWarehouses :many
SELECT * FROM warehouses;

-- name: CreateWarehouse :exec
INSERT INTO warehouses (
    name
) VALUES (
             $1
         );

-- name: DeleteAllWarehouses :exec
DELETE FROM warehouses;