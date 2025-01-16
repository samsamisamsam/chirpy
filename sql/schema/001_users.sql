-- +goose Up
create table USERS (
    ID uuid primary key,
    CREATED_AT timestamp not null,
    UPDATED_AT timestamp not null,
    EMAIL text unique not null
);

-- +goose Down
drop table USERS;
