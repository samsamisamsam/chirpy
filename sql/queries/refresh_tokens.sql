-- name: StoreRefreshToken :exec
insert into REFRESH_TOKENS(TOKEN, CREATED_AT, UPDATED_AT, USER_ID, EXPIRES_AT, REVOKED_AT)
values($1, NOW(), NOW(), $2, NOW() + interval '60 days', NULL);

-- name: GetUserWithToken :one
select TOKEN, CREATED_AT, UPDATED_AT, USER_ID, EXPIRES_AT, REVOKED_AT
from REFRESH_TOKENS
where TOKEN = $1;

-- name: RevokeToken :exec
update REFRESH_TOKENS
set UPDATED_AT = NOW(),
REVOKED_AT = NOW()
where TOKEN = $1;
