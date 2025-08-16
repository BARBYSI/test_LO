package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"task-manager/internal/models"
	"task-manager/internal/services"
)

// Мок для сервиса
type MockTaskService struct {
	CreateTaskFunc func(ctx context.Context, task models.Task) (models.Task, error)
	GetTaskFunc    func(ctx context.Context, id int) (models.Task, error)
	GetTasksFunc   func(ctx context.Context) ([]models.Task, error)
	DeleteTaskFunc func(ctx context.Context, id int) error
}

func (m *MockTaskService) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	return m.CreateTaskFunc(ctx, task)
}

func (m *MockTaskService) GetTask(ctx context.Context, id int) (models.Task, error) {
	return m.GetTaskFunc(ctx, id)
}

func (m *MockTaskService) GetTasks(ctx context.Context) ([]models.Task, error) {
	return m.GetTasksFunc(ctx)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, id int) error {
	return m.DeleteTaskFunc(ctx, id)
}

// Тест CreateTask - успешное создание задачи
func TestCreateTask_Success(t *testing.T) {
	mockSvc := &MockTaskService{
		CreateTaskFunc: func(ctx context.Context, task models.Task) (models.Task, error) {
			task.ID = 1
			return task, nil
		},
	}
	h := NewHandlers(mockSvc)

	reqBody := `{"title":"Test task"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdTask models.Task
	if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
		t.Fatalf("cannot decode response body: %v", err)
	}
	if createdTask.Title != "Test task" {
		t.Errorf("expected name 'Test task', got '%s'", createdTask.Title)
	}
}

func TestGetTasks_Success(t *testing.T) {
	mockTasks := []models.Task{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}
	mockSvc := &MockTaskService{
		GetTasksFunc: func(ctx context.Context) ([]models.Task, error) {
			return mockTasks, nil
		},
	}
	h := NewHandlers(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	h.GetTasks(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var tasks []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Fatalf("cannot decode response body: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestGetTask_Success(t *testing.T) {
	mockTask := models.Task{ID: 1, Title: "Task 1"}
	mockSvc := &MockTaskService{
		GetTaskFunc: func(ctx context.Context, id int) (models.Task, error) {
			return mockTask, nil
		},
	}
	h := NewHandlers(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/tasks?id=1", nil)
	w := httptest.NewRecorder()

	h.GetTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	mockSvc := &MockTaskService{
		DeleteTaskFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}
	h := NewHandlers(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/tasks?id=1", nil)
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	mockSvc := &MockTaskService{
		DeleteTaskFunc: func(ctx context.Context, id int) error {
			return services.ErrTaskNotFound
		},
	}
	h := NewHandlers(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/tasks?id=999", nil)
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestDeleteTask_BadRequest(t *testing.T) {
	mockSvc := &MockTaskService{}
	h := NewHandlers(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/tasks?id=abc", nil)
	w := httptest.NewRecorder()

	h.DeleteTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}
