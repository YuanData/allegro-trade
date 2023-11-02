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

-- name: GetTraderForUpdate :one
SELECT * FROM traders
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

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

-- name: AddTraderRest :one
UPDATE traders
SET rest = rest + sqlc.arg(number)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteTrader :exec
DELETE FROM traders
WHERE id = $1;
