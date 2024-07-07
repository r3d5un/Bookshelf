CREATE TABLE IF NOT EXISTS books.book_authors
(
    book_id UUID NOT NULL REFERENCES books.books (id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES books.authors (id) ON DELETE CASCADE
)