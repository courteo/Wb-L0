create table if not exists users (
    order_id varchar(10) not null unique,
    data json not null
);
