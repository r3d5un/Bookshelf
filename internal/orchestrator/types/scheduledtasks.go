package types

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type ScheduledTask struct {
	ID        uuid.UUID  `json:"id"`
	Name      *string    `json:"queue"`
	State     *string    `json:"state"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	RunAt     *time.Time `json:"runAt"`
	TaskData  *string    `json:"task_data,omitempty"`
}

type ScheduledTaskCollection struct {
	CurrentPage  int              `json:"current_page,omitempty"`
	PageSize     int              `json:"page_size,omitempty"`
	FirstPage    int              `json:"first_page,omitempty"`
	LastPage     int              `json:"last_page,omitempty"`
	TotalRecords int              `json:"total_records,omitempty"`
	OrderBy      string           `json:"order_by,omitempty"`
	Data         []*ScheduledTask `json:"data"`
}

func ReadScheduledTask(
	ctx context.Context,
	models *data.Models,
	taskID uuid.UUID,
) (*ScheduledTask, error) {
	tq, err := models.TaskQueues.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	task := ScheduledTask{
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

func ReadAllScheudledTasks(
	ctx context.Context,
	models *data.Models,
	filters data.Filters,
) (tc *ScheduledTaskCollection, err error) {
	tq, metadata, err := models.TaskQueues.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	var tasks []*ScheduledTask
	for _, t := range tq {
		task := ScheduledTask{
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

	tc = &ScheduledTaskCollection{
		CurrentPage:  metadata.CurrentPage,
		PageSize:     metadata.PageSize,
		FirstPage:    metadata.FirstPage,
		LastPage:     metadata.LastPage,
		TotalRecords: metadata.TotalRecords,
		Data:         tasks,
	}

	return tc, nil
}

// ScheduleTask enqueues a new task to the task queue.
//
// NOTE: The ID will be ignored when creating new tasks. The database generates an ID
// at insertion.
func ScheduleTask(
	ctx context.Context,
	models *data.Models,
	newTask ScheduledTask,
) (createdTask *ScheduledTask, err error) {
	newTaskRow := data.TaskQueue{
		Name:     newTask.Name,
		State:    newTask.State,
		RunAt:    newTask.RunAt,
		TaskData: newTask.TaskData,
	}

	insertedTask, err := models.TaskQueues.Insert(ctx, newTaskRow)
	if err != nil {
		return nil, err
	}

	createdTask = &ScheduledTask{
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

func UpdateScheduledTask(
	ctx context.Context,
	models *data.Models,
	newTaskData ScheduledTask,
) (updatedTask *ScheduledTask, err error) {
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

	updatedTask = &ScheduledTask{
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

func DeleteScheduledTask(
	ctx context.Context,
	models *data.Models,
	id uuid.UUID,
) (deletedTask *ScheduledTask, err error) {
	deletedTaskRow, err := models.TaskQueues.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	task := ScheduledTask{
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

// ClaimScheduledTaskByID selects and locks a task from the queue, marking it with the running state
// before returning the task from the queue to the caller. Note that the task is unlocked once
// the task is set to a running state, and is therefore available to other callers to manipulate.
//
// If a task fails, it needs to be set to a failed state in the queue by the caller in a separate
// function call.
func ClaimScheduledTaskByID(
	ctx context.Context,
	models *data.Models,
	taskID uuid.UUID,
) (*ScheduledTask, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("taskId", taskID.String()))

	logger.Info("testing models")
	_, err := models.TaskQueues.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}

	logger.Info("starting transaction")
	tx, err := models.BeginTx(ctx)
	if err != nil {
		logger.Info("unable to start transaction", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	logger.Info("claiming task")
	taskRow, err := models.TaskQueues.ClaimTx(ctx, tx, taskID)
	if err != nil {
		logger.Info("unable to claim transaction")
		return nil, err
	}

	logger.Info("setting task to a running state")
	state := string(data.RunningTaskState)
	taskRow.State = &state

	taskRow, err = models.TaskQueues.UpdateTx(ctx, tx, *taskRow)
	if err != nil {
		logger.Info("unable to set the task to a running state", "error", err)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	logger.Info("task set to running")

	task := ScheduledTask{
		ID:        taskRow.ID,
		Name:      taskRow.Name,
		State:     taskRow.State,
		CreatedAt: taskRow.CreatedAt,
		UpdatedAt: taskRow.UpdatedAt,
		RunAt:     taskRow.RunAt,
		TaskData:  taskRow.TaskData,
	}

	return &task, nil
}

// SetScheduledTaskState selects and locks a task from the queue, setting it with the desired
// state given by the caller.
func SetScheduledTaskState(
	ctx context.Context,
	models *data.Models,
	taskID uuid.UUID,
	state data.TaskState,
) (*ScheduledTask, error) {
	tx, err := models.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	taskRow, err := models.TaskQueues.ClaimTx(ctx, tx, taskID)
	if err != nil {
		return nil, err
	}

	newState := string(state)
	taskRow.State = &newState

	task := ScheduledTask{
		ID:        taskRow.ID,
		Name:      taskRow.Name,
		State:     taskRow.State,
		CreatedAt: taskRow.CreatedAt,
		UpdatedAt: taskRow.UpdatedAt,
		RunAt:     taskRow.RunAt,
		TaskData:  taskRow.TaskData,
	}

	return &task, nil
}

// DequeueScheduledTask selects and locks a task from the queue, then deletes it. The task is returned
// in as it was in it's last state in the task queue.
func DequeueScheduledTask(
	ctx context.Context,
	models *data.Models,
	taskID uuid.UUID,
) (*ScheduledTask, error) {
	tx, err := models.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	taskRow, err := models.TaskQueues.ClaimTx(ctx, tx, taskID)
	if err != nil {
		return nil, err
	}

	taskRow, err = models.TaskQueues.DequeueTx(ctx, tx, taskID)
	if err != nil {
		return nil, err
	}

	task := ScheduledTask{
		ID:        taskRow.ID,
		Name:      taskRow.Name,
		State:     taskRow.State,
		CreatedAt: taskRow.CreatedAt,
		UpdatedAt: taskRow.UpdatedAt,
		RunAt:     taskRow.RunAt,
		TaskData:  taskRow.TaskData,
	}

	return &task, nil
}
