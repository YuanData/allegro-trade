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
