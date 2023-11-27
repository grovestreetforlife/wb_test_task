CREATE TYPE currency_type AS ENUM (
    'USD',
    'RUB'
);

CREATE TYPE provider_type AS ENUM (
    'wbpay'
);

CREATE TYPE bank_type AS ENUM (
    'alpha'
);

CREATE TABLE transaction (
    id UUID PRIMARY KEY NOT NULL,

    request_id UUID,

    currency currency_type NOT NULL,

    provider provider_type NOT NULL,

    amount DECIMAL NOT NULL,

    payment_dt BIGINT NOT NULL,

    bank bank_type NOT NULL,

    delivery_cost DECIMAL NOT NULL,

    goods_total INTEGER NOT NULL,

    custom_fee INTEGER NOT NULL,

    CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES orders(order_uid) ON DELETE RESTRICT
);