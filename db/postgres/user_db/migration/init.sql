create table IF NOT EXISTS users
(
    id varchar not null primary key,
    name varchar not null unique ,
    Balance int not null ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp
);