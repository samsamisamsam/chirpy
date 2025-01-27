-- name: CreateChirp :one
insert into CHIRPS(ID, CREATED_AT, UPDATED_AT, BODY, USER_ID)
values(gen_random_uuid(), NOW(), NOW(), $1, $2)
returning *;

-- name: GetAllChirps :many
select ID, CREATED_AT, UPDATED_AT, BODY, USER_ID
from CHIRPS
order by CREATED_AT;

-- name: GetChirp :one
select ID, CREATED_AT, UPDATED_AT, BODY, USER_ID
from CHIRPS
where ID = $1;

-- name: GetUserIDFromChirpID :one
select USER_ID
from CHIRPS
where ID = $1;

-- name: DeleteChirp :exec
delete from CHIRPS
where ID = $1;
