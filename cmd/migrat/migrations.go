package main

import (
	"flag"
	"log"

	connectiondb "github.com/k3mpton/todoList"
	"github.com/pressly/goose"
)

var (
	MigMove = flag.String("m", "up", "migration move, up or down")
)

func main() {
	flag.Parse()
	db := connectiondb.Connection()
	defer db.Close()
	if *MigMove != "up" && *MigMove != "down" {
		log.Fatalln("не удалось получить движение миграции")
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	dirMigrations := "./migrations"
	switch *MigMove {
	case "up":
		if err := goose.Up(db, dirMigrations); err != nil {
			log.Fatalf("fail migration up: %v", err)
		}
	default:
		if err := goose.Down(db, dirMigrations); err != nil {
			log.Fatalf("fail migration down: %v", err)
		}
	}
}
