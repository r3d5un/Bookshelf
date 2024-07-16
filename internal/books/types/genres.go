package types

import (
	"context"
	"sync"
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
		Name:        &newGenreData.Name,
		Description: newGenreData.Description,
	}

	insertedGenre, err := models.Genres.Insert(ctx, genreRecord)
	if err != nil {
		return nil, err
	}

	return &insertedGenre.ID, nil
}

func ReadGenre(ctx context.Context, models *data.Models, genreID uuid.UUID) (*Genre, error) {
	genreRecord, err := models.Genres.Get(ctx, genreID)
	if err != nil {
		return nil, err
	}

	genreData := Genre{
		ID:          genreRecord.ID,
		Name:        genreRecord.Name,
		Description: genreRecord.Description,
		CreatedAt:   genreRecord.CreatedAt,
		UpdatedAt:   genreRecord.UpdatedAt,
	}

	bookRecords, totalBookRecords, err := models.Books.GetByGenreID(ctx, genreID)
	if err != nil {
		return nil, err
	}
	if *totalBookRecords < 1 {
		return &genreData, nil
	}

	var wg sync.WaitGroup
	var genreDataMu sync.Mutex

	errorChan := make(chan error, *totalBookRecords)

	for _, bookRecord := range bookRecords {
		wg.Add(1)
		go func(ctx context.Context, models *data.Models, id uuid.UUID) {
			defer wg.Done()

			bookData, err := ReadBook(ctx, models, id)
			if err != nil {
				errorChan <- err
			}

			genreDataMu.Lock()
			genreData.Books = append(genreData.Books, bookData)
			genreDataMu.Unlock()
		}(ctx, models, bookRecord.ID)
	}

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}

	return &genreData, nil
}
