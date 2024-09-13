package data_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestTaskLogModel(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		_, err := models.TaskLogs.Get(context.Background(), uuid.New())
		if err != nil {
			t.Errorf("unable to get task log: %s\n", err)
			return
		}
	})
}
