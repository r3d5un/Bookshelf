package data_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

func TestTaskNotificationModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notificationCh := make(chan pgconn.Notification, 1)
	doneCh := make(chan struct{})
	defer close(doneCh)

	readyCh := make(chan struct{})
	go func() {
		close(readyCh)
		models.TaskNotifications.Listen(ctx, notificationCh, doneCh)
	}()

	<-readyCh

	notification := data.TaskNotification{
		ID:    uuid.New(),
		Queue: "test_queue",
	}

	time.Sleep(1 * time.Second)

	err := models.TaskNotifications.Notify(ctx, notification)
	if err != nil {
		t.Errorf("unable to notify: %s", err)
		return
	}

	select {
	case n := <-notificationCh:
		var receivedNotification data.TaskNotification
		err := json.Unmarshal([]byte(n.Payload), &receivedNotification)
		if err != nil {
			t.Errorf("unable to unmarhsal notification: %s\n", err)
			return
		}
		if notification.ID != receivedNotification.ID {
			t.Errorf(
				"expected notification ID %s, got %s\n",
				notification.ID.String(), receivedNotification.ID.String(),
			)
			return
		}
		if notification.Queue != receivedNotification.Queue {
			t.Errorf(
				"expected notification ID %s, got %s\n",
				notification.Queue, receivedNotification.Queue,
			)
			return
		}
	case <-ctx.Done():
		t.Fatal("Did not receive notification in time")
	}
}
