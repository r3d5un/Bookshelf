CREATE TABLE IF NOT EXISTS books.authors
(
    id          UUID                  DEFAULT gen_random_uuid() PRIMARY KEY,
    name        VARCHAR(512) NOT NULL,
    description TEXT         NULL,
    website     varchar(256) NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);
