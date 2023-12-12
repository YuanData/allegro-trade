-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    membername,
    email,
    secret_code
) VALUES (
    $1, $2, $3
) RETURNING *;
