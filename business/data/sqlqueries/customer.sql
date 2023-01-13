-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAllAccounts :many
SELECT *
FROM accounts
ORDER BY created_on;

-- name: CreateAccount :one
INSERT INTO accounts (full_name, email, username, password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateAccount :many
UPDATE accounts
SET full_name = $2,
    email     = $3,
    username  = $4,
    password  = $5
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE
FROM accounts
WHERE id = $1;

-- name: AddState :one
INSERT INTO address_states (name)
VALUES ($1)
RETURNING *;

-- name: AddAddressToAccount :exec
INSERT INTO addresses (account_id, address_1, address_2,
                       city, state_id, zip_code, recipient_name)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetAccountAddresses :many
SELECT addresses.*, address_states.name AS state_name
FROM addresses
         LEFT JOIN address_states ON addresses.state_id = address_states.id
WHERE account_id = $1;

-- name: DeleteAllAccounts :exec
DELETE
FROM accounts;