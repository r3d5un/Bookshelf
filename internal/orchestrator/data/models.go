package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

var (
	DplicateKeyMatchString = "cannot insert duplicate key object"
)

type Models struct {
	TaskQueues TaskQueueModel
}

func NewModels(db *sql.DB, timeout *time.Duration) Models {
	return Models{
		TaskQueues: TaskQueueModel{DB: db, Timeout: timeout},
	}
}
