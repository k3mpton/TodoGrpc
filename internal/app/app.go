package app

import (
	"log/slog"

	appGrpc "github.com/k3mpton/todoList/internal/app/app"
	"github.com/k3mpton/todoList/internal/service"
	"github.com/k3mpton/todoList/internal/storage"
)

type App struct {
	Appp *appGrpc.Appp
}

func NewApp(
	log *slog.Logger,
	port int,
) *App {

	storage := storage.NewStorage()

	sarvice := service.NewSerice(log, storage, storage, storage)

	app := appGrpc.NewApp(log, port, &sarvice)

	return &App{
		Appp: app,
	}
}
