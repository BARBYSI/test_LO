package services

import (
	"context"
	"errors"
	"testing"

	"task-manager/internal/models"
	"task-manager/internal/repository"

)

type MockTaskRepository struct {
	CreateTaskFunc func(ctx context.Context, task models.Task) (models.Task, error)
	GetTaskFunc    func(ctx context.Context, id int) (models.Task, error)
	GetTasksFunc   func(ctx context.Context) ([]models.Task, error)
	DeleteTaskFunc func(ctx context.Context, id int) error
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	return m.CreateTaskFunc(ctx, task)
}
func (m *MockTaskRepository) GetTask(ctx context.Context, id int) (models.Task, error) {
	return m.GetTaskFunc(ctx, id)
}
func (m *MockTaskRepository) GetTasks(ctx context.Context) ([]models.Task, error) {
	return m.GetTasksFunc(ctx)
}
func (m *MockTaskRepository) DeleteTask(ctx context.Context, id int) error {
	return m.DeleteTaskFunc(ctx, id)
}

func TestCreateTask_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		CreateTaskFunc: func(ctx context.Context, task models.Task) (models.Task, error) {
			task.ID = 1
			return task, nil
		},
	}
	service := NewTaskService(mockRepo)

	task := models.Task{Title: "Test"}
	created, err := service.CreateTask(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}
}

func TestGetTasks_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		GetTasksFunc: func(ctx context.Context) ([]models.Task, error) {
			return []models.Task{{ID: 1, Title: "Task1"}}, nil
		},
	}
	service := NewTaskService(mockRepo)

	tasks, err := service.GetTasks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestDeleteTask_Success(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteTaskFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}
	service := NewTaskService(mockRepo)

	if err := service.DeleteTask(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	mockRepo := &MockTaskRepository{
		DeleteTaskFunc: func(ctx context.Context, id int) error {
			return repository.ErrTaskNotFound
		},
	}
	service := NewTaskService(mockRepo)

	err := service.DeleteTask(context.Background(), 42)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}
