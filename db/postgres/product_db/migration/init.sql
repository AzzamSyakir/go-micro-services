create table IF NOT EXISTS categories
(
    id varchar not null primary key,
    name varchar not null unique ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp
);

create table IF NOT EXISTS products
(
    id varchar not null primary key,
    Sku varchar not null unique ,
    Name varchar not null unique ,
    stock int not null ,
    price int not null ,
    category_id varchar not null ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp,
    foreign key(category_id) references categories(id)
);