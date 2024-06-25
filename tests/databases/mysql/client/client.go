package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"log"
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
	if os.Getenv(DatabaseHost) == "" || os.Getenv(DatabaseName) == "" || os.Getenv(DatabaseUser) == "" || os.Getenv(DatabasePassword) == "" {
		logrus.Error("Missing environment variables")
		return fmt.Errorf("missing environment variables")
	}

	connectionString := buildConnectionString()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Println(err)
		return err
	}

	logrus.Info("Successfully connected to database")
	defer db.Close()

	// Test INSERT permissions
	insertStatement := `INSERT INTO example (entry_time) VALUES (?);`
	now := time.Now()
	_, err = db.ExecContext(ctx, insertStatement, now.UnixNano())

	if err != nil {
		logrus.WithError(err).Error("Unable to perform INSERT operation")
	} else {
		logrus.Info("Successfully INSERTED into our table")
	}

	// Test SELECT permissions
	selectStatement := `SELECT COUNT(*), MAX(entry_time) FROM example limit 3;`
	row := db.QueryRowContext(ctx, selectStatement)

	var count int
	var maxTimestamp *int64

	err = row.Scan(&count, &maxTimestamp)
	if err != nil {
		logrus.WithError(err).Error("Unable to perform SELECT operation")
	} else {
		mostRecentDate := "unknown"
		if maxTimestamp != nil {
			mostRecentDate = time.Unix(0, *maxTimestamp).Format(time.RFC3339)
		}

		logrus.Info("Successfully SELECTED, most recent value: ", mostRecentDate)
	}

	return nil
}

func buildConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s",
		os.Getenv(DatabaseUser),
		url.QueryEscape(os.Getenv(DatabasePassword)),
		os.Getenv(DatabaseHost),
		os.Getenv(DatabaseName))
}
