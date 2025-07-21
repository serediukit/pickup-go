package database

import (
	"database/sql"
	"fmt"
	"pickup-srv/util"

	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
	host := util.GetEnv("DB_HOST", "localhost")
	port := util.GetEnv("DB_PORT", "5432")
	user := util.GetEnv("DB_USER", "postgres")
	password := util.GetEnv("DB_PASSWORD", "password")
	dbname := util.GetEnv("DB_NAME", "pickup_db")
	sslmode := util.GetEnv("DB_SSLMODE", "disable")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
