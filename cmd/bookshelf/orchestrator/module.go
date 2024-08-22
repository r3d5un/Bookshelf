package orchestrator

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/r3d5un/Bookshelf/internal/system"
)

const ModuleName string = "orchestrator"

type Module struct {
	logger *slog.Logger
	db     *sql.DB
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module")

	m.logger.Info("injecting database connection")
	m.db = mono.DB()

	return nil
}

func (m *Module) Shutdown() {
	m.logger.Info("shutting down module", slog.String("module", ModuleName))
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", ModuleName)))
}
