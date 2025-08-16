package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"task-manager/internal/models"
	service "task-manager/internal/services"
)

type TaskService interface {
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	GetTask(ctx context.Context, id int) (models.Task, error)
	GetTasks(ctx context.Context) ([]models.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

type Handlers struct {
	taskSvc TaskService
}

func NewHandlers(services TaskService) *Handlers {
	return &Handlers{
		taskSvc: services,
	}
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := task.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := h.taskSvc.CreateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.taskSvc.GetTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *Handlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskSvc.GetTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.taskSvc.DeleteTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task deleted successfully"))
}
