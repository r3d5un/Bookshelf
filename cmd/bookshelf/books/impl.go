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
