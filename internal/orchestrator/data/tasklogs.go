package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskLog struct {
	ID  uuid.UUID      `json:"name"`
	Log sql.NullString `json:"log"`
}

type TaskLogModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}
