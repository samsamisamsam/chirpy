-- +goose Up
alter table USERS
add column HASHED_PASSWORD text not null default 'unset';

-- +goose Down
alter table USERS
drop column HASHED_PASSWORD;
