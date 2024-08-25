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
	TaskQueues        TaskQueueModel
	TaskNotifications TaskNotificationModel
	pool              *pgxpool.Pool
}

func NewModels(pool *pgxpool.Pool, timeout *time.Duration) Models {
	return Models{
		TaskQueues:        TaskQueueModel{Pool: pool, Timeout: timeout},
		TaskNotifications: TaskNotificationModel{Pool: pool, Timeout: timeout},
	}
}
func (m *Models) BeginTx(ctx context.Context) (tx pgx.Tx, err error) {
	return m.pool.Begin(ctx)
}
