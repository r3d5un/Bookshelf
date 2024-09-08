package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

var (
	DuplicateKeyMatchString = "cannot insert duplicate key object"
)

type Models struct {
	Tasks             TaskModel
	TaskQueues        TaskQueueModel
	TaskNotifications TaskNotificationModel
	SchedulerLock     SchedulerLockModel
	pool              *pgxpool.Pool
}

func NewModels(pool *pgxpool.Pool, timeout *time.Duration) Models {
	return Models{
		Tasks:             TaskModel{Pool: pool, Timeout: timeout},
		TaskQueues:        TaskQueueModel{Pool: pool, Timeout: timeout},
		TaskNotifications: TaskNotificationModel{Pool: pool, Timeout: timeout},
		SchedulerLock:     SchedulerLockModel{Pool: pool, Timeout: timeout},
		pool:              pool,
	}
}
func (m *Models) BeginTx(ctx context.Context) (tx pgx.Tx, err error) {
	return m.pool.Begin(ctx)
}
