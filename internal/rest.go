package rest

import (
	"context"
	"fmt"
	"net/http"

	"task-manager/internal/config"
	"task-manager/internal/handlers"
	"task-manager/internal/repository"
	svc "task-manager/internal/services"
	"task-manager/internal/store"
	"task-manager/pkg/logger"

)

type Rest struct {
	config   *config.Config
	handlers *handlers.Handlers
	srv      *http.Server
	store    *store.Store
	service  *svc.TaskService
	router   *http.ServeMux
}

func NewRest(cfg *config.Config) *Rest {
	store := store.NewStore()
	repository := repository.NewRepository(store)
	taskService := svc.NewTaskService(repository)
	taskHandlers := handlers.NewHandlers(taskService)

	router := initRouter(taskHandlers)

	rest := &Rest{
		config:   cfg,
		handlers: taskHandlers,
		router:   router,
		srv: &http.Server{
			Addr:    cfg.AppPort,
			Handler: router,
		},
		store:   store,
		service: taskService,
	}
	return rest
}

func initRouter(h *handlers.Handlers) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				h.GetTask(w, r)
			} else {
				h.GetTasks(w, r)
			}
		case http.MethodPost:
			h.CreateTask(w, r)
		case http.MethodDelete:
			h.DeleteTask(w, r)
		default:
			w.Header().Set("Allow", "GET, POST, DELETE")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test task for LO"))
	})

	return router
}

func (r *Rest) RunRest() error {
	logger.LogInfo("starting server on port " + r.config.AppPort)
	err := r.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.LogError(fmt.Sprintf("failed to start server: %v", err))
		return err
	}

	return nil
}

func (r *Rest) ShutdownRest(ctx context.Context) error {
	logger.LogInfo("shutting down server")
	return r.srv.Shutdown(ctx)
}
