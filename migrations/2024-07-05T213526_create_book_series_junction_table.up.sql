CREATE TABLE IF NOT EXISTS books.book_series
(
    book_id      UUID  NOT NULL REFERENCES books.books (id) ON DELETE CASCADE,
    series_id    UUID  NOT NULL REFERENCES books.authors (id) ON DELETE CASCADE,
    series_order FLOAT NULL
);