// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: verify_email.sql

package db

import (
	"context"
)

const createVerifyEmail = `-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
    membername,
    email,
    secret_code
) VALUES (
    $1, $2, $3
) RETURNING id, membername, email, secret_code, is_used, created_time, expired_time
`

type CreateVerifyEmailParams struct {
	Membername string `json:"membername"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRowContext(ctx, createVerifyEmail, arg.Membername, arg.Email, arg.SecretCode)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Membername,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedTime,
		&i.ExpiredTime,
	)
	return i, err
}
