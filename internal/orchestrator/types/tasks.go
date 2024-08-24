package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type Task struct {
	ID        uuid.UUID  `json:"id"`
	Queue     *string    `json:"queue"`
	State     *string    `json:"state"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	RunAt     *time.Time `json:"runAt"`
}

type TaskCollection struct {
	CurrentPage  int     `json:"current_page,omitempty"`
	PageSize     int     `json:"page_size,omitempty"`
	FirstPage    int     `json:"first_page,omitempty"`
	LastPage     int     `json:"last_page,omitempty"`
	TotalRecords int     `json:"total_records,omitempty"`
	OrderBy      string  `json:"order_by,omitempty"`
	Data         []*Task `json:"data"`
}

func ReadTask(ctx context.Context, models *data.Models, taskID uuid.UUID) (*Task, error) {
	tq, err := models.TaskQueues.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task := Task{
		ID:        tq.ID,
		Queue:     tq.Queue,
		State:     tq.State,
		CreatedAt: tq.CreatedAt,
		UpdatedAt: tq.UpdatedAt,
		RunAt:     tq.RunAt,
	}

	return &task, nil
}

func ReadAllTasks(
	ctx context.Context,
	models *data.Models,
	filters data.Filters,
) (tc *TaskCollection, err error) {
	tq, metadata, err := models.TaskQueues.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	var tasks []*Task
	for _, t := range tq {
		task := Task{
			ID:        t.ID,
			Queue:     t.Queue,
			State:     t.State,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
			RunAt:     t.RunAt,
		}

		tasks = append(tasks, &task)
	}

	tc = &TaskCollection{
		CurrentPage:  metadata.CurrentPage,
		PageSize:     metadata.PageSize,
		FirstPage:    metadata.FirstPage,
		LastPage:     metadata.LastPage,
		TotalRecords: metadata.TotalRecords,
		Data:         tasks,
	}

	return tc, nil
}

func CreateTask(
	ctx context.Context,
	models *data.Models,
	newTask Task,
) (createdTask *Task, err error) {
	insertedTask, err := models.TaskQueues.Insert(ctx, *newTask.Queue, newTask.State, newTask.RunAt)
	if err != nil {
		return nil, err
	}

	createdTask = &Task{
		ID:        insertedTask.ID,
		Queue:     insertedTask.Queue,
		State:     insertedTask.State,
		CreatedAt: insertedTask.CreatedAt,
		UpdatedAt: insertedTask.UpdatedAt,
		RunAt:     insertedTask.RunAt,
	}

	return createdTask, nil
}

func UpdateTask(
	ctx context.Context,
	models *data.Models,
	newTaskData Task,
) (updatedTask *Task, err error) {
	newTaskRow := data.TaskQueue{
		ID:        newTaskData.ID,
		Queue:     newTaskData.Queue,
		State:     newTaskData.State,
		CreatedAt: newTaskData.CreatedAt,
		UpdatedAt: newTaskData.UpdatedAt,
		RunAt:     newTaskData.RunAt,
	}
	updatedTaskRow, err := models.TaskQueues.Update(ctx, newTaskRow)
	if err != nil {
		return nil, err
	}

	updatedTask = &Task{
		ID:        updatedTaskRow.ID,
		Queue:     updatedTaskRow.Queue,
		State:     updatedTaskRow.State,
		CreatedAt: updatedTaskRow.CreatedAt,
		UpdatedAt: updatedTaskRow.UpdatedAt,
		RunAt:     updatedTaskRow.RunAt,
	}

	return updatedTask, nil
}
