package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
