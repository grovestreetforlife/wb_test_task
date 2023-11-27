BEGIN;

CREATE TABLE product (
    -- chrt_id
    chrt_id BIGINT NOT NULL,

    --track_number номер заказа
    track_number VARCHAR(500) NOT NULL,

    -- price цена без учета скидки
    price DECIMAL NOT NULL,

    -- rid
    rid VARCHAR(500) NOT NULL,

    -- name наименование товара
    name VARCHAR(600) NOT NULL,

    -- sale скидка
    sale INTEGER NOT NULL,

    -- size размер
    size VARCHAR(200) NOT NULL,

    -- total_price цена с учетом скидки
    total_price DECIMAL NOT NULL,

    -- nm_id id товара
    nm_id BIGINT NOT NULL,

    -- brand брэнд товара
    brand VARCHAR(300) NOT NULL,

    -- status стутус обработки заказа
    status INTEGER NOT NULL,

    CONSTRAINT fk_track_number FOREIGN KEY (track_number) REFERENCES orders(track_number) ON DELETE RESTRICT
);

COMMIT;