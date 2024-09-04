package orchestrator

import (
	"context"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

func (m *Module) ReadTask(ctx context.Context, taskID uuid.UUID) (*types.Task, error) {
	return types.ReadTask(ctx, &m.models, taskID)
}

func (m *Module) ReadAllTasks(
	ctx context.Context,
	filters data.Filters,
) (*types.TaskCollection, error) {
	return types.ReadAllTasks(ctx, &m.models, filters)
}

func (m *Module) CreateTask(ctx context.Context, newTask types.Task) (*types.Task, error) {
	return types.CreateTask(ctx, &m.models, newTask)
}

func (m *Module) UpdateTask(ctx context.Context, newTaskData types.Task) (*types.Task, error) {
	return types.UpdateTask(ctx, &m.models, newTaskData)
}

func (m *Module) DeleteTask(ctx context.Context, id uuid.UUID) (*types.Task, error) {
	return types.DeleteTask(ctx, &m.models, id)
}

func (m *Module) ClaimTaskByID(ctx context.Context, taskID uuid.UUID) (*types.Task, error) {
	return types.ClaimTaskByID(ctx, &m.models, taskID)
}

