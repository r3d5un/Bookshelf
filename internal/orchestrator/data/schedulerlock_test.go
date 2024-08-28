package data_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestSchedulerLock(t *testing.T) {
	instanceID := uuid.New()

	t.Run("AcquireLock", func(t *testing.T) {
		lockAcquired, err := models.SchedulerLock.AcquireLock(context.Background(), instanceID)
		if err != nil {
			t.Errorf("an error occurred while acquiring the scheduler lock: %s\n", err)
			return
		}

		if !lockAcquired {
			t.Errorf("unable to acquire scheduler lock: %t", lockAcquired)
			return
		}
	})

	t.Run("MaintainLock", func(t *testing.T) {
		err := models.SchedulerLock.MaintainLock(context.Background(), instanceID)
		if err != nil {
			t.Errorf("unable to maintain scheduler lock: %s\n", err)
			return
		}
	})
}
