package types

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

type NewSeriesData struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type Series struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Books       []*Book    `json:"books,omitempty"`
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

func ReadSeries(ctx context.Context, models *data.Models, seriesID uuid.UUID) (*Series, error) {
	seriesRecord, err := models.Series.Get(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	seriesData := Series{
		ID:          seriesRecord.ID,
		Name:        &seriesRecord.Name,
		Description: seriesRecord.Description,
		CreatedAt:   seriesRecord.CreatedAt,
		UpdatedAt:   seriesRecord.UpdatedAt,
	}

	bookRecords, totalBookRecords, err := models.Books.GetBySeriesID(ctx, seriesID)
	if err != nil {
		return nil, err
	}
	if *totalBookRecords < 1 {
		return &seriesData, nil
	}

	var wg sync.WaitGroup
	var authorDataMu sync.Mutex

	errorChan := make(chan error, *totalBookRecords)

	for _, bookRecord := range bookRecords {
		wg.Add(1)
		go func(ctx context.Context, models *data.Models, id uuid.UUID) {
			defer wg.Done()

			bookData, err := ReadBook(ctx, models, id)
			if err != nil {
				errorChan <- err
			}

			authorDataMu.Lock()
			seriesData.Books = append(seriesData.Books, bookData)
			authorDataMu.Unlock()
		}(ctx, models, bookRecord.ID)
	}

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}

	return &seriesData, nil
}
