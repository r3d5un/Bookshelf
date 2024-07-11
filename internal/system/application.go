package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/r3d5un/Bookshelf/internal/config"
)

type MonolithApplication struct {
	logger  *slog.Logger
	mux     *http.ServeMux
	modules []Module
	db      *sql.DB
	cfg     *config.Config
}

func (app *MonolithApplication) Logger() *slog.Logger {
	return app.logger
}

func (app *MonolithApplication) Mux() *http.ServeMux {
	return app.mux
}

func (app *MonolithApplication) DB() *sql.DB {
	return app.db
}

func (app *MonolithApplication) Config() *config.Config {
	return app.cfg
}

func (app *MonolithApplication) Modules() []Module {
	return app.modules
}

func NewMonolith(
	logger *slog.Logger,
	mux *http.ServeMux,
	modules []Module,
	db *sql.DB,
	cfg *config.Config,
) MonolithApplication {
	return MonolithApplication{
		logger:  logger,
		mux:     mux,
		modules: modules,
		db:      db,
		cfg:     cfg,
	}
}

func (app *MonolithApplication) Serve() error {
	// TODO: Set port and environment through configuration
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 4000),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}

func (app *MonolithApplication) routes() http.Handler {
	app.logger.Info("creating standard middleware chain")
	// TODO: Create middleware: recoverPanic
	standard := alice.New(app.logRequest)

	handler := standard.Then(app.Mux())
	return handler
}

func (app *MonolithApplication) SetupModules(ctx context.Context) error {
	for _, v := range app.modules {
		if err := v.Startup(ctx, app); err != nil {
			return err
		}
	}

	return nil
}

func (app *MonolithApplication) ShutdownModules() error {
	for _, v := range app.modules {
		v.Shutdown()
	}

	return nil
}
