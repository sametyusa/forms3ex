package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

const maxAttempts = 5
const dbName = "innsecure"
const connectionFormat = "host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable"

func NewConnection(host, user, password string) (*sql.DB, error) {
	psqlconn := fmt.Sprintf(connectionFormat, host, user, password, dbName)
	fmt.Println(psqlconn)
	var err error
	for i := 0; i < maxAttempts; i++ {
		db, err := sql.Open("postgres", psqlconn)
		if err == nil {
			return db, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		fmt.Printf("failed to connect to database with string: %s\n", psqlconn)
	}
	return nil, err
}
