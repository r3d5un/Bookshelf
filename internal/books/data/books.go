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
	Description *string    `json:"description,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
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
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			slog.String("id", id.String()),
		),
	)

	b = &Book{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&b.ID,
		&b.Title,
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
       description,
       published,
       created_at,
       updated_at
FROM books.books
WHERE ($1::uuid IS NULL OR id = $1::uuid)
  AND ($2::text = '' OR title LIKE '%' || $2::text || '%')
  AND ($3::text = '' OR description LIKE '%' || $3::text || '%')
  AND ($4::timestamp IS NULL OR published >= $4::timestamp)
  AND ($5::timestamp IS NULL OR published < $5::timestamp)
  AND ($6::timestamp IS NULL OR created_at >= $6::timestamp)
  AND ($7::timestamp IS NULL OR created_at < $7::timestamp)
  AND ($8::timestamp IS NULL OR updated_at >= $8::timestamp)
  AND ($9::timestamp IS NULL OR updated_at < $9::timestamp)
` + database.CreateOrderByClause(filters.OrderBy) + `
OFFSET $10 FETCH NEXT $11 ROWS ONLY;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"filters", filters,
		),
	)

	books = []*Book{}

	logger.Info("performing query")
	rows, err := m.DB.QueryContext(
		qCtx,
		query,
		filters.ID,
		filters.Title,
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
                         description,
                         published,
                         created_at,
                         updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        NOW(),
        NOW())
RETURNING
    id,
    title,
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

	b = &Book{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newBook.ID,
		newBook.Title,
		newBook.Description,
		newBook.Published,
	).Scan(
		&b.ID,
		&b.Title,
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
SET title       = CASE WHEN $2 = '' THEN title ELSE COALESCE($2, title) END,
    description = COALESCE($3, description),
    published   = COALESCE($4, published),
    created_at  = COALESCE($5, created_at),
    updated_at  = NOW()
WHERE id = $1
RETURNING
    id,
    title,
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

	b = &Book{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newBook.ID,
		newBook.Title,
		newBook.Description,
		newBook.Published,
		newBook.CreatedAt,
	).Scan(
		&b.ID,
		&b.Title,
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
                         description,
                         published,
                         created_at,
                         updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        NOW())
ON CONFLICT (id)
    DO UPDATE SET id          = excluded.id,
                  title       = excluded.title,
                  description = excluded.description,
                  published   = excluded.published,
                  created_at  = excluded.created_at,
                  updated_at  = excluded.updated_at
RETURNING id,
          title,
          description,
          published,
          created_at,
          updated_at;
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

	b = &Book{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newBook.ID,
		newBook.Title,
		newBook.Description,
		newBook.Published,
		newBook.CreatedAt,
	).Scan(
		&b.ID,
		&b.Title,
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
			slog.String("id", id.String()),
		),
	)

	b = &Book{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&b.ID,
		&b.Title,
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

func (m *BookModel) GetByAuthorID(
	ctx context.Context,
	id uuid.UUID,
) (books []*Book, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT b.id,
       b.title,
       b.description,
       b.published,
       b.created_at,
       b.updated_at
FROM books.books b
         INNER JOIN
     books.book_authors ba ON b.id = ba.author_id
         INNER JOIN
     books.authors a ON a.id = ba.book_id
WHERE a.id = $1
ORDER BY a.id;
`
	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"id", id.String(),
		),
	)

	books = []*Book{}

	logger.Info("performing query")
	rows, err := m.DB.QueryContext(
		qCtx,
		query,
		id,
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

func (m *BookModel) GetBySeriesID(
	ctx context.Context,
	id uuid.UUID,
) (books []*Book, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT b.id,
       b.title,
       b.description,
       b.published,
       b.created_at,
       b.updated_at
FROM books.books b
         INNER JOIN
     books.book_series bs ON b.id = bs.book_id
         INNER JOIN
     books.series s ON s.id = bs.series_id
WHERE s.id = $1
ORDER BY bs.series_order;
`
	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"id", id.String(),
		),
	)

	books = []*Book{}

	logger.Info("performing query")
	rows, err := m.DB.QueryContext(
		qCtx,
		query,
		id,
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

func (m *BookModel) GetByGenreID(
	ctx context.Context,
	id uuid.UUID,
) (books []*Book, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT b.id,
       b.title,
       b.description,
       b.published,
       b.created_at,
       b.updated_at
FROM books.books b
         INNER JOIN
     books.book_genres bg ON b.id = bg.book_id
         INNER JOIN
     books.genres s ON s.id = bg.genres_id
WHERE s.id = $1
ORDER BY b.published;
`
	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"id", id.String(),
		),
	)

	books = []*Book{}

	logger.Info("performing query")
	rows, err := m.DB.QueryContext(
		qCtx,
		query,
		id,
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
