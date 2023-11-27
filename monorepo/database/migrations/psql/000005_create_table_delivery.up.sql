BEGIN;

CREATE TABLE delivery (
    order_uid UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    zip VARCHAR(20) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address VARCHAR(200) NOT NULL,
    region VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,

    CONSTRAINT pk_order_id PRIMARY KEY (order_uid),
    CONSTRAINT fk_order_uid FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE RESTRICT
);

COMMIT;