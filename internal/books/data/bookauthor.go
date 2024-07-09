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

type BookAuthor struct {
	BookID   uuid.UUID `json:"bookId"`
	AuthorID uuid.UUID `json:"authorId"`
}

type BookAuthorModel struct {
	DB      *sql.DB
	Timeout *time.Duration
}

func (m BookAuthorModel) Insert(
	ctx context.Context,
	bookID uuid.UUID,
	authorID uuid.UUID,
) (ba *BookAuthor, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.book_authors (book_id,
                                author_id)
VALUES ($1,
        $2)
RETURNING book_id,
          author_id;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			slog.String("bookId", bookID.String()),
			slog.String("authorID", authorID.String()),
		),
	)

	ba = &BookAuthor{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		bookID,
		authorID,
	).Scan(
		&ba.BookID,
		&ba.AuthorID,
	)

	if err != nil {
		logger.Error("unable to insert record", "error", err)
		return nil, err
	}

	logger.Info("returning inserted book", "insertedBook", ba)
	return ba, nil
}
