create table IF NOT EXISTS orders
(
    id varchar not null primary key ,
    user_id varchar not null ,
    total_price int not null ,
    total_paid int not null ,
    total_return int not null ,
    receipt_code varchar not null  unique ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp
);

create table IF NOT EXISTS order_products
(
    id varchar not null primary key ,
    order_id varchar not null ,
    product_id varchar not null  ,
    total_price int not null ,
    qty int not null ,
    created_at timestamp not null ,
    updated_at timestamp not null ,
    deleted_at timestamp,
    foreign key(order_id) references orders(id)
    );