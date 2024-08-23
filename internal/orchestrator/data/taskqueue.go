package data

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
)

const (
	WaitingTaskState  string = "waiting"
	RunningTaskState  string = "running"
	CompleteTaskState string = "complete"
	StoppedTaskState  string = "stopped"
)

type TaskQueue struct {
	ID        uuid.UUID  `json:"id"`
	Queue     *string    `json:"queue"`
	State     *string    `json:"state"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	RunAt     *time.Time `json:"runAt"`
}

type TaskQueueModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}

func (m *TaskQueueModel) Get(ctx context.Context, id uuid.UUID) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       queue,
       state,
       created_at,
       updated_at,
       run_at
FROM orchestrator.tasks
WHERE id = '$1';
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			slog.String("id", id.String()),
		),
	)

	task = &TaskQueue{}

	logger.Info("performing query")
	err = m.Pool.QueryRow(qCtx, query, id.String()).Scan(
		&task.ID,
		&task.Queue,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found", slog.String("taskId", id.String()))
			return nil, ErrRecordNotFound
		default:
			logger.Info("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil
}

func (m *TaskQueueModel) GetAll(
	ctx context.Context,
	filters Filters,
) (tasks []*TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT COUNT(*) OVER() AS total,
       id,
       queue,
       state,
       created_at,
       updated_at,
       run_at
FROM orchestrator.tasks
WHERE ($1::uuid = '' OR id = $1::uuid)
  AND ($2::text = '' OR queue = $2::text)
  AND ($3::text = '' OR state = $3::text)
  AND ($4::timestamp IS NULL OR created_at >= $4::timestamp)
  AND ($5::timestamp IS NULL OR created_at < $5::timestamp)
  AND ($6::timestamp IS NULL OR updated_at >= $6::timestamp)
  AND ($7::timestamp IS NULL OR updated_at < $7::timestamp)
  AND ($6::timestamp IS NULL OR run_at >= $6::timestamp)
  AND ($7::timestamp IS NULL OR run_at < $7::timestamp)
` + database.CreateOrderByClause(filters.OrderBy) + `
OFFSET $8 FETCH NEXT $9 ROWS ONLY;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"filters", filters,
		),
	)

	tasks = []*TaskQueue{}
	totalResults := 0

	logger.Info("performing query")
	rows, err := m.Pool.Query(
		qCtx,
		query,
		filters.ID,
		filters.Queue,
		filters.State,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
		filters.RunAtFrom,
		filters.RunAtTo,
	)
	defer rows.Close()

	for rows.Next() {
		var task TaskQueue

		err := rows.Scan(
			&totalResults,
			&task.ID,
			&task.Queue,
			&task.State,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.RunAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		logger.Error("an error occurred while parsing query results", "error", err)
		return nil, err
	}

	logger.Info("calculating metadata")
	metadata := calculateMetadata(totalResults, filters.Page, filters.PageSize, filters.OrderBy)
	logger.Info("metadata calculated", "metadata", metadata)

	logger.Info("returning records")
	return tasks, nil
}

func (m *TaskQueueModel) Insert(
	ctx context.Context,
	taskQueue string,
	state *string,
	runAt *time.Time,
) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO orchestrator.task (queue,
                               state,
                               run_at)
VALUES ($1::TEXT,
        COALESCE($2::TEXT, 'waiting'),
        COALESCE($3::TIMESTAMP, CURRENT_TIMESTAMP))
RETURNING
    id,
    queue,
    state,
    created_at,
    updated_at,
    run_at;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			slog.String("taskQueue", taskQueue),
		),
	)

	task = &TaskQueue{}

	logger.Info("performing query")
	err = m.Pool.QueryRow(qCtx, query, taskQueue, state, runAt).Scan(
		&task.ID,
		&task.Queue,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.Info("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil
}

func (m *TaskQueueModel) Update(
	ctx context.Context,
	taskQueue TaskQueue,
) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
UPDATE orchestrator.task
SET queue = COALESCE($2, queue),
	state = COALESCE($2, state),
	created_at = COALESCE($3, created_at),
	run_at = COALESCE($4, run_at)
WHERE id = $1
RETURNING
    id,
    queue,
    state,
    created_at,
    updated_at,
    run_at;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"task", taskQueue,
		),
	)

	task = &TaskQueue{}

	logger.Info("performing query")
	err = m.Pool.QueryRow(qCtx, query, task.ID).Scan(
		&task.ID,
		&task.Queue,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.Info("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil
}

func (m *TaskQueueModel) Delete(ctx context.Context, id uuid.UUID) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
DELETE FROM orchestrator.task
WHERE id = $1
RETURNING
    id,
    queue,
    state,
    created_at,
    updated_at,
    run_at;

`
	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			slog.String("id", id.String()),
		),
	)

	task = &TaskQueue{}

	logger.Info("performing query")
	err = m.Pool.QueryRow(qCtx, query, task.ID).Scan(
		&task.ID,
		&task.Queue,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found")
			return nil, ErrRecordNotFound
		default:
			logger.Info("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil
}

func (m *TaskQueueModel) Dequeue(ctx context.Context, id uuid.UUID) (*TaskQueue, error) {
	// TODO: Implement
	return nil, nil
}

func (m *TaskQueueModel) Notify(ctx context.Context, queue string) error {
	logger := logging.LoggerFromContext(ctx)

	query := "NOTIFY task_queue_notification, $1;"

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
		),
	)

	conn, err := m.Pool.Acquire(context.Background())
	if err != nil {
		logger.Error("unable to acquire connection", "error", err)
		return err
	}
	defer conn.Release()

	logger.Info("executing notification statement")
	_, err = conn.Exec(context.Background(), query, queue)
	if err != nil {
		logger.Error("unable to execute notification statement", "error", err)
		return err
	}

	return nil
}

func (m *TaskQueueModel) Listen(
	ctx context.Context,
	notificationCh chan<- pgconn.Notification,
	done <-chan bool,
) {
	logger := logging.LoggerFromContext(ctx)

	query := `LISTEN task_queue_notification;`

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
		),
	)

	for {
		conn, err := m.Pool.Acquire(ctx)
		if err != nil {
			logger.Error("unable to acquire connection", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Info("listening for task notifications")
		_, err = conn.Exec(context.Background(), query)
		if err != nil {
			logger.Error("unable to listen", "error", err)
			conn.Release()
			continue
		}

		for {
			select {
			case <-done:
				logger.Info("received done signal, stopping listener.")
				return
			default:
				if notification, err := conn.Conn().WaitForNotification(ctx); err != nil {
					logger.Info("unable to receive notification", "error", err)
					break
				} else {
					select {
					case notificationCh <- *notification:
						logger.Info("notification sent to channel")
					case <-done:
						logger.Info("received done signal, stopping listener.")
						return
					default:
						logger.Error("channel full, notification not sent")
						time.Sleep(5 * time.Second)
						continue
					}
				}
			}
		}

	}
}
