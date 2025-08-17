package repo

import (
	"lo/internal/model"
	"sync"
)

type TaskStorageInterface interface {
	GetAll(status string) []model.Task
	GetById(id int) model.Task
	Update(task model.Task) model.Task
	Create(task model.Task) model.Task
}

type TaskStorage struct {
	mu     sync.RWMutex
	task   map[int]model.Task
	nextID int
}

func NewTaskStorage() TaskStorageInterface {
	return &TaskStorage{
		task:   make(map[int]model.Task),
		nextID: 1,
	}
}

func (r *TaskStorage) Create(task model.Task) model.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.Id = r.nextID
	r.nextID++
	r.task[task.Id] = task
	return task
}

func (r *TaskStorage) Update(task model.Task) model.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.task[task.Id] = task
	return task
}

func (r *TaskStorage) GetById(id int) model.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task := r.task[id]
	return task
}

func (r *TaskStorage) GetAll(status string) []model.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tasks []model.Task
	for _, task := range r.task {
		if status == "" || task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
