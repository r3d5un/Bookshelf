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
	orchestratorData "github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	orchestratorTypes "github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

type Monolith interface {
	Context() context.Context
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
	Books        Books
	UI           UI
	Orchestrator Orchestrator
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
	// Genre
	CreateGenre(ctx context.Context, newGenreData types.NewGenreData) (*uuid.UUID, error)
	ReadGenre(ctx context.Context, genreID uuid.UUID) (*types.Genre, error)
	ReadAllGenre(ctx context.Context, filters data.Filters) ([]*types.Genre, error)
	UpdateGenre(ctx context.Context, newGenreData types.Genre) (*types.Genre, error)
	DeleteGenre(ctx context.Context, id uuid.UUID) error
	// Books
	CreateBook(ctx context.Context, newBookData types.Book) (*uuid.UUID, error)
	ReadBook(ctx context.Context, genreID uuid.UUID) (*types.Book, error)
	ReadAllBook(ctx context.Context, filters data.Filters) ([]*types.Book, error)
	ReadBooksBySeries(ctx context.Context, seriesID uuid.UUID) ([]*types.Book, error)
	UpdateBook(ctx context.Context, newBookDAta types.Book) (*types.Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
}

type UI interface{}

type Orchestrator interface {
	ReadScheduledTask(
		ctx context.Context,
		taskID uuid.UUID,
	) (*orchestratorTypes.ScheduledTask, error)
	ReadAllScheduledTasks(
		ctx context.Context,
		filters orchestratorData.Filters,
	) (*orchestratorTypes.ScheduledTaskCollection, error)
	CreateScheduledTask(
		ctx context.Context,
		newTask orchestratorTypes.ScheduledTask,
	) (*orchestratorTypes.ScheduledTask, error)
	UpdateScheduledTask(
		ctx context.Context,
		newTaskData orchestratorTypes.ScheduledTask,
	) (*orchestratorTypes.ScheduledTask, error)
	DeleteScheduledTask(ctx context.Context, id uuid.UUID) (*orchestratorTypes.ScheduledTask, error)
	ClaimScheduledTaskByID(
		ctx context.Context,
		taskID uuid.UUID,
	) (*orchestratorTypes.ScheduledTask, error)
	SetScheduledTaskState(
		ctx context.Context,
		taskID uuid.UUID,
		state orchestratorData.TaskState,
	) (*orchestratorTypes.ScheduledTask, error)
	DequeueScheduledTask(
		ctx context.Context,
		taskID uuid.UUID,
	) (*orchestratorTypes.ScheduledTask, error)
}
