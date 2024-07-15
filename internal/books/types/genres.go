package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

type NewGenreData struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type Genre struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Books       []*Book    `json:"books,omitempty"`
}

func CreateGenre(
	ctx context.Context,
	models *data.Models,
	newGenreData NewGenreData,
) (*uuid.UUID, error) {
	genreRecord := data.Genre{
		ID:          uuid.New(),
		Name:        newGenreData.Name,
		Description: newGenreData.Description,
	}

	insertedGenre, err := models.Genres.Insert(ctx, genreRecord)
	if err != nil {
		return nil, err
	}

	return &insertedGenre.ID, nil
}
