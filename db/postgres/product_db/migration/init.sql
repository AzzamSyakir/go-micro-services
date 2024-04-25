create table IF NOT EXISTS categories
(
    id uuid not null primary key,
    name text not null unique ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp
);

create table IF NOT EXISTS products
(
    id uuid not null primary key,
    Sku text not null unique ,
    Name text not null unique ,
    stock int not null ,
    price int not null ,
    category_id uuid not null ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp,
    foreign key(category_id) references categories(id)
);