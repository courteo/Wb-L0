create table if not exists orders (
    order_id varchar(10) not null unique,
    data json not null
);
