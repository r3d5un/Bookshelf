package orchestrator

import (
	"context"
	"errors"
)

var (
	ErrNoTask = errors.New("task does not exist")
)

type Task func(context.Context) error

type Collection map[string]Task

// Run is executes a function that matches the given name, injecting the
// given context as a function parameter.
func (otc *Collection) Run(ctx context.Context, taskName string) error {
	if task, exists := (*otc)[taskName]; exists {
		return task(ctx)
	}

	return ErrNoTask
}

// Adds a new task to the task collection
func (otc *Collection) Add(taskName string, newTask Task) {
	(*otc)[taskName] = newTask
}

// Get a task function without executing it
func (otc *Collection) Get(taskName string) (ot *Task, err error) {
	if task, exists := (*otc)[taskName]; exists {
		return &task, nil
	}

	return nil, ErrNoTask
}
