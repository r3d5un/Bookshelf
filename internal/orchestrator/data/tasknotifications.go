package data

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/logging"
)

// TaskNotification is used by the TaskNotificationModel to create a payload for the notification.
type TaskNotification struct {
	// ID task causing the notification
	ID uuid.UUID `json:"id"`
	// Queue refers to the task causing the notification
	Queue string `json:"queue"`
}

type TaskNotificationModel struct {
	Timeout *time.Duration
	Pool    *pgxpool.Pool
}

// Notify sends a notification on the task_queue_notification PostgreSQL channel. Takes in a TaskNotification
// which is encoded to JSON as a payload for the notification.
func (m *TaskNotificationModel) Notify(ctx context.Context, notification TaskNotification) error {
	logger := logging.LoggerFromContext(ctx)

	logger.Info("jsonifying notification", "notification", notification)
	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	query := "SELECT pg_notify($1, $2);"

	logger = logger.With(
		slog.Group(
			"query",
			slog.String("statement", database.MinifySQL(query)),
		),
	)
	conn, err := m.Pool.Acquire(context.Background())
	if err != nil {
		logger.Error("unable to acquire connection", "error", err)
		return err
	}
	defer conn.Release()

	logger.Info("executing notification statement")
	_, err = conn.Exec(context.Background(), query, "task_queue_notification", string(payload))
	if err != nil {
		logger.Error("unable to execute notification statement", "error", err)
		return err
	}

	return nil
}

// Listen for notifications on the task_queue_notification PostgreSQL channel. The listener
// takes a notification channel that the caller can read from, getting notifications as soon
// as they are received, and trigger tasks as needed. The payload can be decoded to a TaskNotification.
//
// A done channel is needed to perform a clean shutdown.
func (m *TaskNotificationModel) Listen(
	ctx context.Context,
	notificationCh chan<- pgconn.Notification,
	done <-chan bool,
) {
	logger := logging.LoggerFromContext(ctx)

	query := `LISTEN task_queue_notification;`

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
		defer conn.Release()

		logger.Info("listening for task notifications")
		_, err = conn.Exec(ctx, query)
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
