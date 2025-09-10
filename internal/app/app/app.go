package appGrpc

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/k3mpton/todoList/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Appp struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int, service server.Service) *Appp {
	grpcS := grpc.NewServer()

	reflection.Register(grpcS)

	server.NewServer(grpcS, service)

	return &Appp{
		log:        log,
		grpcServer: grpcS,
		port:       port,
	}
}

func (a *Appp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *Appp) Run() error {
	const op = "appGrpc.Run"

	conn, err := net.Listen("tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	if err := a.grpcServer.Serve(conn); err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}

func (a *Appp) StopApp() {

	const op = "grpcApp.StopApp"

	a.log.With(
		"op", op,
	).Info("stopped grpc server...", slog.Int("port", a.port))

	a.grpcServer.GracefulStop()
}
