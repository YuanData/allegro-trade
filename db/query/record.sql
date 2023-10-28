-- name: CreateRecord :one
INSERT INTO records (
  trader_id,
  number
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetRecord :one
SELECT * FROM records
WHERE id = $1 LIMIT 1;

-- name: ListRecords :many
SELECT * FROM records
WHERE trader_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
