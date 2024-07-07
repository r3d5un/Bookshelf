package data

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
)

type Book struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	AuthorID    *uuid.UUID `json:"authorId"`
	Description *string    `json:"description"`
	Published   *time.Time `json:"published"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type BookModel struct {
	DB      *sql.DB
	Timeout *time.Duration
}

func (m *BookModel) Get(ctx context.Context, id uuid.UUID) (b *Book, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       title,
       author_id,
       description,
       published,
       created_at,
       updated_at
FROM books.books
WHERE id = $1;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		"query",
		slog.String("statement", database.MinifySQL(query)),
		slog.String("id", id.String()),
	)

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&b.ID,
		&b.Title,
		&b.AuthorID,
		&b.Description,
		&b.Published,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Info("no rows found", "group_id", id.String())
			return nil, ErrRecordNotFound
		default:
			logger.Info("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning book")
	return b, nil
}

func (m *BookModel) GetAll(
	ctx context.Context,
	filters Filters,
) (books []*Book, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       title,
       author_id,
       description,
       published,
       created_at,
       updated_at
FROM books.books
WHERE ($1 IS NULL OR id = $1)
  AND ($2 IS NULL OR author_id = $2)
  AND ($3 IS NULL OR description LIKE '%' || $3 || '%')
  AND ($4 IS NULL OR published >= $4)
  AND ($5 IS NULL OR published < $5)
  AND ($6 IS NULL OR created_at >= $6)
  AND ($7 IS NULL OR created_at < $7)
  AND ($8 IS NULL OR updated_at >= $8)
  AND ($9 IS NULL OR updated_at < $9)
` + database.CreateOrderByClause(filters.OrderBy) + `
OFFSET $10 FETCH NEXT $11 ROWS ONLY;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		"query",
		slog.String("statement", database.MinifySQL(query)),
		"filters", filters,
	)

	logger.Info("performing query")
	rows, err := m.DB.QueryContext(
		qCtx,
		query,
		filters.ID,
		filters.AuthorID,
		filters.Description,
		filters.PublishedFrom,
		filters.PublishedTo,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
		filters.offset(),
		filters.limit(),
	)
	if err != nil {
		logger.Error("error performing query", "error", err)
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.AuthorID,
			&book.Description,
			&book.Published,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		books = append(books, &book)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, nil, err
	}
	numberOfRecords := len(books)

	logger.Info("returning records", slog.Int("records", numberOfRecords))
	return books, &numberOfRecords, nil
}

func (m *BookModel) Insert(ctx context.Context, newBook Book) (b *Book, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.books (id,
                         title,
                         author_id,
                         description,
                         published,
                         created_at,
                         updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7)
RETURNING
    id,
    title,
    author_id,
    description,
    published,
    created_at,
    updated_at;
`
	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newBook", newBook,
		),
	)

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newBook.ID,
		newBook.Title,
		newBook.AuthorID,
		newBook.Description,
		newBook.Published,
		newBook.CreatedAt,
		newBook.UpdatedAt,
	).Scan(
		&b.ID,
		&b.Title,
		&b.AuthorID,
		&b.Description,
		&b.Published,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		logger.Error("unable to insert record", "error", err)
		return nil, err
	}

	logger.Info("returning inserted book", "insertedBook", b)
	return b, nil
}

func (m *BookModel) Update(ctx context.Context, newBook Book) (b *Book, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
UPDATE books.books
SET id          = COALESCE($1, id),
    title       = COALESCE($2, title),
    author_id   = COALESCE($3, author_id),
    description = COALESCE($4, description),
    published   = COALESCE($5, published),
    created_at  = COALESCE($6, created_at),
    updated_at  = COALESCE($7, updated_at)
WHERE id = $1
RETURNING
    id,
    title,
    author_id,
    description,
    published,
    created_at,
    updated_at;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newBook", newBook,
		),
	)

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newBook.ID,
		newBook.Title,
		newBook.AuthorID,
		newBook.Description,
		newBook.Published,
		newBook.CreatedAt,
		newBook.UpdatedAt,
	).Scan(
		&b.ID,
		&b.Title,
		&b.AuthorID,
		&b.Description,
		&b.Published,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Info("no record found", "error", err)
			return nil, ErrRecordNotFound
		default:
			logger.Error("unable to perform query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning updated book", "updatedBook", b)
	return b, nil
}

func (m *BookModel) Upsert(ctx context.Context, newBook Book) (b *Book, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.books (id,
                         title,
                         author_id,
                         description,
                         published,
                         created_at,
                         updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8)
ON CONFLICT (id)
    DO UPDATE SET id          = excluded.id,
                  title       = excluded.title,
                  author_id   = excluded.author_id,
                  description = excluded.description,
                  published   = excluded.published,
                  created_at  = excluded.created_at,
                  updated_at  = excluded.updated_at;
`

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newBook", newBook,
		),
	)

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newBook.ID,
		newBook.Title,
		newBook.AuthorID,
		newBook.Description,
		newBook.Published,
		newBook.CreatedAt,
		newBook.UpdatedAt,
	).Scan(
		&b.ID,
		&b.Title,
		&b.AuthorID,
		&b.Description,
		&b.Published,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		logger.Info("an error occurred while executing query", "error", err)
		return nil, err
	}
	logger.Info("returning upserted book")
	return b, nil
}

func (m *BookModel) Delete(ctx context.Context, id uuid.UUID) (b *Book, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
DELETE FROM books.books
WHERE id = $1
RETURNING
	id,
	title,
	author_id,
	description,
	published,
	created_at,
	updated_at;
`
	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		"query",
		slog.String("statement", database.MinifySQL(query)),
		slog.String("id", id.String()),
	)

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&b.ID,
		&b.Title,
		&b.AuthorID,
		&b.Description,
		&b.Published,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Info("no rows found", "group_id", id.String())
			return nil, ErrRecordNotFound
		default:
			logger.Info("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning deleted book")
	return b, nil
}
