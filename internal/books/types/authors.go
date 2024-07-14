package types

import (
	"context"
	"sync"
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

func ReadAuthor(ctx context.Context, models *data.Models, authorID uuid.UUID) (*Author, error) {
	authorRecord, err := models.Authors.Get(ctx, authorID)
	if err != nil {
		return nil, err
	}

	authorData := Author{
		ID:          authorRecord.ID,
		Name:        authorRecord.Name,
		Description: authorRecord.Description,
		Website:     authorRecord.Website,
		CreatedAt:   authorRecord.CreatedAt,
		UpdatedAt:   authorRecord.UpdatedAt,
	}

	bookRecords, totalBookRecords, err := models.Books.GetByAuthorID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	if *totalBookRecords < 1 {
		return &authorData, nil
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
			authorData.Books = append(authorData.Books, bookData)
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

	return &authorData, nil
}

func ListAuthor(ctx context.Context, models *data.Models, filters data.Filters) ([]*Author, error) {

}

func UpdateAuthor(ctx context.Context, models *data.Models, newAuthorData Author) (*Author, error) {
	authorRecord := data.Author{
		ID:          newAuthorData.ID,
		Name:        newAuthorData.Name,
		Description: newAuthorData.Description,
		Website:     newAuthorData.Website,
	}

	updatedAuthor, err := models.Authors.Update(ctx, authorRecord)
	if err != nil {
		return nil, err
	}

	updatedAuthorData, err := ReadAuthor(ctx, models, updatedAuthor.ID)
	if err != nil {
		return nil, err
	}

	return updatedAuthorData, nil
}

func DeleteAuthor(ctx context.Context, models *data.Models, id uuid.UUID) error {
	_, err := models.Books.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
