package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"task-manager/internal/models"
	"task-manager/internal/store"

)

func TestCreateTask(t *testing.T) {
	st := store.NewStore()
	repo := NewRepository(st)

	ctx := context.Background()
	task := models.Task{
		Title:       "Test Task",
		Description: "Desc",
	}

	createdTask, err := repo.CreateTask(ctx, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if createdTask.ID == 0 {
		t.Error("expected non-zero task ID")
	}
	if createdTask.Title != task.Title {
		t.Errorf("expected title %q, got %q", task.Title, createdTask.Title)
	}
}

func TestGetTask(t *testing.T) {
	st := store.NewStore()
	repo := NewRepository(st)

	ctx := context.Background()

	task, _ := repo.CreateTask(ctx, models.Task{Title: "Test Task"})
	retrievedTask, err := repo.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if retrievedTask.ID != task.ID {
		t.Errorf("expected ID %d, got %d", task.ID, retrievedTask.ID)
	}
	if retrievedTask.Title != task.Title {
		t.Errorf("expected title %q, got %q", task.Title, retrievedTask.Title)
	}
}

func TestGetTasks(t *testing.T) {
	st := store.NewStore()
	repo := NewRepository(st)

	ctx := context.Background()

	repo.CreateTask(ctx, models.Task{Title: "Task1"})
	repo.CreateTask(ctx, models.Task{Title: "Task2"})

	tasks, err := repo.GetTasks(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestDeleteTask(t *testing.T) {
	st := store.NewStore()
	repo := NewRepository(st)

	ctx := context.Background()

	task, _ := repo.CreateTask(ctx, models.Task{Title: "To delete"})

	err := repo.DeleteTask(ctx, task.ID)
	if err != nil {
		t.Errorf("unexpected error on delete: %v", err)
	}

	tasks, _ := repo.GetTasks(ctx)
	for _, tsk := range tasks {
		if tsk.ID == task.ID {
			t.Error("deleted task still present in store")
		}
	}
}

func TestCreateTask_ContextCanceled(t *testing.T) {
	st := store.NewStore()
	repo := NewRepository(st)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	task := models.Task{Title: "Invalid due to ctx cancel"}
	_, err := repo.CreateTask(ctx, task)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestGetTasks_ContextTimeout(t *testing.T) {
	st := store.NewStore()
	repo := NewRepository(st)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond)

	_, err := repo.GetTasks(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded error, got %v", err)
	}
}

func TestDeleteTask_ContextCanceled(t *testing.T) {
	st := store.NewStore()
	repo := NewRepository(st)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.DeleteTask(ctx, 123)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}
