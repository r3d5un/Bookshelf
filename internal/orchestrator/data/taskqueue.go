package data

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
)

const (
	WaitingTaskState  string = "waiting"
	RunningTaskState  string = "running"
	CompleteTaskState string = "complete"
	StoppedTaskState  string = "stopped"
	ErrorTaskState    string = "error"
)

type TaskQueue struct {
	ID        uuid.UUID  `json:"id"`
	Name      *string    `json:"name"`
	State     *string    `json:"state"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	RunAt     *time.Time `json:"runAt"`
	TaskData  *string    `json:"task_data,omitempty"`
}

type TaskQueueModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}

func (m *TaskQueueModel) Get(ctx context.Context, id uuid.UUID) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       state,
       created_at,
       updated_at,
       run_at,
       task_data
FROM orchestrator.tasks
WHERE id = $1;
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
		&task.Name,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
		&task.TaskData,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found", slog.String("taskId", id.String()))
			return nil, ErrRecordNotFound
		default:
			logger.Error("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil
}

func (m *TaskQueueModel) GetAll(
	ctx context.Context,
	filters Filters,
) (tasks []*TaskQueue, metadata *Metadata, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT COUNT(*) OVER() AS total,
       id,
       name,
       state,
       created_at,
       updated_at,
       run_at,
       task_data
FROM orchestrator.tasks
WHERE ($1::uuid IS NULL OR id = $1::uuid)
  AND ($2::text IS NULL OR name = $2::text)
  AND ($3::task_state IS NULL OR state = $3::task_state)
  AND ($4::timestamp IS NULL OR created_at >= $4::timestamp)
  AND ($5::timestamp IS NULL OR created_at < $5::timestamp)
  AND ($6::timestamp IS NULL OR updated_at >= $6::timestamp)
  AND ($7::timestamp IS NULL OR updated_at < $7::timestamp)
  AND ($8::timestamp IS NULL OR run_at >= $8::timestamp)
  AND ($9::timestamp IS NULL OR run_at < $9::timestamp)
` + database.CreateOrderByClause(filters.OrderBy) + `
OFFSET $10 FETCH NEXT $11 ROWS ONLY;
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
		filters.Name,
		filters.State,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
		filters.RunAtFrom,
		filters.RunAtTo,
		filters.offset(),
		filters.limit(),
	)
	defer rows.Close()

	for rows.Next() {
		var task TaskQueue

		err := rows.Scan(
			&totalResults,
			&task.ID,
			&task.Name,
			&task.State,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.RunAt,
			&task.TaskData,
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

func (m *TaskQueueModel) Insert(
	ctx context.Context,
	newTask TaskQueue,
) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
INSERT INTO orchestrator.tasks (name,
                               state,
                               run_at,
                               task_data)
VALUES ($1::TEXT,
        COALESCE($2::task_state, 'waiting'),
        COALESCE($3::TIMESTAMP, CURRENT_TIMESTAMP),
        COALESCE($4::JSONB, NULL))
RETURNING
    id,
    name,
    state,
    created_at,
    updated_at,
    run_at,
    task_data;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"newTask", newTask,
		),
	)

	task = &TaskQueue{}

	logger.Info("performing query")
	err = m.Pool.QueryRow(qCtx, query, newTask.Name, newTask.State, newTask.RunAt, newTask.RunAt).
		Scan(
			&task.ID,
			&task.Name,
			&task.State,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.RunAt,
			&task.TaskData,
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

func (m *TaskQueueModel) Update(
	ctx context.Context,
	newTaskData TaskQueue,
) (updatedTask *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
UPDATE orchestrator.tasks
SET name = COALESCE($2::text, name),
	state = COALESCE($3::task_state, state),
	created_at = COALESCE($4::timestamp, created_at),
	run_at = COALESCE($5::timestamp, run_at),
	task_data = COALESCE($6::JSONB, task_data)
WHERE id = $1::uuid
RETURNING
    id,
    name,
    state,
    created_at,
    updated_at,
    run_at,
    task_data;
`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
			"task", newTaskData,
		),
	)

	updatedTask = &TaskQueue{}

	logger.Info("performing query")
	err = m.Pool.QueryRow(
		qCtx,
		query,
		newTaskData.ID.String(),
		newTaskData.Name,
		newTaskData.State,
		newTaskData.CreatedAt,
		newTaskData.RunAt,
		newTaskData.TaskData,
	).Scan(
		&updatedTask.ID,
		&updatedTask.Name,
		&updatedTask.State,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
		&updatedTask.RunAt,
		&updatedTask.TaskData,
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
	return updatedTask, nil
}

