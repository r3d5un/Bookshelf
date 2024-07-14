package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

type NewAuthorData struct {
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Website     *string `json:"website,omitempty"`
}

type Author struct {
	ID          uuid.UUID  `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Website     *string    `json:"website"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	Books       []*Book    `json:"books"`
}

func CreateAuthor(
	ctx context.Context,
	models *data.Models,
	newAuthorData NewAuthorData,
) (*uuid.UUID, error) {
	authorRecord := data.Author{
		ID:          uuid.New(),
		Name:        &newAuthorData.Name,
		Description: newAuthorData.Description,
		Website:     newAuthorData.Website,
	}

	insertedAuthor, err := models.Authors.Insert(ctx, authorRecord)
	if err != nil {
		return nil, err
	}

	return &insertedAuthor.ID, nil
}

