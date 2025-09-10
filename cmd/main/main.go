package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/k3mpton/todoList/internal/app"
	"github.com/k3mpton/todoList/internal/config"
	"github.com/k3mpton/todoList/pkg/logger"
)

func main() {
	cfg := config.MustConfigLoad()
	log := logger.InitLogger(cfg.Env)
	log.Info("renato")
	application := app.NewApp(log, cfg.Grpc.Port)

	go application.Appp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Stopped server", slog.Int("port", cfg.Grpc.Port))

	application.Appp.StopApp()

	log.Info("Server stop!")
}
