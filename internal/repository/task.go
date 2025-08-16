package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"task-manager/internal/models"
	"task-manager/internal/store"
	"task-manager/pkg/logger"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type Repository struct {
	store *store.Store
}

func NewRepository(store *store.Store) *Repository {
	return &Repository{store: store}
}

func (r *Repository) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return models.Task{}, ctx.Err()
	default:
		task.ID = r.store.NextID()
		r.store.Set(task.ID, task)
		logger.LogInfo(fmt.Sprintf("task with title %s created", task.Title))
		return task, nil
	}
}

func (r *Repository) GetTask(ctx context.Context, id int) (models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return models.Task{}, ctx.Err()
	default:
		task, ok := r.store.Get(id)
		if !ok {
			logger.LogInfo(fmt.Sprintf("task with id %d not found", id))
			return models.Task{}, ErrTaskNotFound
		}
		return task, nil
	}
}

func (r *Repository) GetTasks(ctx context.Context) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		tasks := r.store.GetAllTasks()
		return tasks, nil
	}
}

func (r *Repository) DeleteTask(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		deleted := r.store.Delete(id)
		if !deleted {
			logger.LogInfo(fmt.Sprintf("task %d not found", id))
			return ErrTaskNotFound
		}
		logger.LogInfo(fmt.Sprintf("task %d deleted", id))
		return nil
	}
}
