-- +goose Up
alter table USERS
add column IS_CHIRPY_RED bool not null default false;

-- +goose Down
alter table USERS
drop column IS_CHIRPY_RED;
