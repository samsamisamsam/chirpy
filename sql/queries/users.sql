-- name: CreateUser :one
insert into USERS(ID, CREATED_AT, UPDATED_AT, EMAIL)
values(gen_random_uuid(), NOW(), NOW(), $1)
returning *;
