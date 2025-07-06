-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR PRIMARY KEY,
    track_number VARCHAR NOT NULL,
    entry VARCHAR NOT NULL,
    locale VARCHAR NOT NULL,
    internal_signature VARCHAR,
    customer_id VARCHAR NOT NULL,
    delivery_service VARCHAR NOT NULL,
    shardkey VARCHAR NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR PRIMARY KEY REFERENCES orders(order_uid),
    name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    zip VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    region VARCHAR NOT NULL,
    email VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS payment (
    order_uid VARCHAR PRIMARY KEY REFERENCES orders(order_uid),
    transaction VARCHAR NOT NULL,
    request_id VARCHAR,
    currency VARCHAR NOT NULL,
    provider VARCHAR NOT NULL,
    amount INT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL
);

CREATE TABLE IF NOT EXISTS item (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR REFERENCES orders(order_uid),
    chrt_id INT NOT NULL,
    track_number VARCHAR NOT NULL,
    price INT NOT NULL,
    rid VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    sale INT NOT NULL,
    size VARCHAR NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR NOT NULL,
    status INT NOT NULL
);

CREATE INDEX idx_items_order_uid ON item(order_uid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS delivery;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
