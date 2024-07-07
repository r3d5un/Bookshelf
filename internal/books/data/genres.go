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
	Name        string     `json:"name"`
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
		"query",
		slog.String("statement", database.MinifySQL(query)),
		slog.String("id", id.String()),
	)

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

	logger.Info("returning author")
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
WHERE ($1 IS NULL OR id = $1)
  AND ($2 IS NULL OR name LIKE '%' || $2 || '%')
  AND ($3 IS NULL OR description LIKE '%' || $3 || '%')
  AND ($4 IS NULL OR created_at >= $4)
  AND ($5 IS NULL OR created_at < $5)
  AND ($6 IS NULL OR updated_at >= $6)
  AND ($7 IS NULL OR updated_at < $7)
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
