-- name: CreateMember :one
INSERT INTO members (
  membername,
  password_hash,
  name_entire,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetMember :one
SELECT * FROM members
WHERE membername = $1 LIMIT 1;

-- name: UpdateMember :one
UPDATE members
SET
  password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
  password_changed_time = COALESCE(sqlc.narg(password_changed_time), password_changed_time),
  name_entire = COALESCE(sqlc.narg(name_entire), name_entire),
  email = COALESCE(sqlc.narg(email), email)
WHERE
  membername = sqlc.arg(membername)
RETURNING *;
