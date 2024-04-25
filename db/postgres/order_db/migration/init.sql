create table IF NOT EXISTS orders
(
    id uuid not null primary key ,
    user_id uuid not null ,
    total_price int not null ,
    total_paid int not null ,
    total_return int not null ,
    receipt_code text not null  unique ,
    created_at timestamptz not null ,
    updated_at timestamptz not null ,
    deleted_at timestamptz
);

create table IF NOT EXISTS order_products
(
    id uuid not null primary key ,
    order_id uuid not null ,
    product_id uuid not null  ,
    total_price int not null ,
    qty int not null ,
    created_at timestamptz not null ,
    updated_at timestamptz not null ,
    deleted_at timestamptz,
    foreign key(order_id) references orders(id)
);
