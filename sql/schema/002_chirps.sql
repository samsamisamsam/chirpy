-- +goose Up
create table CHIRPS (
    ID uuid primary key,
    CREATED_AT timestamp not null,
    UPDATED_AT timestamp not null,
    BODY text not null,
    USER_ID uuid not null references USERS (ID) on delete cascade
);

-- +goose Down
drop table CHIRPS;
