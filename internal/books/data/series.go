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

type Series struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type SeriesModel struct {
	DB      *sql.DB
	Timeout *time.Duration
}

func (m *SeriesModel) Get(ctx context.Context, id uuid.UUID) (series *Series, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       description,
       created_at,
       updated_at
FROM books.series
WHERE id = $1;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		"query",
		slog.String("statement", database.MinifySQL(query)),
		slog.String("id", id.String()),
	)

	series = &Series{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&series.ID,
		&series.Name,
		&series.Description,
		&series.CreatedAt,
		&series.UpdatedAt,
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

	logger.Info("returning series")
	return series, nil
}

func (m *SeriesModel) GetAll(
	ctx context.Context,
	filters Filters,
) (series []*Series, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       description,
       created_at,
       updated_at
FROM books.series
WHERE ($1::uuid IS NULL OR id = $1::uuid)
  AND ($2::text IS NULL OR name LIKE '%' || $2::text || '%')
  AND ($3::text IS NULL OR description LIKE '%' || $3::text || '%')
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
		"query",
		slog.String("statement", database.MinifySQL(query)),
		"filters", filters,
	)

	series = []*Series{}

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
		var s Series

		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		series = append(series, &s)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, nil, err
	}
	numberOfRecords := len(series)

	logger.Info("returning records", slog.Int("records", numberOfRecords))
	return series, &numberOfRecords, nil
}

func (m *SeriesModel) Insert(ctx context.Context, newSeries Series) (series *Series, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.series (id,
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
			"newSeries", newSeries,
		),
	)

	series = &Series{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newSeries.ID,
		newSeries.Name,
		newSeries.Description,
	).Scan(
		&series.ID,
		&series.Name,
		&series.Description,
		&series.CreatedAt,
		&series.UpdatedAt,
	)
	if err != nil {
		logger.Error("unable to insert record", "error", err)
		return nil, err
	}

	logger.Info("returning inserted series", "insertedSeries", series)
	return series, nil
}

func (m *SeriesModel) Update(ctx context.Context, newSeries Series) (series *Series, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
UPDATE books.series
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
			"newSeries", newSeries,
		),
	)

	series = &Series{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newSeries.ID,
		newSeries.Name,
		newSeries.Description,
		newSeries.CreatedAt,
	).Scan(
		&series.ID,
		&series.Name,
		&series.Description,
		&series.CreatedAt,
		&series.UpdatedAt,
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

	logger.Info("returning updated series", "updatedSeries", series)
	return series, nil
}

func (m *SeriesModel) Upsert(ctx context.Context, newSeries Series) (series *Series, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.series  (id,
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
			"newSeries", newSeries,
		),
	)

	series = &Series{}

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		newSeries.ID,
		newSeries.Name,
		newSeries.Description,
		newSeries.CreatedAt,
	).Scan(
		&series.ID,
		&series.Name,
		&series.Description,
		&series.CreatedAt,
		&series.UpdatedAt,
	)
	if err != nil {
		logger.Info("an error occurred while executing query", "error", err)
		return nil, err
	}

	logger.Info("returning upserted author")
	return series, nil
}

func (m *SeriesModel) Delete(ctx context.Context, id uuid.UUID) (series *Series, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
DELETE FROM books.series
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
		"query",
		slog.String("statement", database.MinifySQL(query)),
		slog.String("id", id.String()),
	)

	series = &Series{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(qCtx, query, id.String()).Scan(
		&series.ID,
		&series.Name,
		&series.Description,
		&series.CreatedAt,
		&series.UpdatedAt,
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

	logger.Info("returning deleted series")
	return series, nil
}

func (m *SeriesModel) GetByBookID(
	ctx context.Context,
	id uuid.UUID,
) (series []*Series, totalResults *int, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT s.id,
       s.name,
       s.description,
       s.created_at,
       s.updated_at
FROM books.series s
         INNER JOIN
     books.book_series bs ON s.id = bs.series_id
         INNER JOIN
     books.books b ON b.id = bs.book_id
WHERE b.id = $1
ORDER BY bs.series_order;
`
	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		"query",
		slog.String("statement", database.MinifySQL(query)),
		"bookId", id.String(),
	)

	series = []*Series{}

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
		var s Series

		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		series = append(series, &s)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, nil, err
	}
	numberOfRecords := len(series)

	logger.Info("returning records", slog.Int("records", numberOfRecords))
	return series, &numberOfRecords, nil
}
