package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	WaitingTaskState  string = "waiting"
	RunningTaskState  string = "running"
	CompleteTaskState string = "complete"
	StoppedTaskState  string = "stopped"
)

type TaskQueue struct {
	ID        uuid.UUID `json:"id"`
	Queue     string    `json:"queue"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RunAt     time.Time `json:"runAt"`
}

type TaskQueueModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}
}
