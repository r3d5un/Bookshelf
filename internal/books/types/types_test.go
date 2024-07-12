package types_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
	"github.com/r3d5un/Bookshelf/internal/database"
	tt "github.com/r3d5un/Bookshelf/internal/testing"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB
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
	db, err = database.OpenPool(connString, 15, 15, "15m", duration)
	if err != nil {
		slog.Error("unable to open the database connection pool", "error", err)
		os.Exit(1)
	}

	newModels := data.NewModels(db, &duration)
	models = &newModels

	// Run tests
	exitCode := m.Run()
	defer os.Exit(exitCode)
}

func TestComplexBookTypes(t *testing.T) {
	description := "This is a test description."
	timestamp := time.Now()

	bookRecord := data.Book{
		ID:          uuid.New(),
		Title:       "TestBookTitle",
		Description: &description,
		Published:   &timestamp,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	genreRecord := data.Genre{
		ID:          uuid.New(),
		Name:        "TestGenreName",
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}
	_, err := models.Genres.Insert(context.Background(), genreRecord)
	if err != nil {
		t.Errorf("unable to insert genre: %s\n", err)
		return
	}

	seriesRecord := data.Series{
		ID:          uuid.New(),
		Name:        "TestSeriesName",
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}
	_, err = models.Series.Insert(context.Background(), seriesRecord)
	if err != nil {
		t.Errorf("unable to insert series: %s\n", err)
		return
	}

	authorName := "TestAuthorName"
	website := "www.testwebsite.com"
	authorRecord := data.Author{
		ID:          uuid.New(),
		Name:        &authorName,
		Description: &description,
		Website:     &website,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}
	_, err = models.Authors.Insert(context.Background(), authorRecord)
	if err != nil {
		t.Errorf("unable to insert series: %s\n", err)
		return
	}

	book := types.Book{
		ID:          nil,
		Title:       &bookRecord.Title,
		Description: bookRecord.Description,
		Published:   bookRecord.Published,
		CreatedAt:   bookRecord.CreatedAt,
		UpdatedAt:   bookRecord.UpdatedAt,
		Authors:     []*data.Author{&authorRecord},
		Genres:      []*data.Genre{&genreRecord},
		BookSeries: []*data.BookSeries{
			{
				BookID:      bookRecord.ID,
				SeriesID:    seriesRecord.ID,
				SeriesOrder: 1.0,
			},
		},
	}

	var insertedBook *data.Book

	t.Run("TestCreateBook", func(t *testing.T) {
		_, err := types.CreateBook(context.Background(), models, book)
		if err != nil {
			t.Errorf("error occurred when registering new book: %s\n", err)
			return
		}
	})

	t.Run("TestReadBook", func(t *testing.T) {
		bookRecord := data.Book{
			ID:          uuid.New(),
			Title:       "TestGetBookTitle",
			Description: &description,
			Published:   &timestamp,
			CreatedAt:   &timestamp,
			UpdatedAt:   &timestamp,
		}
		insertedBook, err = models.Books.Insert(context.Background(), bookRecord)
		if err != nil {
			t.Errorf("unable to insert book: %s\n", err)
			return
		}
		if _, err := types.ReadBook(context.Background(), models, bookRecord.ID); err != nil {
			t.Errorf("error occurred while retrieving book: %s\n", err)
			return
		}
	})

	t.Run("TestReadAllBooks", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 10,
		}
		bookList, err := types.ReadAllBooks(context.Background(), models, filters)
		if err != nil {
			t.Errorf("unable to read books: %s\n", err)
			return
		}
		if len(bookList) < 1 {
			t.Error("no books returned")
			return
		}
	})

	t.Run("TestGetNonExistingBook", func(t *testing.T) {
		if _, err := types.ReadBook(context.Background(), models, uuid.New()); err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				// desired result. caller is tasked with handling the ErrRecordNotFound error
				return
			default:
				t.Errorf("error occurred while retrieving book: %s\n", err)
				return
			}
		}
	})

	t.Run("TestDeleteBook", func(t *testing.T) {
		if err := types.DeleteBook(context.Background(), models, insertedBook.ID); err != nil {
			t.Errorf("unable to delete book: %s\n", err)
			return
		}
	},
	)
}
