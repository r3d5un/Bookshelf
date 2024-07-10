package types

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

type Book struct {
	data.Book
	Authors    []*data.Author     `json:"author,omitempty"`
	Genres     []*data.Genre      `json:"genre,omitempty"`
	Series     []*data.Series     `json:"series,omitempty"`
	BookSeries []*data.BookSeries `json:"bookSeries,omitempty"`
}

func GetBook(
	ctx context.Context,
	models *data.Models,
	bookID uuid.UUID,
) (book *Book, err error) {
	bookCh := make(chan bookDataResult, 1)
	defer close(bookCh)
	authorCh := make(chan authorDataResult, 1)
	defer close(authorCh)
	seriesCh := make(chan seriesDataResult, 1)
	defer close(seriesCh)

	go func() {
		go getBookData(ctx, models, bookID, bookCh)
		close(bookCh)
	}()
	go func() {
		go getBookAuthorData(ctx, models, bookID, authorCh)
		close(authorCh)
	}()
	go func() {
		go getBookSeriesData(ctx, models, bookID, seriesCh)
		close(seriesCh)
	}()

	bookData := <-bookCh
	authorData := <-authorCh
	seriesData := <-seriesCh

	if bookData.err != nil {
		return nil, bookData.err
	}
	if authorData.err != nil {
		return nil, authorData.err
	}
	if seriesData.err != nil {
		return nil, seriesData.err
	}

	book = &Book{
		Book:    *bookData.book,
		Authors: authorData.authors,
		Series:  seriesData.series,
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

func NewBook(ctx context.Context, models *data.Models, newBook Book) error {
	insertedBook, err := models.Books.Insert(ctx, newBook.Book)
	if err != nil {
		return err
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
			if _, err := models.BookSeries.Insert(ctx, series.BookID, series.SeriesID, series.SeriesOrder); err != nil {
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
			return err
		}
	}

	return nil
}
