package orchestrator

import (
	"context"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

func (m *Module) ReadScheduledTask(
	ctx context.Context,
	taskID uuid.UUID,
) (*types.ScheduledTask, error) {
	return types.ReadScheduledTask(ctx, &m.models, taskID)
}

func (m *Module) ReadAllScheduledTasks(
	ctx context.Context,
	filters data.Filters,
) (*types.TaskCollection, error) {
	return types.ReadAllScheudledTasks(ctx, &m.models, filters)
}

func (m *Module) CreateScheduledTask(
	ctx context.Context,
	newTask types.ScheduledTask,
) (*types.ScheduledTask, error) {
	return types.ScheduleTask(ctx, &m.models, newTask)
}

func (m *Module) UpdateScheduledTask(
	ctx context.Context,
	newTaskData types.ScheduledTask,
) (*types.ScheduledTask, error) {
	return types.UpdateScheduledTask(ctx, &m.models, newTaskData)
}

func (m *Module) DeleteScheduledTask(
	ctx context.Context,
	id uuid.UUID,
) (*types.ScheduledTask, error) {
	return types.DeleteScheduledTask(ctx, &m.models, id)
}

func (m *Module) ClaimScheduledTaskByID(
	ctx context.Context,
	taskID uuid.UUID,
) (*types.ScheduledTask, error) {
	return types.ClaimScheduledTaskByID(ctx, &m.models, taskID)
}

func (m *Module) SetScheduledTaskState(
	ctx context.Context,
	taskID uuid.UUID,
	state data.TaskState,
) (*types.ScheduledTask, error) {
	return types.SetScheduledTaskState(ctx, &m.models, taskID, state)
}

func (m *Module) DequeueScheduledTask(
	ctx context.Context,
	taskID uuid.UUID,
) (*types.ScheduledTask, error) {
	return types.DequeueScheduledTask(ctx, &m.models, taskID)
}
