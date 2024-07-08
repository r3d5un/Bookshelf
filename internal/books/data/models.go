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
	Books   BookModel
	Authors AuthorModel
	Series  SeriesModel
}

func NewModels(db *sql.DB, timeout *time.Duration) Models {
	return Models{
		Books:   BookModel{DB: db, Timeout: timeout},
		Authors: AuthorModel{DB: db, Timeout: timeout},
		Series:  SeriesModel{DB: db, Timeout: timeout},
	}
}
