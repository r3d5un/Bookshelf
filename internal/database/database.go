package database

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"time"
)

func OpenPool(
	connString string,
	maxOpenConns int,
	maxIdleConns int,
	maxIdleTime string,
	timeout time.Duration,
) (*sql.DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MinifySQL(query string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(query, " "))
}

func CreateOrderByClause(orderBy []string) string {
	if len(orderBy) == 0 {
		return "ORDER BY id"
	}

	var orderClauses []string
	for _, item := range orderBy {
		if strings.HasPrefix(item, "-") {
			orderClauses = append(orderClauses, strings.TrimPrefix(item, "-")+" DESC")
		} else {
			orderClauses = append(orderClauses, item+" ASC")
		}
	}

	return "ORDER BY " + strings.Join(orderClauses, ", ")
}
