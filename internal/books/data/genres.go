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

type Genre struct {
	ID          uuid.UUID  `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type GenreModel struct {
	DB      *sql.DB
	Timeout *time.Duration
}

func (m *GenreModel) Get(ctx context.Context, id uuid.UUID) (genre *Genre, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       description,
       created_at,
       updated_at
FROM books.genres
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

	genre = &Genre{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&genre.ID,
		&genre.Name,
		&genre.Description,
		&genre.CreatedAt,
		&genre.UpdatedAt,
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

	logger.Info("returning genre")
	return genre, nil
}

func (m *GenreModel) GetAll(
	ctx context.Context,
	filters Filters,
) (genres []*Genre, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       description,
       created_at,
       updated_at
FROM books.genres
WHERE ($1::uuid IS NULL OR id = $1::uuid)
  AND ($2::text = '' OR name LIKE '%' || $2::text || '%')
  AND ($3::text = '' OR description LIKE '%' || $3::text || '%')
  AND ($4::timestamp IS NULL OR created_at >= $4::timestamp)
  AND ($5::timestamp IS NULL OR created_at < $5::timestamp)
  AND ($6::timestamp IS NULL OR updated_at >= $6::timestamp)
  AND ($7::timestamp IS NULL OR updated_at < $7::timestamp)
` + database.CreateOrderByClause(filters.OrderBy) + `
OFFSET $8 FETCH NEXT $9 ROWS ONLY;
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

	genres = []*Genre{}

	logger.Info("performing query")
	rows, err := m.DB.QueryContext(
		qCtx,
		query,
		filters.ID,
		filters.Name,
		filters.Description,
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
		var genre Genre

		err := rows.Scan(
			&genre.ID,
			&genre.Name,
			&genre.Description,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		genres = append(genres, &genre)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, nil, err
	}
	numberOfRecords := len(genres)

	logger.Info("returning records", slog.Int("records", numberOfRecords))
	return genres, &numberOfRecords, nil
}

func (m *GenreModel) Insert(ctx context.Context, newGenre Genre) (genre *Genre, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.genres (id,
                          name,
                          description,
                          created_at,
                          updated_at)
VALUES ($1,
        $2,
        $3,
        NOW(),
        NOW())
RETURNING
    id,
    name,
    description,
    created_at,
    updated_at;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newGenre", newGenre,
		),
	)

	genre = &Genre{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newGenre.ID,
		newGenre.Name,
		newGenre.Description,
	).Scan(
		&genre.ID,
		&genre.Name,
		&genre.Description,
		&genre.CreatedAt,
		&genre.UpdatedAt,
	)
	if err != nil {
		logger.Error("unable to insert record", "error", err)
		return nil, err
	}

	logger.Info("returning inserted genre", "insertedGenre", genre)
	return genre, nil
}

func (m *GenreModel) Update(ctx context.Context, newGenre Genre) (genre *Genre, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
UPDATE books.genres
SET id          = COALESCE($1, id),
    name        = COALESCE($2, name),
    description = COALESCE($3, description),
    created_at  = COALESCE($4, created_at),
    updated_at  = NOW()
WHERE id = $1
RETURNING
    id,
    name,
    description,
    created_at,
    updated_at;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newGenre", newGenre,
		),
	)

	genre = &Genre{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newGenre.ID,
		newGenre.Name,
		newGenre.Description,
		newGenre.CreatedAt,
	).Scan(
		&genre.ID,
		&genre.Name,
		&genre.Description,
		&genre.CreatedAt,
		&genre.UpdatedAt,
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

	logger.Info("returning updated genre", "updatedGenre", genre)
	return genre, nil
}

func (m *GenreModel) Upsert(ctx context.Context, newGenre Genre) (genre *Genre, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.genres (id,
                          name,
                          description,
                          created_at,
                          updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        NOW())
ON CONFLICT (id)
    DO UPDATE SET id          = excluded.id,
                  name        = excluded.name,
                  description = excluded.description,
                  created_at  = excluded.created_at,
                  updated_at  = excluded.updated_at
RETURNING id,
          name,
          description,
          created_at,
          updated_at;
`

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newGenre", newGenre,
		),
	)

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	genre = &Genre{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newGenre.ID,
		newGenre.Name,
		newGenre.Description,
		newGenre.CreatedAt,
	).Scan(
		&genre.ID,
		&genre.Name,
		&genre.Description,
		&genre.CreatedAt,
		&genre.UpdatedAt,
	)
	if err != nil {
		logger.Info("an error occurred while executing query", "error", err)
		return nil, err
	}

	logger.Info("returning upserted genre")
	return genre, nil
}

func (m *GenreModel) Delete(ctx context.Context, id uuid.UUID) (genre *Genre, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
DELETE FROM books.genres
WHERE id = $1
RETURNING
	id,
	name,
	description,
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

	genre = &Genre{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&genre.ID,
		&genre.Name,
		&genre.Description,
		&genre.CreatedAt,
		&genre.UpdatedAt,
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

	logger.Info("returning deleted genre")
	return genre, nil
}

func (m *GenreModel) GetByBookID(
	ctx context.Context,
	id uuid.UUID,
) (genres []*Genre, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT g.id,
       g.name,
       g.description,
       g.created_at,
       g.updated_at
FROM books.genres g
         INNER JOIN
     books.book_genres bg ON g.id = bg.genres_id
         INNER JOIN
     books.books b ON b.id = bg.book_id
WHERE b.id = $1
ORDER BY b.id;
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

	genres = []*Genre{}

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
		var genre Genre

		err := rows.Scan(
			&genre.ID,
			&genre.Name,
			&genre.Description,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		genres = append(genres, &genre)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, nil, err
	}
	numberOfRecords := len(genres)

	logger.Info("returning records", slog.Int("records", numberOfRecords))
	return genres, &numberOfRecords, nil
}
