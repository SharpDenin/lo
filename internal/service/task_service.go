package service

import (
	"lo/internal/repo"
)

type TaskService struct {
	repo repo.TaskStorageInterface
	log  chan string
}

func NewTaskService(repo repo.TaskStorageInterface, log chan string) *TaskService {
	return &TaskService{
		repo: repo,
		log:  log,
	}
}
