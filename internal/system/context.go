package system

import (
	"context"

	"github.com/google/uuid"
)

// InstanceFromContext returns the UUID from the given context. If
// instanceID is not set, or is not a UUID, a new UUID is returned.
func InstanceFromContext(ctx context.Context) uuid.UUID {
	if instanceID, ok := ctx.Value("instanceID").(uuid.UUID); ok {
		return instanceID
	}

	return uuid.New()
}
