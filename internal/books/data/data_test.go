package data_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

func TestMain(m *testing.M) {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	jsonLogger := slog.New(handler)
	slog.SetDefault(jsonLogger)

	dbName := "bookshelf_testing"
	dbUser := "postgres"
	dbPassword := "postgres"

	migrations, err := listUpMigrationScrips(
		"/home/r3d5un/Development/Projects/Bookshelf/migrations",
	)
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
	db, err = database.OpenPool(connString, 15, 15, "15m", duration)
	if err != nil {
		slog.Error("unable to open the database connection pool", "error", err)
		os.Exit(1)
	}
}

func listUpMigrationScrips(dirPath string) (migrations []string, err error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			migrations = append(migrations, fmt.Sprintf("%s/%s", dirPath, file.Name()))
		}
	}

	return migrations, nil
}
