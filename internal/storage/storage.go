package storage

import (
	"context"
	"database/sql"
	"fmt"

	connectiondb "github.com/k3mpton/todoList"
	"github.com/k3mpton/todoList/internal/models"
	todopb "github.com/k3mpton/todoList/protoc/gen/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Sql struct {
	db *sql.DB
}

func NewStorage() *Sql {
	db := connectiondb.Connection()
	return &Sql{
		db: db,
	}
}

func (s *Sql) GetTask(
	ctx context.Context,
	task_id int64,
	user_id int64,
) (*todopb.Task, error) {
	const op = "storage.GetTask"

	query := `select * from tasks where id = $1 and user_id = $2`
	row := s.db.QueryRowContext(ctx, query, task_id, user_id)
	if row.Err() != nil {
		return nil, fmt.Errorf("%v: %v", op, row.Err())
	}

	var task models.Task
	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DueDate,
		&task.UserId,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}

	return &todopb.Task{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     timestamppb.New(task.DueDate),
		UserId:      task.UserId,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}, nil
}

func (s *Sql) TasksList(
	ctx context.Context,
	user_id int64,
) ([]*todopb.Task, error) {
	const op = "storage.TasksList"

	query := `select tasks.* from tasks where user_id = $1 
	order by created_at desc`

	rows, err := s.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", op, err)
	}
	defer rows.Close()

	// var tasksForUser []*todopb.Task
	// for rows.Next() {
	// 	var task todopb.Task
	// 	err := rows.Scan(
	// 		&task.Id,
	// 		&task.Title,
	// 		&task.Description,
	// 		&task.Status,
	// 		&task.DueDate,
	// 		&task.UserId,
	// 		&task.CreatedAt,
	// 		&task.UpdatedAt,
	// 	)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%v: %v", op, err)
	// 	}
	// 	tasksForUser = append(tasksForUser, &task)
	// }

	var tasksForUser []*todopb.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.DueDate,
			&task.UserId,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", op, err)
		}

		taskTodopd := todopb.Task{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			DueDate:     timestamppb.New(task.DueDate),
			UserId:      task.UserId,
			CreatedAt:   timestamppb.New(task.CreatedAt),
			UpdatedAt:   timestamppb.New(task.UpdatedAt),
		}

		tasksForUser = append(tasksForUser, &taskTodopd)
	}
	return tasksForUser, nil
}

func (s *Sql) DeleteTask(
	ctx context.Context,
	task_id int64,
	user_id int64,
) error {
	const op = "storage.DeleteTask"

	query := `delete from Tasks where id = $1 and user_id = $2`
	r, err := s.db.Exec(query, task_id, user_id)

	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	rows, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	if rows == 0 {
		return fmt.Errorf("%v: failed search your task((", op)
	}

	return nil
}

func (s *Sql) UpdateTask(
	ctx context.Context,
	NewTitle string,
	NewDescript string,
	task_id int64,
	user_id int64,
) (string, string, error) {
	const op = "storage.UpdateTask"

	query := `update Tasks
	set title = $1, description = $2,
	updated_at = Now() where id = $3 and user_id = $4
	 returning title, description`

	var updatedTitle, UpdatedDescription string
	err := s.db.QueryRowContext(ctx, query, NewTitle,
		NewDescript, task_id, user_id).Scan(&updatedTitle, &UpdatedDescription)
	if err != nil {
		return "", "", fmt.Errorf("%v: %v", op, err)
	}

	return updatedTitle, UpdatedDescription, nil
}

func (s *Sql) CreateTask(
	ctx context.Context,
	taskModel *todopb.Task,
) error {
	const op = "storage.CreateTask"

	query := `insert into 
	tasks(title, description, status, due_date, user_id)
	values ($1, $2, $3, $4, $5)`

	_, err := s.db.ExecContext(ctx, query,
		taskModel.Title,
		taskModel.Description,
		taskModel.Status,
		taskModel.DueDate.AsTime(),
		taskModel.UserId,
	)

	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	return nil
}

func (s *Sql) MarkTaskAsDone(
	ctx context.Context,
	task_id int64,
	user_id int64,
) error {
	const op = "storage.MarkTaskAsDone"

	query := `
		update Tasks 
		set status = true
		where user_id = $1 and id = $2
	`

	r, err := s.db.ExecContext(ctx, query, user_id, task_id)
	if err != nil {
		return fmt.Errorf("%v: %v", op, err)
	}

	rows, _ := r.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("%v: %v", op, "could not find user or task")
	}
	return nil
}
