package types

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type Task struct {
	ID        uuid.UUID  `json:"id"`
	Name      *string    `json:"queue"`
	State     *string    `json:"state"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	RunAt     *time.Time `json:"runAt"`
	TaskData  *string    `json:"task_data,omitempty"`
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
		Name:      tq.Name,
		State:     tq.State,
		CreatedAt: tq.CreatedAt,
		UpdatedAt: tq.UpdatedAt,
		RunAt:     tq.RunAt,
		TaskData:  tq.TaskData,
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
			Name:      t.Name,
			State:     t.State,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
			RunAt:     t.RunAt,
			TaskData:  t.TaskData,
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
	insertedTask, err := models.TaskQueues.Insert(
		ctx, *newTask.Name, newTask.State, newTask.RunAt, nil,
	)
	if err != nil {
		return nil, err
	}

	createdTask = &Task{
		ID:        insertedTask.ID,
		Name:      insertedTask.Name,
		State:     insertedTask.State,
		CreatedAt: insertedTask.CreatedAt,
		UpdatedAt: insertedTask.UpdatedAt,
		RunAt:     insertedTask.RunAt,
		TaskData:  insertedTask.TaskData,
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
		Name:      newTaskData.Name,
		State:     newTaskData.State,
		CreatedAt: newTaskData.CreatedAt,
		UpdatedAt: newTaskData.UpdatedAt,
		RunAt:     newTaskData.RunAt,
		TaskData:  newTaskData.TaskData,
	}
	updatedTaskRow, err := models.TaskQueues.Update(ctx, newTaskRow)
	if err != nil {
		return nil, err
	}

	updatedTask = &Task{
		ID:        updatedTaskRow.ID,
		Name:      updatedTaskRow.Name,
		State:     updatedTaskRow.State,
		CreatedAt: updatedTaskRow.CreatedAt,
		UpdatedAt: updatedTaskRow.UpdatedAt,
		RunAt:     updatedTaskRow.RunAt,
		TaskData:  updatedTaskRow.TaskData,
	}

	return updatedTask, nil
}

func DeleteTask(
	ctx context.Context,
	models *data.Models,
	id uuid.UUID,
) (deletedTask *Task, err error) {
	deletedTaskRow, err := models.TaskQueues.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	task := Task{
		ID:        deletedTaskRow.ID,
		Name:      deletedTaskRow.Name,
		State:     deletedTaskRow.State,
		CreatedAt: deletedTaskRow.CreatedAt,
		UpdatedAt: deletedTaskRow.UpdatedAt,
		RunAt:     deletedTaskRow.RunAt,
		TaskData:  deletedTaskRow.TaskData,
	}

	return &task, nil
}
