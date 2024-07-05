CREATE TABLE IF NOT EXISTS books.book_genres
(
    book_id   UUID NOT NULL REFERENCES books.books (id) ON DELETE CASCADE,
    genres_id UUID NOT NULL REFERENCES books.genres (id) ON DELETE CASCADE
);
