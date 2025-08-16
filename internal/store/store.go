package store

import (
	"sync"

	"task-manager/internal/models"
)

type Store struct {
	tasks  map[int]models.Task
	mu     sync.RWMutex
	nextID int
}

func NewStore() *Store {
	return &Store{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *Store) NextID() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	return id
}

func (s *Store) Get(key int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.tasks[key]
	return value, ok
}

func (s *Store) GetAllTasks() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *Store) Set(key int, value models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[key] = value
}

func (s *Store) Delete(key int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[key]
	if !ok {
		return false
	}
	delete(s.tasks, key)
	return true
}
