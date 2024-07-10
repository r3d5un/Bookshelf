package data

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
)

type BookGenre struct {
	BookID  uuid.UUID `json:"bookId"`
	GenreID uuid.UUID `json:"genreId"`
}

type BookGenreModel struct {
	DB      *sql.DB
	Timeout *time.Duration
}

func (m BookGenreModel) Insert(
	ctx context.Context,
	bookID uuid.UUID,
	genreID uuid.UUID,
) (bg *BookGenre, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.book_genres (book_id,
                               genres_id)
VALUES ($1,
        $2)
RETURNING book_id,
          genres_id;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			slog.String("bookId", bookID.String()),
			slog.String("genreId", genreID.String()),
		),
	)

	bg = &BookGenre{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		bookID,
		genreID,
	).Scan(
		&bg.BookID,
		&bg.GenreID,
	)

	if err != nil {
		logger.Error("unable to insert record", "error", err)
		return nil, err
	}

	logger.Info("returning inserted book genre", "insertedBook", bg)
	return bg, nil
}
