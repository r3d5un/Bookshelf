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

	return b, nil
}
