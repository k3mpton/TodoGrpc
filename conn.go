package connectiondb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func Connection() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("fail read .env: %v", err)
	}

	dbUrl := os.Getenv("DbSourse")

	db, err := sql.Open("pgx", dbUrl)

	if err != nil {
		log.Fatalf("fail open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("fail ping db: %v", err)
	}

	return db
}
