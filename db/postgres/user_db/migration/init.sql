create table IF NOT EXISTS users
(
    id uuid not null primary key,
    name text not null unique ,
    password TEXT not null  ,
    email TEXT not null unique ,
    Balance int not null ,
    created_at timestamptz not null ,
    updated_at timestamptz not null
);