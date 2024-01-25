CREATE DATABASE IF NOT EXISTS defaultdb;
set database = defaultdb;
CREATE TABLE if not exists transactions  (
                                             id UUID primary key,
                                             customer_id UUID not null,
                                             product_id UUID not null,
                                             quantity int not null,
                                             amount double precision not null,
                                             created timestamp default current_timestamp
)