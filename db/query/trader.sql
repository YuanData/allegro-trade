-- name: CreateTrader :one
INSERT INTO traders (
  holder,
  rest,
  symbol
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTrader :one
SELECT * FROM traders
WHERE id = $1 LIMIT 1;

-- name: ListTraders :many
SELECT * FROM traders
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTrader :one
UPDATE traders
SET rest = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTrader :exec
DELETE FROM traders
WHERE id = $1;
