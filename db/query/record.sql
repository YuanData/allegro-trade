-- name: CreateRecord :one
INSERT INTO records (
  from_trader_id,
  to_trader_id,
  number
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetRecord :one
SELECT * FROM records
WHERE id = $1 LIMIT 1;

-- name: ListRecords :many
SELECT * FROM records
WHERE 
    from_trader_id = $1 OR
    to_trader_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;
