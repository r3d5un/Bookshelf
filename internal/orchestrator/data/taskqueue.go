package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
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

func (m *TaskQueueModel) Listen(
	ctx context.Context,
	notificationCh chan<- pgconn.Notification,
	done <-chan bool,
) {
	logger := logging.LoggerFromContext(ctx)

	query := `
LISTEN task_queue_notification;
`

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
		),
	)

	for {
		conn, err := m.Pool.Acquire(ctx)
		if err != nil {
			logger.Error("unable to acquire connection", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Info("listening for task notifications")
		_, err = conn.Exec(context.Background(), query)
		if err != nil {
			logger.Error("unable to listen", "error", err)
			conn.Release()
			continue
		}

		for {
			select {
			case <-done:
				logger.Info("received done signal, stopping listener.")
				return
			default:
				if notification, err := conn.Conn().WaitForNotification(ctx); err != nil {
					logger.Info("unable to receive notification", "error", err)
					break
				} else {
					select {
					case notificationCh <- *notification:
						logger.Info("notification sent to channel")
					case <-done:
						logger.Info("received done signal, stopping listener.")
						return
					default:
						logger.Error("channel full, notification not sent")
						time.Sleep(5 * time.Second)
						continue
					}
				}
			}
		}

	}
}
