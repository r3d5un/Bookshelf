package data

import (
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	Name      sql.NullString `json:"name"`
	CronExpr  sql.NullString `json:"cronExpr"`
	Enabled   sql.NullBool   `json:"enabled"`
	Deleted   sql.NullBool   `json:"deleted"`
	Timestamp sql.NullTime   `json:"timestamp"`
}

type TaskModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}
