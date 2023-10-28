-- name: CreateDetail :one
INSERT INTO details (
  from_trader_id,
  to_trader_id,
  number
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetDetail :one
SELECT * FROM details
WHERE id = $1 LIMIT 1;

-- name: ListDetails :many
SELECT * FROM details
WHERE 
    from_trader_id = $1 OR
    to_trader_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;
