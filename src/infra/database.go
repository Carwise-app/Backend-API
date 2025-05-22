package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func ConnectDb() *sql.DB {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		strings.TrimSpace(os.Getenv("DB_HOST")),
		strings.TrimSpace(os.Getenv("DB_PORT")),
		strings.TrimSpace(os.Getenv("DB_USER")),
		strings.TrimSpace(os.Getenv("DB_PASSWD")),
		strings.TrimSpace(os.Getenv("DB_NAME")),
	)
	db, err := sql.Open("postgres", connectionUrl)
	if err != nil {
		log.Fatal("failed to connect to database: %w", err)
	}

	return db
}
