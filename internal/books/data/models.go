package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

var (
	DplicateKeyMatchString = "cannot insert duplicate key object"
)

type Models struct {
	Authors     AuthorModel
	Books       BookModel
	BookAuthors BookAuthorModel
	BookGenres  BookGenreModel
	BookSeries  BookSeriesModel
	Genres      GenreModel
	Series      SeriesModel
}

func NewModels(db *sql.DB, timeout *time.Duration) Models {
	return Models{
		Authors:     AuthorModel{DB: db, Timeout: timeout},
		Books:       BookModel{DB: db, Timeout: timeout},
		BookAuthors: BookAuthorModel{DB: db, Timeout: timeout},
		BookGenres:  BookGenreModel{DB: db, Timeout: timeout},
		BookSeries:  BookSeriesModel{DB: db, Timeout: timeout},
		Genres:      GenreModel{DB: db, Timeout: timeout},
		Series:      SeriesModel{DB: db, Timeout: timeout},
	}
}
