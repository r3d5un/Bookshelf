package types

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

type TaskLog struct {
	ID     uuid.UUID `json:"id"`
	TaskID uuid.UUID `json:"taskId"`
	Log    string    `json:"log"`
}

type TaskLogWriter struct {
	taskID    uuid.UUID
	Done      chan struct{}
	logBuffer chan TaskLog
	models    *data.Models
	wg        *sync.WaitGroup
}

func NewTaskLogWriter(
	ctx context.Context,
	models *data.Models,
	taskID uuid.UUID,
	logBufferSize int,
) TaskLogWriter {
	var wg sync.WaitGroup

	logWriter := TaskLogWriter{
		taskID:    taskID,
		logBuffer: make(chan TaskLog, logBufferSize),
		Done:      make(chan struct{}, 1),
		models:    models,
		wg:        &wg,
	}

	go func() {
		logWriter.LogSink(ctx)
	}()

	return logWriter
}

func (tlw *TaskLogWriter) Write(p []byte) (n int, err error) {
	log := TaskLog{
		ID:     uuid.New(),
		TaskID: tlw.taskID,
		Log:    string(p),
	}

	tlw.wg.Add(1)
	go func() {
		defer tlw.wg.Done()
		tlw.logBuffer <- log
	}()

	return len(p), nil
}

func (tlw *TaskLogWriter) LogSink(ctx context.Context) {
	for {
		select {
		case <-tlw.Done:
			return
		case log, ok := <-tlw.logBuffer:
			if !ok {
				slog.Error("unable to read logs from log buffer")
				return
			}
			_, err := CreateTaskLog(ctx, tlw.models, log)
			if err != nil {
				slog.Error("unable to create log records", "error", err)
				return
			}
		}
	}
}

func (tlw *TaskLogWriter) Stop() {
	tlw.wg.Wait()
	close(tlw.logBuffer)
	close(tlw.Done)
}

func CreateTaskLog(ctx context.Context, models *data.Models, log TaskLog) (*TaskLog, error) {
	logRow := data.TaskLog{
		ID:     log.ID,
		TaskID: log.TaskID,
		Log:    log.Log,
	}

	insertedLogRow, err := models.TaskLogs.Insert(ctx, logRow)
	if err != nil {
		return nil, err
	}

	createdLog := TaskLog{
		ID:     insertedLogRow.ID,
		TaskID: insertedLogRow.TaskID,
		Log:    insertedLogRow.Log,
	}

	return &createdLog, nil
}

func ReadLogsByTaskQueueID(
	ctx context.Context,
	models *data.Models,
	taskID uuid.UUID,
) ([]*TaskLog, error) {
	logRows, err := models.TaskLogs.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	logs := make([]*TaskLog, len(logRows))

	for _, logRow := range logRows {
		log := TaskLog{
			ID:     logRow.ID,
			TaskID: logRow.TaskID,
			Log:    logRow.Log,
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

// NewTaskLogger creates a new logger that writes logs to stdout and the task log
// database table. It also returns a stop function, which stop associated goroutines
// that writes to the database in a non-blocking manner.
func NewTaskLogger(
	ctx context.Context,
	models *data.Models,
	taskName string,
	taskQueueID uuid.UUID,
) (logger *slog.Logger, stop func()) {
	logWriter := NewTaskLogWriter(ctx, models, taskQueueID, 100)

	handler := slog.NewJSONHandler(io.MultiWriter(&logWriter, os.Stdout), nil)
	logger = slog.New(handler).With(slog.Group(
		"task",
		slog.String("taskName", taskName),
		slog.String("taskId", taskQueueID.String())),
	)

	return logger, logWriter.Stop
}
