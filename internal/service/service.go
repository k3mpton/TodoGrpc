package service

import (
	"context"
	"fmt"
	"log/slog"

	todopb "github.com/k3mpton/todoList/protoc/gen/go"
)

type TaskProvider interface {
	GetTask(
		ctx context.Context,
		task_id int64,
		user_id int64,
	) (*todopb.Task, error)

	TasksList(
		ctx context.Context,
		user_id int64,
	) ([]*todopb.Task, error)
}

type TaskChanges interface {
	DeleteTask(
		ctx context.Context,
		task_id int64,
		user_id int64,
	) error

	UpdateTask(
		ctx context.Context,
		NewTitle string,
		NewDescript string,
		task_id int64,
		user_id int64,
	) (string, string, error)

	MarkTaskAsDone(
		ctx context.Context,
		task_id int64,
		user_id int64,
	) error
}

type TaskSaver interface {
	CreateTask(
		ctx context.Context,
		taskModel *todopb.Task,
	) error
}

type service struct {
	log      *slog.Logger
	saver    TaskSaver
	provider TaskProvider
	changer  TaskChanges
}

func NewSerice(
	log *slog.Logger,
	t TaskSaver,
	t1 TaskProvider,
	t2 TaskChanges,
) service {
	return service{
		log:      log,
		saver:    t,
		provider: t1,
		changer:  t2,
	}
}

func (s *service) CreateTask(
	ctx context.Context,
	task *todopb.Task,
) error {
	const op = "service.CreateTask"

	log := s.log.With(
		"op", op,
		"user_id", task.UserId,
	)

	log.Info("Creat task.....")

	err := s.saver.CreateTask(
		ctx,
		task,
	)

	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Append New Task!")

	return nil
}

func (s *service) GetTask(
	ctx context.Context,
	task_id int64,
	user_id int64,
) (*todopb.Task, error) {
	const op = "service.GetTask"

	log := s.log.With(
		"op", op,
		"user_id", user_id,
	)

	log.Info("getting task...")

	task, err := s.provider.GetTask(
		ctx,
		task_id,
		user_id,
	)

	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Access get yor task")

	return task, nil
}

func (s *service) ListTask(
	ctx context.Context,
	user_id int64,
) ([]*todopb.Task, error) {
	const op = "service.ListTask"

	log := s.log.With(
		"op", op,
		"user_id", user_id,
	)

	log.Info("getting user tasks....")

	tasksUser, err := s.provider.TasksList(
		ctx,
		user_id,
	)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Completed getting tasks!")

	return tasksUser, nil
}

func (s *service) UpdateTask(
	ctx context.Context,
	NewDescript string,
	NewTitle string,
	task_id int64,
	user_id int64,
) (string, string, error) {
	const op = "service.UpdateTask"

	log := s.log.With(
		"op", op,
	)

	log.Info("Update Description Task...")

	newDesc, newTitle, err := s.changer.UpdateTask(ctx, NewTitle, NewDescript, task_id, user_id)

	if err != nil {
		return "", "", fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Completed updated task!")

	return newDesc, newTitle, nil
}

func (s *service) DeleteTask(
	ctx context.Context,
	task_id int64,
	user_id int64,
) error {
	const op = "service.DeleteTask"

	log := s.log.With(
		"op", op,
		"user_id", user_id,
	)

	log.Info("Start del task...")

	if err := s.changer.DeleteTask(ctx, task_id, user_id); err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	log.Info("Completed deleted task!")

	return nil
}

func (s *service) MarkTaskAsDone(
	ctx context.Context,
	task_i int64,
	user_id int64,
) error {
	const op = "service.MarkTaskAsDone"

	log := s.log.With(
		"op", op,
	)

	log.Info("the task status has been changed to completed")

	err := s.changer.MarkTaskAsDone(ctx, task_i, user_id)

	if err != nil {
		return fmt.Errorf("%v:%v", op, err)
	}

	log.Info("Ok Task As done ")

	return nil
}
