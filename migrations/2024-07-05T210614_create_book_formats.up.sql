CREATE TABLE IF NOT EXISTS books.books
(
    id        UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    book_id   UUID         NOT NULL REFERENCES books.books (id),
    type      VARCHAR(256) NOT NULL,
    published TIMESTAMP    NULL,
    publisher VARCHAR(128) NULL,
    isbn      VARCHAR(32)  NULL,
    isbn10    VARCHAR(32)  NULL,
    language  VARCHAR(2)   NULL,
    pages     INT          NULL,
    duration  INTERVAL     NULL
);