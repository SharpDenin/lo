package service

import (
	"errors"
	"fmt"
	"lo/internal/model"
	"lo/internal/repo"
	"math/rand"
	"time"
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

func (s *TaskService) Create(m model.Task) (model.Task, error) {
	if m.Title == "" {
		s.logAction("Failed to create task: title is required")
		return model.Task{}, errors.New("title is required")
	}

	created := s.repo.Create(m)
	if created.Id == 0 {
		s.logAction("Failed to create task: repository error")
		return model.Task{}, errors.New("repository failed to create task")
	}
	created.Status = "Pending"
	created.Retries = 0
	s.repo.Update(created)

	go s.process(created.Id)

	s.logAction(fmt.Sprintf("Task created with status Pending"))
	return created, nil
}

func (s *TaskService) GetById(id int) model.Task {
	return s.repo.GetById(id)
}

func (s *TaskService) GetAll(status string) []model.Task {
	return s.repo.GetAll(status)
}

func (s *TaskService) process(id int) {
	task := s.GetById(id)
	if task.Id == 0 {
		return
	}

	rand.Seed(time.Now().UnixNano())

	for task.Retries <= 2 {
		task.Status = "Pending"
		s.repo.Update(task)
		s.logAction(fmt.Sprintf("Task %d status changed to Pending (attempt %d)", id, task.Retries+1))

		sleepDuration := time.Duration(rand.Intn(61)+60) * time.Second
		time.Sleep(sleepDuration)

		if rand.Float64() < 0.2 {
			task.Status = "Failed"
			s.repo.Update(task)
			s.logAction(fmt.Sprintf("Task %d status changed to Failed", id))

			if task.Retries == 2 {
				task.Status = "Error"
				s.repo.Update(task)
				s.logAction(fmt.Sprintf("Task %d status changed to Error", id))
				break
			}

			task.Retries++
		} else {
			task.Status = "Completed"
			s.repo.Update(task)
			s.logAction(fmt.Sprintf("Task %d status changed to Completed", id))
			break
		}
	}
}

func LogWorker(logChan <-chan string) {
	for logMsg := range logChan {
		fmt.Printf("[LOG] %s\n", logMsg)
	}
}

func (s *TaskService) logAction(msg string) {
	select {
	case s.log <- msg:
	default:
		// Drop if channel full
	}
}
