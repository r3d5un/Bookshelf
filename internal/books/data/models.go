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
	Authors AuthorModel
	Books   BookModel
	Genres  GenreModel
	Series  SeriesModel
}

func NewModels(db *sql.DB, timeout *time.Duration) Models {
	return Models{
		Authors: AuthorModel{DB: db, Timeout: timeout},
		Books:   BookModel{DB: db, Timeout: timeout},
		Genres:  GenreModel{DB: db, Timeout: timeout},
		Series:  SeriesModel{DB: db, Timeout: timeout},
	}
}
