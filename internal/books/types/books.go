package types

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

type Book struct {
	ID          *uuid.UUID         `json:"id,omitempty"`
	Title       *string            `json:"title,omitempty"`
	Description *string            `json:"description,omitempty"`
	Published   *time.Time         `json:"published,omitempty"`
	CreatedAt   *time.Time         `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time         `json:"updatedAt,omitempty"`
	Authors     []*data.Author     `json:"author,omitempty"`
	Genres      []*data.Genre      `json:"genre,omitempty"`
	Series      []*data.Series     `json:"series,omitempty"`
	BookSeries  []*data.BookSeries `json:"bookSeries,omitempty"`
}

// Retrieves and builds a Book object containing the complete dataset for a single book.
//
// If the book does not exist, nil and an ErrRecordNotFound error will be returned.
func ReadBook(ctx context.Context, models *data.Models, bookID uuid.UUID) (*Book, error) {
	bookCh := make(chan bookDataResult, 1)
	authorCh := make(chan authorDataResult, 1)
	seriesCh := make(chan seriesDataResult, 1)
	genreCh := make(chan genreDataResult, 1)
	errCh := make(chan error, 4)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		getBookData(ctx, models, bookID, bookCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		getBookAuthorData(ctx, models, bookID, authorCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		getBookSeriesData(ctx, models, bookID, seriesCh)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		getBookGenreData(ctx, models, bookID, genreCh)
	}()

	go func() {
		wg.Wait()
		close(bookCh)
		close(authorCh)
		close(seriesCh)
		close(genreCh)
		close(errCh)
	}()

	var bookData bookDataResult
	var authorData authorDataResult
	var seriesData seriesDataResult
	var genreData genreDataResult

	// Collect results from channels
	for i := 0; i < 4; i++ {
		select {
		case bd := <-bookCh:
			bookData = bd
			if bd.err != nil {
				errCh <- bd.err
			}
		case ad := <-authorCh:
			authorData = ad
			if ad.err != nil {
				errCh <- ad.err
			}
		case sd := <-seriesCh:
			seriesData = sd
			if sd.err != nil {
				errCh <- sd.err
			}
		case gd := <-genreCh:
			genreData = gd
			if gd.err != nil {
				errCh <- gd.err
			}
		}
	}

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	book := &Book{
		ID:          &bookData.book.ID,
		Title:       &bookData.book.Title,
		Description: bookData.book.Description,
		Published:   bookData.book.Published,
		CreatedAt:   bookData.book.CreatedAt,
		UpdatedAt:   bookData.book.UpdatedAt,
		Authors:     authorData.authors,
		Series:      seriesData.series,
		Genres:      genreData.genres,
	}

	return book, nil
}

type bookDataResult struct {
	book *data.Book
	err  error
}

func getBookData(
	ctx context.Context,
	models *data.Models,
	bookID uuid.UUID,
	bookCh chan<- bookDataResult,
) {
	data, err := models.Books.Get(ctx, bookID)
	bookCh <- bookDataResult{book: data, err: err}
}

type authorDataResult struct {
	authors []*data.Author
	err     error
}

func getBookAuthorData(
	ctx context.Context,
	models *data.Models,
	bookID uuid.UUID,
	authorCh chan<- authorDataResult,
) {
	data, _, err := models.Authors.GetByBookID(ctx, bookID)
	authorCh <- authorDataResult{authors: data, err: err}
}

type seriesDataResult struct {
	series []*data.Series
	err    error
}

func getBookSeriesData(
	ctx context.Context,
	models *data.Models,
	bookID uuid.UUID,
	seriesCh chan<- seriesDataResult,
) {
	data, _, err := models.Series.GetByBookID(ctx, bookID)
	seriesCh <- seriesDataResult{series: data, err: err}
}

type genreDataResult struct {
	genres []*data.Genre
	err    error
}

func getBookGenreData(
	ctx context.Context,
	models *data.Models,
	bookID uuid.UUID,
	genreCh chan<- genreDataResult,
) {
	data, _, err := models.Genres.GetByBookID(ctx, bookID)
	genreCh <- genreDataResult{genres: data, err: err}
}

func CreateBook(ctx context.Context, models *data.Models, newBook Book) (*uuid.UUID, error) {
	insertedBook, err := models.Books.Insert(ctx, data.Book{
		ID:          uuid.New(),
		Title:       *newBook.Title,
		Description: newBook.Description,
		Published:   newBook.Published,
		CreatedAt:   newBook.CreatedAt,
		UpdatedAt:   newBook.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errChanLength := len(newBook.Genres) + len(newBook.Series) + len(newBook.Authors)
	errCh := make(chan error, errChanLength)

	for _, genre := range newBook.Genres {
		wg.Add(1)
		go func(genre *data.Genre) {
			defer wg.Done()
			if _, err := models.BookGenres.Insert(ctx, insertedBook.ID, genre.ID); err != nil {
				errCh <- err
			}
		}(genre)
	}

	for _, series := range newBook.BookSeries {
		wg.Add(1)
		go func(series *data.BookSeries) {
			defer wg.Done()
			if _, err := models.BookSeries.Insert(ctx, insertedBook.ID, series.SeriesID, series.SeriesOrder); err != nil {
				errCh <- err
			}
		}(series)
	}

	for _, author := range newBook.Authors {
		wg.Add(1)
		go func(author *data.Author) {
			defer wg.Done()
			if _, err := models.BookAuthors.Insert(ctx, insertedBook.ID, author.ID); err != nil {
				errCh <- err
			}
		}(author)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return &insertedBook.ID, nil
}

func UpdateBook(ctx context.Context, models *data.Models, newBookData Book) (*Book, error) {
	bookRecord := data.Book{
		ID:          *newBookData.ID,
		Title:       *newBookData.Title,
		Description: newBookData.Description,
		Published:   newBookData.Published,
		CreatedAt:   newBookData.CreatedAt,
		UpdatedAt:   newBookData.UpdatedAt,
	}

	updatedBook, err := models.Books.Update(ctx, bookRecord)
	if err != nil {
		return nil, err
	}

	updatedBookData, err := ReadBook(ctx, models, updatedBook.ID)
	if err != nil {
		return nil, err
	}

	return updatedBookData, nil
}

func DeleteBook(ctx context.Context, models *data.Models, id uuid.UUID) error {
	_, err := models.Books.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func ReadAllBooks(ctx context.Context, models *data.Models, filters data.Filters) ([]*Book, error) {
	bookListData, totalResults, err := models.Books.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}
	if *totalResults < 1 {
		return []*Book{}, nil
	}

	var wg sync.WaitGroup
	var booksMu sync.Mutex

	var books []*Book
	errorChan := make(chan error, *totalResults)

	for _, bookData := range bookListData {
		wg.Add(1)
		go func(ctx context.Context, models *data.Models, id uuid.UUID) {
			defer wg.Done()

			b, err := ReadBook(ctx, models, id)
			if err != nil {
				errorChan <- err
			}

			booksMu.Lock()
			books = append(books, b)
			booksMu.Unlock()
		}(ctx, models, bookData.ID)
	}

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}

	return books, nil
}

func AddAuthorsToBook(
	ctx context.Context,
	models *data.Models,
	bookID uuid.UUID,
	authorIDs []*uuid.UUID,
) error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(authorIDs))

	for _, authorID := range authorIDs {
		wg.Add(1)
		go func(ctx context.Context, models *data.Models, id *uuid.UUID) {
			defer wg.Done()

			_, err := models.BookAuthors.Insert(ctx, bookID, *authorID)
			if err != nil {
				errorChan <- err
			}
		}(ctx, models, authorID)
	}

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadBooksbySeries(
	ctx context.Context,
	models *data.Models,
	seriesID uuid.UUID,
) ([]*Book, error) {
	booksInSeries, totalResults, err := models.Books.GetBySeriesID(ctx, seriesID)
	if err != nil {
		return nil, err
	}
	if *totalResults < 1 {
		return []*Book{}, nil
	}

	var wg sync.WaitGroup
	var booksMu sync.Mutex

	var books []*Book
	errorChan := make(chan error, *totalResults)

	for _, bookData := range booksInSeries {
		wg.Add(1)

		go func(ctx context.Context, models *data.Models, id uuid.UUID) {
			defer wg.Done()

			b, err := ReadBook(ctx, models, id)
			if err != nil {
				errorChan <- err
			}

			booksMu.Lock()
			books = append(books, b)
			booksMu.Unlock()
		}(ctx, models, bookData.ID)
	}

	wg.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}

	return books, nil
}
