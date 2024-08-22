package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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
	DB      *sql.DB
	Timeout *time.Duration
}
