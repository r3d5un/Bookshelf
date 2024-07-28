package system

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
	"github.com/r3d5un/Bookshelf/internal/config"
)

type Monolith interface {
	Logger() *slog.Logger
	Mux() *http.ServeMux
	DB() *sql.DB
	Config() *config.Config
	Modules() *Modules
}

type Module interface {
	Startup(context.Context, Monolith) error
	Shutdown()
}

type Modules struct {
	Books Books
	UI    UI
}

type Books interface {
	// Authors
	CreateAuthor(ctx context.Context, data types.NewAuthorData) (*uuid.UUID, error)
	ReadAuthor(ctx context.Context, id uuid.UUID) (*types.Author, error)
	ReadAllAuthors(ctx context.Context, filters data.Filters) ([]*types.Author, error)
	UpdateAuthor(ctx context.Context, data types.Author) (*types.Author, error)
	DeleteAuthor(ctx context.Context, id uuid.UUID) error
	// Series
	CreateSeries(ctx context.Context, newSeriesData types.NewSeriesData) (*uuid.UUID, error)
	ReadSeries(ctx context.Context, seriesID uuid.UUID) (*types.Series, error)
	ReadAllSeries(ctx context.Context, filters data.Filters) ([]*types.Series, error)
	UpdateSeries(ctx context.Context, newSeriesData types.Series) (*types.Series, error)
	DeleteSeries(ctx context.Context, id uuid.UUID) error
}

type UI interface{}
