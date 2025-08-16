package services

import (
	"context"
	"errors"

	"task-manager/internal/models"
	"task-manager/internal/repository"

)

var (
	ErrTaskNotFound = repository.ErrTaskNotFound
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	GetTask(ctx context.Context, id int) (models.Task, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

type TaskService struct {
	rep TaskRepository
}

func NewTaskService(repository TaskRepository) *TaskService {
	return &TaskService{
		rep: repository,
	}
}

func (t *TaskService) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	createdTask, err := t.rep.CreateTask(ctx, task)
	if err != nil {
		return models.Task{}, err
	}
	return createdTask, nil
}

func (t *TaskService) GetTask(ctx context.Context, id int) (models.Task, error) {
	task, err := t.rep.GetTask(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}

	return task, nil
}

func (t *TaskService) GetTasks(ctx context.Context) ([]models.Task, error) {
	return t.rep.GetTasks(ctx)
}

func (t *TaskService) DeleteTask(ctx context.Context, id int) error {

	err := t.rep.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		return err
	}

	return nil
}
