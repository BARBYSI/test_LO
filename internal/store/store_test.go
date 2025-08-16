package store

import (
	"fmt"
	"sync"
	"testing"

	"task-manager/internal/models"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store.tasks == nil {
		t.Errorf("expected store.tasks to be initialized")
	}
}

func TestGet(t *testing.T) {
	store := NewStore()
	task := models.Task{ID: 1, Title: "Task 1"}
	store.Set(1, task)

	got, ok := store.Get(1)
	if !ok {
		t.Errorf("expected task to be found")
	}
	if got.ID != task.ID || got.Title != task.Title {
		t.Errorf("expected task to be %v, got %v", task, got)
	}
}

func TestGetNotFound(t *testing.T) {
	store := NewStore()
	_, ok := store.Get(1)
	if ok {
		t.Errorf("expected task not to be found")
	}
}

func TestSet(t *testing.T) {
	store := NewStore()
	task := models.Task{ID: 1, Title: "Task 1"}
	store.Set(1, task)

	got, ok := store.Get(1)
	if !ok {
		t.Errorf("expected task to be found")
	}
	if got.ID != task.ID || got.Title != task.Title {
		t.Errorf("expected task to be %v, got %v", task, got)
	}
}

func TestDelete(t *testing.T) {
	store := NewStore()
	task := models.Task{ID: 1, Title: "Task 1"}
	store.Set(1, task)

	store.Delete(1)

	_, ok := store.Get(1)
	if ok {
		t.Errorf("expected task to be deleted")
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := NewStore()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			task := models.Task{ID: i, Title: fmt.Sprintf("Task %d", i)}
			store.Set(i, task)
		}(i)
	}

	wg.Wait()

	for i := 0; i < 10; i++ {
		got, ok := store.Get(i)
		if !ok {
			t.Errorf("expected task to be found")
		}
		if got.ID != i || got.Title != fmt.Sprintf("Task %d", i) {
			t.Errorf("expected task to be %v, got %v", i, got)
		}
	}
}
