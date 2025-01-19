// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
insert into USERS(ID, CREATED_AT, UPDATED_AT, EMAIL, HASHED_PASSWORD)
values(gen_random_uuid(), NOW(), NOW(), $1, $2)
returning ID, CREATED_AT, UPDATED_AT, EMAIL, IS_CHIRPY_RED
`

type CreateUserParams struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

type CreateUserRow struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select ID, CREATED_AT, UPDATED_AT, EMAIL, HASHED_PASSWORD, IS_CHIRPY_RED
from USERS
where EMAIL = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
update USERS
set EMAIL = $1,
HASHED_PASSWORD = $2,
UPDATED_AT = NOW()
where ID = $3
returning ID, CREATED_AT, UPDATED_AT, EMAIL, IS_CHIRPY_RED
`

type UpdateUserParams struct {
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	ID             uuid.UUID `json:"id"`
}

type UpdateUserRow struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (UpdateUserRow, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.Email, arg.HashedPassword, arg.ID)
	var i UpdateUserRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.IsChirpyRed,
	)
	return i, err
}

const upgradeUser = `-- name: UpgradeUser :exec
update USERS
set IS_CHIRPY_RED = true
where ID = $1
`

func (q *Queries) UpgradeUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgradeUser, id)
	return err
}
