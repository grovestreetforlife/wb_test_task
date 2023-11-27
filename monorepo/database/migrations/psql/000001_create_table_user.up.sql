BEGIN;

CREATE TABLE users (
    -- id пользователя
    id VARCHAR(500) NOT NULL,

    CONSTRAINT pk_id PRIMARY KEY (id)
);

COMMIT;