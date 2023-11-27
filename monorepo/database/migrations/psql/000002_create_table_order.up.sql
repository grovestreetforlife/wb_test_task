BEGIN;

CREATE TYPE locale_type AS ENUM (
    'en'
);

CREATE TABLE orders (
    -- order_uid id заказа
    order_uid UUID NOT NULL,

    -- track_number номер отслеживания
    track_number VARCHAR(500) NOT NULL UNIQUE,

    --entry платформа
    entry VARCHAR(500) NOT NULL,

    -- locale
    locale locale_type NOT NULL,

    internal_signature VARCHAR(500),

    -- customer_id id покупателя
    customer_id VARCHAR(500) NOT NULL,

    -- delivery_service сервис доставки
    delivery_service VARCHAR(100) NOT NULL,

    -- shardkey
    shardkey VARCHAR(20) NOT NULL,

    -- sm_id
    sm_id BIGINT NOT NULL,

    -- date_created дата создания
    date_created TIMESTAMP WITHOUT TIME ZONE,

    -- oof_shard
    oof_shard VARCHAR(20) NOT NULL,

    CONSTRAINT pk_order_uid PRIMARY KEY (order_uid),
    CONSTRAINT fk_customer_id FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE RESTRICT
);

COMMIT;