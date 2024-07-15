package types

import (
	"context"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

type NewSeriesData struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func CreateSeries(
	ctx context.Context,
	models *data.Models,
	newSeriesData NewSeriesData,
) (*uuid.UUID, error) {
	seriesRecord := data.Series{
		ID:          uuid.New(),
		Name:        newSeriesData.Name,
		Description: newSeriesData.Description,
	}

	insertedSeries, err := models.Series.Insert(ctx, seriesRecord)
	if err != nil {
		return nil, err
	}

	return &insertedSeries.ID, nil
}
