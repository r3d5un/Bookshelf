package orchestrator

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/config"
	"github.com/r3d5un/Bookshelf/internal/system"
)

const ModuleName string = "orchestrator"

type Module struct {
	logger *slog.Logger
	cfg    *config.Config
	db     *pgxpool.Pool
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module")

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

	return nil
}

func (m *Module) Shutdown() {
	m.logger.Info("shutting down module")

	m.logger.Info("closing module connection pool")
	m.db.Close()

	m.logger.Info("module shutdown complete")
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", ModuleName)))
}
