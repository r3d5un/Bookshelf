package books

import (
	"context"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func (m *Module) CreateAuthor(ctx context.Context, data types.NewAuthorData) (*uuid.UUID, error) {
	id, err := types.CreateAuthor(ctx, &m.models, data)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (m *Module) ReadAuthor(ctx context.Context, id uuid.UUID) (*types.Author, error) {
	a, err := types.ReadAuthor(ctx, &m.models, id)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) ReadAllAuthors(
	ctx context.Context,
	filters data.Filters,
) ([]*types.Author, error) {
	a, err := types.ReadAllAuthors(ctx, &m.models, filters)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) UpdateAuthor(ctx context.Context, data types.Author) (*types.Author, error) {
	a, err := types.UpdateAuthor(ctx, &m.models, data)
	if err != nil {
		return a, nil
	}

	return a, nil
}

func (m *Module) DeleteAuthor(ctx context.Context, id uuid.UUID) error {
	err := types.DeleteAuthor(ctx, &m.models, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) CreateSeries(ctx context.Context, data types.NewSeriesData) (*uuid.UUID, error) {
	id, err := types.CreateSeries(ctx, &m.models, data)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (m *Module) ReadSeries(ctx context.Context, id uuid.UUID) (*types.Series, error) {
	a, err := types.ReadSeries(ctx, &m.models, id)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) ReadAllSeries(
	ctx context.Context,
	filters data.Filters,
) ([]*types.Series, error) {
	a, err := types.ReadAllSeries(ctx, &m.models, filters)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) UpdateSeries(ctx context.Context, data types.Series) (*types.Series, error) {
	a, err := types.UpdateSeries(ctx, &m.models, data)
	if err != nil {
		return a, nil
	}

	return a, nil
}

func (m *Module) DeleteSeries(ctx context.Context, id uuid.UUID) error {
	err := types.DeleteSeries(ctx, &m.models, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) CreateGenre(ctx context.Context, data types.NewGenreData) (*uuid.UUID, error) {
	id, err := types.CreateGenre(ctx, &m.models, data)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (m *Module) ReadGenre(ctx context.Context, id uuid.UUID) (*types.Genre, error) {
	a, err := types.ReadGenre(ctx, &m.models, id)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) ReadAllGenre(
	ctx context.Context,
	filters data.Filters,
) ([]*types.Genre, error) {
	a, err := types.ReadAllGenre(ctx, &m.models, filters)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) UpdateGenre(ctx context.Context, data types.Genre) (*types.Genre, error) {
	a, err := types.UpdateGenre(ctx, &m.models, data)
	if err != nil {
		return a, nil
	}

	return a, nil
}

func (m *Module) DeleteGenre(ctx context.Context, id uuid.UUID) error {
	err := types.DeleteGenre(ctx, &m.models, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) CreateBook(ctx context.Context, data types.Book) (*uuid.UUID, error) {
	id, err := types.CreateBook(ctx, &m.models, data)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (m *Module) ReadBook(ctx context.Context, id uuid.UUID) (*types.Book, error) {
	a, err := types.ReadBook(ctx, &m.models, id)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) ReadAllBook(
	ctx context.Context,
	filters data.Filters,
) ([]*types.Book, error) {
	a, err := types.ReadAllBooks(ctx, &m.models, filters)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (m *Module) ReadBooksBySeries(ctx context.Context, seriesID uuid.UUID) ([]*types.Book, error) {
	books, err := types.ReadBooksBySeries(ctx, &m.models, seriesID)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (m *Module) UpdateBook(ctx context.Context, data types.Book) (*types.Book, error) {
	a, err := types.UpdateBook(ctx, &m.models, data)
	if err != nil {
		return a, nil
	}

	return a, nil
}

func (m *Module) DeleteBook(ctx context.Context, id uuid.UUID) error {
	err := types.DeleteBook(ctx, &m.models, id)
	if err != nil {
		return err
	}

	return nil
}
