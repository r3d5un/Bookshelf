package orchestrator

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/config"
	"github.com/r3d5un/Bookshelf/internal/orchestrator"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/system"
)

const ModuleName string = "orchestrator"

type Module struct {
	schedulerID         uuid.UUID
	logger              *slog.Logger
	cfg                 *config.Config
	db                  *pgxpool.Pool
	models              data.Models
	scheduler           *orchestrator.CronScheduler
	done                chan struct{}
	taskNotificationCh  chan pgconn.Notification
	taskCollection      orchestrator.Collection
	wg                  sync.WaitGroup
	isSchedulerMasterCh chan bool
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module")

	m.logger.Info("setting scheduler ID from instance")
	m.schedulerID = system.InstanceFromContext(mono.Context())
	m.logger.Info("scheduler instance ID set", "id", m.schedulerID)

	m.logger.Info("injecting database connection")
	m.cfg = mono.Config()

	dbConfig, err := pgxpool.ParseConfig(m.cfg.DB.DSN)
	if err != nil {
		m.logger.Error("unable to parse postgresql pool configuration", "error", err)
		return nil
	}
	m.db, err = pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		m.logger.Error("unable to create connection pool", "error", err)
		return nil
	}
	m.logger.Info("connection pool established")

	m.logger.Info("initializing channels")
	m.wg = sync.WaitGroup{}
	m.done = make(chan struct{})
	m.taskNotificationCh = make(chan pgconn.Notification, 100)
	m.isSchedulerMasterCh = make(chan bool, 1)

	timeout := time.Duration(m.cfg.DB.Timeout) * time.Second
	m.models = data.NewModels(m.db, &timeout)

	m.logger.Info("creating task collection")
	m.taskCollection = orchestrator.Collection{}

	m.logger.Info("creating task runner")
	m.wg.Add(1)
	go m.taskRunner(ctx)

	m.logger.Info("creating task scheduler")
	m.scheduler = orchestrator.NewScheduler(&m.models)

	m.logger.Info("adding tasks")
	err = m.addTasks(ctx)
	if err != nil {
		m.logger.Error("unable to add tasks", "error", err)
		return
	}

	m.wg.Add(1)
	go func() { // Goroutine for checking the scheduler lock
		defer m.wg.Done()
		m.checkSchedulerLock(ctx)
	}()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.manageScheduler(ctx)
	}()

	// TODO: Handle due/stale tasks in queue without notifications

	m.logger.Info("startup complete")

	return nil
}

func (m *Module) Shutdown() {
	m.logger.Info("shutting down module")

	m.logger.Info("stopping scheduler")
	m.scheduler.Stop()
	m.wg.Done()

	m.logger.Info("closing notification channel")
	close(m.taskNotificationCh)

	m.logger.Info("sending stop signal")
	close(m.done)
	close(m.isSchedulerMasterCh)

	m.logger.Info("waiting for scheduler background processes to complete")
	m.wg.Wait()

	// TODO: Listener refusing to let go of connection. Hang causes unclean shutdown.
	m.logger.Info("closing module connection pool")
	activeConns := m.db.Stat().AcquiredConns()
	if activeConns > 0 {
		m.logger.Error(
			"acquired connections not released, unclean shutdown",
			"activeConnections", activeConns,
		)
	} else {
		m.db.Close()
	}

	m.logger.Info("module shutdown complete")
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", ModuleName)))
}