func (m *TaskQueueModel) Delete(ctx context.Context, id uuid.UUID) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
DELETE FROM orchestrator.tasks
WHERE id = $1
RETURNING
    id,
    name,
    state,
    created_at,
    updated_at,
    run_at,
    task_data;

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
		&task.Name,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
		&task.TaskData,
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

func (m *TaskQueueModel) ConsumeByID(
	ctx context.Context,
	taskCh chan<- TaskQueue,
	taskRunResultCh <-chan error,
	id uuid.UUID,
) error {
	logger := logging.LoggerFromContext(ctx)

	logger.Info("creating transaction for consuming task", "id", id)
	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		slog.Error("unable to begin transaction", "error", err)
		return err
	}
	defer tx.Commit(ctx)

	var commitSuccessful bool
	defer func() {
		if !commitSuccessful {
			tx.Rollback(context.Background())
		}
	}()

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	slog.Info("reading and locking task", "id", id)
	task, err := m.ClaimTx(qCtx, tx, id)
	if err != nil {
		return err
	}

	// send task to consumer for processing
	logger.Info("sending task for processing")
	taskCh <- *task
	logger.Info("task sent")

	// receive reply from task processor
	logger.Info("waiting for reply")
	err = <-taskRunResultCh
	// An error should only occur when the task was not able to complete.
	// Any failed tasks should be marked with an error.
	if err != nil {
		logger.Error("task unsuccssful, setting task state to error", slog.Any("error", err))
		errorState := ErrorTaskState
		task.State = &errorState
		_, err = m.UpdateTx(ctx, tx, *task)
		if err != nil {
			logger.Error(
				"an error occurred while updating the task state",
				"task", task,
				"error", err,
			)
			return err
		}
		return err
	}
	logger.Info("task complete")

	// dequeue after task run, mark with error if task failed
	logger.Info("dequeueing task")
	_, err = m.DequeueTx(ctx, tx, id)
	if err != nil {
		logger.Error("unable to dequeue item", "id", id, "error", err)
		return err
	}

	commitSuccessful = true
	if commitErr := tx.Commit(ctx); commitErr != nil {
		logger.Error("failed to commit transaction", "commitError", commitErr)
		return commitErr
	}

	logger.Info("task processed and dequeued")

	return nil
}

// Uses a preexisting transaction to select a task and lock a row by it's ID.
// The row cannot be changed while the transaction is active.
func (m *TaskQueueModel) ClaimTx(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
SELECT id,
       name,
       state,
       created_at,
       updated_at,
       run_at,
       task_data
FROM orchestrator.tasks
WHERE id = $1::uuid
	AND run_at >= NOW()
	AND state = 'waiting'
ORDER BY created_at
    FOR UPDATE SKIP LOCKED
LIMIT 1;
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
	err = tx.QueryRow(qCtx, query, id.String()).Scan(
		&task.ID,
		&task.Name,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
		&task.TaskData,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("no rows found", slog.String("taskId", id.String()))
			return nil, ErrRecordNotFound
		default:
			logger.Error("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil

}

// Uses a preexisting transaction to delete a queue row by it's ID.
//
// Differs from a typical delete because it can read a row that has
// been marked by PostgreSQL with `FOR UPDATE SKIP LOCKED`.
func (m *TaskQueueModel) DequeueTx(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
DELETE FROM orchestrator.tasks
WHERE id = $1
RETURNING
    id,
    name,
    state,
    created_at,
    updated_at,
    run_at,
    task_data;
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
	err = tx.QueryRow(qCtx, query, id.String()).Scan(
		&task.ID,
		&task.Name,
		&task.State,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.RunAt,
		&task.TaskData,
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

func (m *TaskQueueModel) UpdateTx(
	ctx context.Context,
	tx pgx.Tx,
	taskQueue TaskQueue,
) (task *TaskQueue, err error) {
	logger := logging.LoggerFromContext(ctx)

	query := `
UPDATE orchestrator.task
SET name = COALESCE($2, name),
	state = COALESCE($2, state),
	created_at = COALESCE($3, created_at),
	run_at = COALESCE($4, run_at),
	task_data = COALESCE($5::JSONB, task_data),
WHERE id = $1
RETURNING
    id,
    name,
    state,
    created_at,
    updated_at,
    run_at,
	task_data;
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
	err = tx.QueryRow(qCtx, query, task.ID).Scan(
		&task.ID,
		&task.Name,
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
			logger.Error("an error occurred while performing query", "error", err)
			return nil, err
		}
	}

	logger.Info("returning task")
	return task, nil
}
