package orchestrator

import (
	"context"
	"errors"
)

var (
	ErrNoTask = errors.New("task does not exist")
)

type OrchestratorTask func(context.Context) error

type OrchestratorTaskCollection map[string]OrchestratorTask

// Run is executes a function that matches the given name, injecting the
// given context as a function parameter.
func (otc *OrchestratorTaskCollection) Run(ctx context.Context, taskName string) error {
	if task, exists := (*otc)[taskName]; exists {
		return task(ctx)
	}

	return ErrNoTask
}

// Adds a new task to the task collection
func (otc *OrchestratorTaskCollection) Add(taskName string, newTask OrchestratorTask) {
	(*otc)[taskName] = newTask
}

// Get a task function without executing it
func (otc *OrchestratorTaskCollection) Get(taskName string) (ot *OrchestratorTask, err error) {
	if task, exists := (*otc)[taskName]; exists {
		return &task, nil
	}

	return nil, ErrNoTask
}
