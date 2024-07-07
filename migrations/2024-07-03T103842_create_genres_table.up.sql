CREATE TABLE IF NOT EXISTS books.genres
(
    id           UUID                  DEFAULT gen_random_uuid() PRIMARY KEY,
    name         VARCHAR(256) NOT NULL UNIQUE,
    description  TEXT         NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);