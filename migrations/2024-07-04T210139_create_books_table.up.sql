CREATE TABLE IF NOT EXISTS books.books
(
    id          UUID                  DEFAULT gen_random_uuid() PRIMARY KEY,
    title       VARCHAR(512) NOT NULL,
    description TEXT         NULL,
    published   TIMESTAMP    NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
