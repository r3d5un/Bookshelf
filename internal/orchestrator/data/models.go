package data

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

var (
	DuplicateKeyMatchString = "cannot insert duplicate key object"
)

type Models struct {
	TaskQueues TaskQueueModel
}

func NewModels(pool *pgxpool.Pool, timeout *time.Duration) Models {
	return Models{
		TaskQueues: TaskQueueModel{Pool: pool, Timeout: timeout},
	}
}
