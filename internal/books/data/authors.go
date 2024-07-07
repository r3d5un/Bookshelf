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

type Author struct {
	ID          uuid.UUID  `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Website     *string    `json:"website"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type AuthorModel struct {
	DB      *sql.DB
	Timeout *time.Duration
}

func (m *AuthorModel) Get(ctx context.Context, id uuid.UUID) (author *Author, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       description,
       website,
       created_at,
       updated_at
FROM books.authors
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
		&author.ID,
		&author.Name,
		&author.Description,
		&author.Website,
		&author.CreatedAt,
		&author.UpdatedAt,
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

	logger.Info("returning author")
	return author, nil
}

func (m *AuthorModel) GetAll(
	ctx context.Context,
	filters Filters,
) (authors []*Author, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       description,
       website,
       created_at,
       updated_at
FROM books.authors
WHERE ($1 IS NULL OR id = $1)
  AND ($2 IS NULL OR name LIKE '%' || $2 || '%')
  AND ($3 IS NULL OR description LIKE '%' || $3 || '%')
  AND ($4 IS NULL OR website LIKE '%' || $4 || '%')
  AND ($5 IS NULL OR created_at >= $5)
  AND ($6 IS NULL OR created_at < $6)
  AND ($7 IS NULL OR updated_at >= $7)
  AND ($8 IS NULL OR updated_at < $8)
` + database.CreateOrderByClause(filters.OrderBy) + `
OFFSET $9 FETCH NEXT $10 ROWS ONLY;
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
		filters.Name,
		filters.Description,
		filters.Website,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
	)
	if err != nil {
		logger.Error("error performing query", "error", err)
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var author Author

		err := rows.Scan(
			&author.ID,
			&author.Name,
			&author.Description,
			&author.Website,
			&author.CreatedAt,
			&author.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		authors = append(authors, &author)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, nil, err
	}
	numberOfRecords := len(authors)

	logger.Info("returning records", slog.Int("records", numberOfRecords))
	return authors, &numberOfRecords, nil
}

func (m *AuthorModel) Insert(ctx context.Context, newAuthor Author) (author *Author, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.authors (id,
                           name,
                           description,
                           website,
                           created_at,
                           updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6)
RETURNING
    id,
    name,
    description,
    website,
    created_at,
    updated_at;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newAuthor", newAuthor,
		),
	)

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newAuthor.ID,
		newAuthor.Name,
		newAuthor.Description,
		newAuthor.Website,
		newAuthor.CreatedAt,
		newAuthor.UpdatedAt,
	).Scan(
		&author.ID,
		&author.Name,
		&author.Description,
		newAuthor.Website,
		newAuthor.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		logger.Error("unable to insert record", "error", err)
		return nil, err
	}

	logger.Info("returning inserted author", "insertedAuthor", author)
	return author, nil
}

func (m *AuthorModel) Update(ctx context.Context, newAuthor Author) (author *Author, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
UPDATE books.authors
SET id          = COALESCE($1, id),
    name        = COALESCE($2, name),
    description = COALESCE($3, description),
    website     = COALESCE($4, website),
    created_at  = COALESCE($5, created_at),
    updated_at  = COALESCE($6, updated_at)
WHERE id = $1
RETURNING
    id,
    name,
    description,
    website,
    created_at,
    updated_at;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newAuthor", newAuthor,
		),
	)

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newAuthor.ID,
		newAuthor.Name,
		newAuthor.Description,
		newAuthor.Website,
		newAuthor.CreatedAt,
		newAuthor.UpdatedAt,
	).Scan(
		&author.ID,
		&author.Name,
		&author.Description,
		&author.Website,
		&author.CreatedAt,
		&author.UpdatedAt,
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

	logger.Info("returning updated author", "updatedAuthor", author)
	return author, nil
}

func (m *AuthorModel) Upsert(ctx context.Context, newAuthor Author) (author *Author, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.authors (id,
                          name,
                          description,
                          website,
                          created_at,
                          updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6)
ON CONFLICT (id)
    DO UPDATE SET id          = excluded.id,
                  name        = excluded.name,
                  description = excluded.description,
                  website     = excluded.website,
                  created_at  = excluded.created_at,
                  updated_at  = excluded.updated_at;
`

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newAuthor", newAuthor,
		),
	)

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newAuthor.ID,
		newAuthor.Name,
		newAuthor.Description,
		newAuthor.Website,
		newAuthor.CreatedAt,
		newAuthor.UpdatedAt,
	).Scan(
		&author.ID,
		&author.Name,
		&author.Description,
		&author.Website,
		&author.CreatedAt,
		&author.UpdatedAt,
	)
	if err != nil {
		logger.Info("an error occurred while executing query", "error", err)
		return nil, err
	}

	logger.Info("returning upserted author")
	return author, nil
}

func (m *AuthorModel) Delete(ctx context.Context, id uuid.UUID) (author *Author, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
DELETE FROM books.authors
WHERE id = $1
RETURNING
	id,
	name,
	description,
	website,
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
		&author.ID,
		&author.Name,
		&author.Description,
		&author.Website,
		&author.CreatedAt,
		&author.UpdatedAt,
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

	logger.Info("returning deleted author")
	return author, nil
}
