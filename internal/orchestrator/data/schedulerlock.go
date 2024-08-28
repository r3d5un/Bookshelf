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

type SchedulerLock struct {
	ID            uuid.UUID `json:"id"`
	InstanceID    uuid.UUID `json:"instanceId"`
	LastHeartbeat time.Time `json:"lastHeartbeat"`
}

type SchedulerLockModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}

func (m *SchedulerLockModel) AcquireLock(ctx context.Context, instanceID uuid.UUID) (bool, error) {
	query := `
UPDATE orchestrator.scheduler_lock
SET instance_id    = $1,
    last_heartbeat = NOW()
WHERE (instance_id = '' OR last_heartbeat < NOW() - INTERVAL '10 seconds')
RETURNING
    id,
    instance_id,
    last_heartbeat;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", database.MinifySQL(query)),
		slog.String("instanceId", instanceID.String())),
	)

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var acquiredLock SchedulerLock
	err := m.Pool.QueryRow(qCtx, query, instanceID.String()).
		Scan(&acquiredLock.ID, &acquiredLock.InstanceID, &acquiredLock.LastHeartbeat)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("unable to acquire lock, lock updated within the last 10 seconds")
			return false, nil
		default:
			logger.Error("error occurred while performing query", "error", err)
			return false, err
		}
	}
	if acquiredLock.InstanceID != instanceID {
		logger.Info(
			"unable to acquire lock, returend ID does not match instance ID",
			"instanceId", instanceID,
			"acquiredLock", acquiredLock,
		)
	}

	logger.Info("lock acquired", "acquiredLock", acquiredLock)
	return true, nil
}

func (m *SchedulerLockModel) MaintainLock(ctx context.Context, instanceID uuid.UUID) error {
	query := `
UPDATE orchestrator.scheduler_lock
SET last_heartbeat = NOW()
WHERE instance_id = $1;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", database.MinifySQL(query)),
		slog.String("instanceId", instanceID.String())),
	)

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	res, err := m.Pool.Exec(qCtx, query, instanceID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			logger.Info("unable to acquire lock, lock updated within the last 10 seconds")
			return nil
		default:
			logger.Error("error occurred while performing query", "error", err)
			return err
		}
	}

	logger.Info("query completed", "rowsAffected", res.RowsAffected())

	return nil
}
