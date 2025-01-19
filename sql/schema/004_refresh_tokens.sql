-- +goose Up
create table REFRESH_TOKENS (
    TOKEN text primary key,
    CREATED_AT timestamp not null,
    UPDATED_AT timestamp not null,
    USER_ID uuid not null references USERS (ID) on delete cascade,
    EXPIRES_AT timestamp not null,
    REVOKED_AT timestamp
);

-- +goose Down
drop table REFRESH_TOKENS;
