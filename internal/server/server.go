package server

import (
	"context"
	"fmt"
	"log/slog"

	todopb "github.com/k3mpton/todoList/protoc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	CreateTask(
		ctx context.Context,
		task *todopb.Task,
	) error

	GetTask(
		ctx context.Context,
		task_id int64,
		user_id int64,
	) (*todopb.Task, error)

	ListTask(
		ctx context.Context,
		user_id int64,
	) ([]*todopb.Task, error)

	UpdateTask(
		ctx context.Context,
		NewDescript string,
		NewTitle string,
		task_id int64,
		user_id int64,
	) (string, string, error)

	DeleteTask(
		ctx context.Context,
		task_id int64,
		user_id int64,
	) error

	MarkTaskAsDone(
		ctx context.Context,
		task_id int64,
		user_id int64,
	) error
}

type ApiServer struct {
	todopb.UnimplementedTaskServiceServer
	service Service
}

func NewServer(serverGrpc *grpc.Server, service Service) {
	todopb.RegisterTaskServiceServer(
		serverGrpc,
		&ApiServer{
			service: service,
		},
	)
}

// var OldDescriptionAndTitle map[int]string // id and description

func (s *ApiServer) CreateTask(ctx context.Context, r *todopb.CreateTaskRequest) (*todopb.CreateTaskResponse, error) {
	const op = "server.CreateTask"

	log := slog.With(
		"op", op,
		"title", r.Task.Title,
	)

	log.Info("Created task.....")

	task := todopb.Task{
		Id:          r.Task.Id,
		Title:       r.Task.Title,
		Description: r.Task.Description,
		Status:      r.Task.Status,
		DueDate:     r.Task.DueDate,
		UserId:      r.Task.UserId,
		CreatedAt:   r.Task.CreatedAt,
		UpdatedAt:   r.Task.UpdatedAt,
	}

	fmt.Println("sdjfkasdjfkjaskdjf")

	if err := s.service.CreateTask(ctx, &task); err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	return &todopb.CreateTaskResponse{
		Task: r.Task,
	}, nil
}

func (s *ApiServer) GetTask(ctx context.Context, r *todopb.GetTaskRequest) (*todopb.GetTaskResponse, error) {
	const op = "server.GetTask"

	log := slog.With(
		"op", op,
	)

	log.Info("Getting task...")

	task, err := s.service.GetTask(ctx, r.Id, r.UserId)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	return &todopb.GetTaskResponse{
		Task: task,
	}, nil
}

func (s *ApiServer) ListTasks(ctx context.Context, r *todopb.ListTasksRequest) (*todopb.ListTasksResponse, error) {
	const op = "server.ListTasks"

	log := slog.With(
		"op", op,
		"userId", r.UserId,
	)

	log.Info("List tasks..d.")

	tasks, err := s.service.ListTask(ctx, r.UserId)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Congratulations! Your list tasks access")

	return &todopb.ListTasksResponse{
		Tasks: tasks,
	}, nil
}

func (s *ApiServer) UpdateTask(ctx context.Context, r *todopb.UpdateTaskRequest) (*todopb.UpdateTaskResponse, error) {
	const op = "server.UpdateTask"

	log := slog.With(
		"op", op,
		"taskid", r.Task.Id,
		"userId", r.Task.UserId,
	)

	log.Info("Updating test...")

	newTitle, newDescript, err := s.service.UpdateTask(ctx, r.Task.Description, r.Task.Title, r.Task.Id, r.Task.UserId)

	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Access, test is Updating! ")

	return &todopb.UpdateTaskResponse{
		Task: &todopb.Task{
			Title:       newTitle,
			Description: newDescript,
		},
	}, nil
}

func (s *ApiServer) DeleteTask(ctx context.Context, r *todopb.DeleteTaskRequest) (*emptypb.Empty, error) {
	const op = "server.DeleteTask"

	log := slog.With(
		"op", op,
		"task_id", r.Id,
	)

	log.Info("Processing delete task...")

	if err := s.service.DeleteTask(ctx, r.Id, r.UserId); err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Access! Delete task!")

	return nil, nil
}

func (s *ApiServer) MarkTaskAsDone(ctx context.Context, r *todopb.MarkTaskAsDoneRequest) (*todopb.MarkTaskAsDoneResponse, error) {
	const op = "server.MarkTaskAsDone"

	log := slog.With(
		"op", op,
		"task_id", r.Id,
	)

	log.Info("Swapping task false -> true")

	if err := s.service.MarkTaskAsDone(ctx, r.Id, r.UserId); err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Acess Swap fasle -> true task!")

	return &todopb.MarkTaskAsDoneResponse{}, nil
}
