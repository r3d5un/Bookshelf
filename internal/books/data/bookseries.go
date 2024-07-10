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

type BookSeries struct {
	BookID      uuid.UUID `json:"bookId"`
	SeriesID    uuid.UUID `json:"seriesId"`
	SeriesOrder float32   `json:"seriesOrder"`
}

type BookSeriesModel struct {
	DB      *sql.DB
	Timeout *time.Duration
}

func (m BookSeriesModel) Insert(
	ctx context.Context,
	bookID uuid.UUID,
	seriesID uuid.UUID,
	order float32,
) (ba *BookSeries, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO books.book_series (book_id,
                               series_id,
                               series_order)
VALUES ($1,
        $2,
        $3)
RETURNING book_id,
          series_id,
          series_order;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			slog.String("bookId", bookID.String()),
			slog.String("seriesId", seriesID.String()),
		),
	)

	ba = &BookSeries{}

	logger.Info("performing query")
	err = m.DB.QueryRowContext(
		qCtx,
		query,
		bookID,
		seriesID,
		order,
	).Scan(
		&ba.BookID,
		&ba.SeriesID,
		&ba.SeriesOrder,
	)

	if err != nil {
		logger.Error("unable to insert record", "error", err)
		return nil, err
	}

	logger.Info("returning inserted book", "insertedBook", ba)
	return ba, nil
}
