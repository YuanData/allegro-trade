-- name: CreateDetail :one
INSERT INTO details (
  trader_id,
  number
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetDetail :one
SELECT * FROM details
WHERE id = $1 LIMIT 1;

-- name: ListDetails :many
SELECT * FROM details
WHERE trader_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
