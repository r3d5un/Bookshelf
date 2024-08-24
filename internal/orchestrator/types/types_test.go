package types_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	tt "github.com/r3d5un/Bookshelf/internal/testing"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *pgxpool.Pool
var models *data.Models

func TestMain(m *testing.M) {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	jsonLogger := slog.New(handler)
	slog.SetDefault(jsonLogger)

	dbName := "bookshelf_testing"
	dbUser := "postgres"
	dbPassword := "postgres"

	projectRoot, err := tt.FindProjectRoot()
	if err != nil {
		slog.Error("unable to get project root")
		os.Exit(1)
	}

	migrations, err := tt.ListUpMigrationScrips(fmt.Sprintf("%s/migrations", projectRoot))
	if err != nil {
		slog.Error("unable to list up migrations", "error", err)
		os.Exit(1)
	}

	postgresContainer, err := postgres.Run(
		context.Background(),
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.WithInitScripts(migrations...),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		slog.Error("error occurred while setting up postgres test container", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := postgresContainer.Terminate(context.Background()); err != nil {
			slog.Error("error occurred while terminating up postgres test container", "error", err)
			os.Exit(1)
		}
	}()

	host, err := postgresContainer.Host(context.Background())
	if err != nil {
		slog.Error("unable to get the host from the postgres container", "error", err)
		os.Exit(1)
	}

	port, err := postgresContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		slog.Error("unable to get the host from the postgres container", "error", err)
		os.Exit(1)
	}
	connString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		dbUser, dbPassword, host, port.Port(), dbName,
	)
	slog.Info("DSN", "connString", connString)

	duration := time.Second * 5
	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		slog.Error("unable to parse postgresql pool configuration", "error", err)
		os.Exit(1)
	}
	db, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		slog.Error("unable to create connection pool", "error", err)
		os.Exit(1)
	}
	slog.Info("connection pool established")

	newModels := data.NewModels(db, &duration)
	models = &newModels

	// Run tests
	exitCode := m.Run()
	defer os.Exit(exitCode)
}
