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
		Name:        &newSeriesData.Name,
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
		Name:        seriesRecord.Name,
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
	var seriesDataMu sync.Mutex

	errorChan := make(chan error, *totalBookRecords)

	for _, bookRecord := range bookRecords {
		wg.Add(1)
		go func(ctx context.Context, models *data.Models, id uuid.UUID) {
			defer wg.Done()

			bookData, err := ReadBook(ctx, models, id)
			if err != nil {
				errorChan <- err
			}

			seriesDataMu.Lock()
			seriesData.Books = append(seriesData.Books, bookData)
			seriesDataMu.Unlock()
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

func ReadAllSeries(
	ctx context.Context,
	models *data.Models,
	filters data.Filters,
) ([]*Series, error) {
	seriesListData, totalResults, err := models.Series.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	if *totalResults < 1 {
		return []*Series{}, nil
	}

	var wg sync.WaitGroup
	var seriesMu sync.Mutex

	var series []*Series
	errorChan := make(chan error, *totalResults)

	for _, seriesData := range seriesListData {
		wg.Add(1)
		go func(ctx context.Context, models *data.Models, id uuid.UUID) {
			defer wg.Done()

			s, err := ReadSeries(ctx, models, id)
			if err != nil {
				errorChan <- err
			}

			seriesMu.Lock()
			series = append(series, s)
			seriesMu.Unlock()
		}(ctx, models, seriesData.ID)
	}

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}

	return series, nil
}

func UpdateSeries(ctx context.Context, models *data.Models, newSeriesData Series) (*Series, error) {
	seriesRecord := data.Series{
		ID:          newSeriesData.ID,
		Name:        newSeriesData.Name,
		Description: newSeriesData.Description,
		CreatedAt:   newSeriesData.CreatedAt,
		UpdatedAt:   newSeriesData.UpdatedAt,
	}

	updatedSeries, err := models.Series.Update(ctx, seriesRecord)
	if err != nil {
		return nil, err
	}

	updatedSeriesData, err := ReadSeries(ctx, models, updatedSeries.ID)
	if err != nil {
		return nil, err
	}

	return updatedSeriesData, nil
}

func DeleteSeries(ctx context.Context, models *data.Models, id uuid.UUID) error {
	_, err := models.Series.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
