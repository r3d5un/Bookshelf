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
	Timestamp sql.NullTime   `json:"timestamp"`
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
		&task.Timestamp,
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

func (m *TaskModel) GetAll(ctx context.Context, name *string) (task []*Task, err error) {
	return nil, nil
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
