package data

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
)

type Task struct {
	Name      sql.NullString `json:"name"`
	CronExpr  sql.NullString `json:"cronExpr"`
	Enabled   sql.NullBool   `json:"enabled"`
	Deleted   sql.NullBool   `json:"deleted"`
	UpdatedAt sql.NullTime   `json:"timestamp"`
}

type TaskModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}

func (m *TaskModel) Get(ctx context.Context, name string) (task *Task, err error) {
	query := `
SELECT
    name,
    cron_expr,
    enabled,
    deleted,
    updated_at
FROM orchestrator.tasks
WHERE name = $1;
`

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", database.MinifySQL(query)),
		"name", slog.String("name", name),
	))

	task = &Task{}

	logger.Info("performing query")
	err = m.Pool.QueryRow(ctx, query, name).Scan(
		&task.Name,
		&task.CronExpr,
		&task.Enabled,
		&task.Deleted,
		&task.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.Error("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil
}

func (m *TaskModel) GetAll(
	ctx context.Context,
	filters Filters,
) (tasks []*Task, metadata *Metadata, err error) {
	query := `
SELECT COUNT(*) OVER() AS total,
       name,
       cron_expr,
       enabled,
       deleted,
       updated_at
FROM orchestrator.tasks
WHERE ($1::text IS NULL OR name = $1::text)
  AND ($2::text IS NULL OR cron_expr = $2::text)
  AND ($3::boolean IS NULL OR enabled = $3::boolean)
  AND ($4::boolean IS NULL OR deleted >= $4::boolean)
  AND ($5::timestamp IS NULL OR updated_at >= $5::timestamp)
  AND ($6::timestamp IS NULL OR updated_at < $6::timestamp)
` + database.CreateOrderByClause(filters.OrderBy) + `
OFFSET $7 FETCH NEXT $8 ROWS ONLY;
`

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", database.MinifySQL(query)),
		"name", slog.Any("filters", filters),
	))

	tasks = []*Task{}
	totalResults := 0

	logger.Info("performing query")
	rows, err := m.Pool.Query(
		ctx,
		query,
		filters.Name,
		filters.CronExpr,
		filters.Enabled,
		filters.Deleted,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
		filters.offset(),
		filters.limit(),
	)
	defer rows.Close()

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&totalResults,
			&task.Name,
			&task.CronExpr,
			&task.Enabled,
			&task.Deleted,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, nil, err
	}

	logger.Info("calculating metadata")
	md := calculateMetadata(totalResults, filters.Page, filters.PageSize, filters.OrderBy)
	logger.Info("metadata calculated", "metadata", metadata)

	logger.Info("returning records")
	return tasks, &md, nil
}

func (m *TaskModel) Insert(ctx context.Context, newTask Task) (task []*Task, err error) {
	return nil, nil
}

func (m *TaskModel) Update(ctx context.Context, newData Task) (task *Task, err error) {
	return nil, nil
}

func (m *TaskModel) Upsert(ctx context.Context, newTask Task) (task *Task, err error) {
	return nil, nil
}

func (m *TaskModel) Delete(ctx context.Context, name string) (task *Task, err error) {
	return nil, nil
}
