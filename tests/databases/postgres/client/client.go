package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

const (
	DatabaseHost     = "DATABASE_HOST"
	DatabasePort     = "DATABASE_PORT"
	DatabaseNAME     = "DATABASE_NAME"
	DatabaseUser     = "DATABASE_USER"
	DatabasePassword = "DATABASE_PASSWORD"
)

func main() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelFunc()

	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			go func() {
				err := checkDatabaseAccess()
				if err != nil {
					logrus.WithError(err).Error("Database query loop failed")
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}

func checkDatabaseAccess() error {
	psqlInfo := getConnectionInfo()

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	// Test INSERT permissions
	insertStatement := `INSERT INTO example (upload_time) VALUES ($1)`
	now := time.Now()
	_, err = db.Exec(insertStatement, now.UnixNano())
	if err != nil {
		fmt.Println("Unable to perform INSERT operation")
	} else {
		fmt.Println("Successfully INSERTED into our table")
	}

	// Test SELECT permissions
	selectStatement := `SELECT COUNT(*), MAX(upload_time) FROM example limit 3;`
	row := db.QueryRow(selectStatement)

	var count int
	var maxTimestamp int64

	err = row.Scan(&count, &maxTimestamp)
	if err != nil {
		fmt.Println("Unable to perform SELECT operation")
	} else {
		mostRecentDate := time.Unix(0, maxTimestamp).Format(time.RFC3339)
		fmt.Println("Successfully SELECTED, most recent value: ", mostRecentDate)
	}

	return nil
}

func getConnectionInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s, password=%s dbname=%s sslmode=disable",
		os.Getenv(DatabaseHost), os.Getenv(DatabasePort), os.Getenv(DatabaseUser), os.Getenv(DatabasePassword), os.Getenv(DatabaseNAME))
}
