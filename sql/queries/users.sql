-- name: CreateUser :one
insert into USERS(ID, CREATED_AT, UPDATED_AT, EMAIL, HASHED_PASSWORD)
values(gen_random_uuid(), NOW(), NOW(), $1, $2)
returning ID, CREATED_AT, UPDATED_AT, EMAIL, IS_CHIRPY_RED;

-- name: GetUserByEmail :one
select ID, CREATED_AT, UPDATED_AT, EMAIL, HASHED_PASSWORD, IS_CHIRPY_RED
from USERS
where EMAIL = $1;

-- name: UpdateUser :one
update USERS
set EMAIL = $1,
HASHED_PASSWORD = $2,
UPDATED_AT = NOW()
where ID = $3
returning ID, CREATED_AT, UPDATED_AT, EMAIL, IS_CHIRPY_RED;

-- name: UpgradeUser :exec
update USERS
set IS_CHIRPY_RED = true
where ID = $1;
