package main

import (
	"context"
	_ "database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"time"
)

const (
	DatabaseHost     = "DATABASE_HOST"
	DatabaseName     = "DATABASE_NAME"
	DatabaseUser     = "DATABASE_USER"
	DatabasePassword = "DATABASE_PASSWORD"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			go func() {
				err := checkDatabaseAccess(ctx)
				if err != nil {
					logrus.WithError(err).Error("Database query loop failed")
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}

func checkDatabaseAccess(ctx context.Context) error {
	connectionString := buildConnectionString()
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return err
	}

	logrus.Info("Successfully connected to database")
	defer conn.Close(ctx)
	// Test INSERT permissions
	insertStatement := `INSERT INTO example (entry_time) VALUES ($1)`
	now := time.Now()
	_, err = conn.Exec(ctx, insertStatement, now.UnixNano())
	if err != nil {
		logrus.Info("Unable to perform INSERT operation")
	} else {
		logrus.Info("Successfully INSERTED into our table")
	}

	// Test SELECT permissions
	selectStatement := `SELECT COUNT(*), MAX(entry_time) FROM example limit 3;`
	row := conn.QueryRow(ctx, selectStatement)

	var count int
	var maxTimestamp int64

	err = row.Scan(&count, &maxTimestamp)
	if err != nil {
		logrus.Info("Unable to perform SELECT operation")
	} else {
		mostRecentDate := time.Unix(0, maxTimestamp).Format(time.RFC3339)
		logrus.Info("Successfully SELECTED, most recent value: ", mostRecentDate)
	}

	return nil
}

func buildConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		os.Getenv(DatabaseUser),
		url.QueryEscape(os.Getenv(DatabasePassword)),
		os.Getenv(DatabaseHost),
		os.Getenv(DatabaseName))
}
