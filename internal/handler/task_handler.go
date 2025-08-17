package handler

import (
	"encoding/json"
	"lo/internal/model"
	"lo/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.GetTasks(w, r)
	case http.MethodPost:
		h.CreateTask(w, r)
	default:
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var m model.Task
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if m.Title == "" {
		sendError(w, http.StatusBadRequest, "Title is required")
		return
	}
	task, err := h.taskService.Create(m)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	if task.Id == 0 {
		sendError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}
	sendSuccess(w, http.StatusCreated, task, nil)
}

func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(path)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	task := h.taskService.GetById(id)
	if task.Id == 0 {
		sendError(w, http.StatusNotFound, "Task not found")
		return
	}

	sendSuccess(w, http.StatusOK, task, nil)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status != "" && !isValidStatus(status) {
		sendError(w, http.StatusBadRequest, "Invalid status. Use: Pending, Completed, Failed,akaan Error")
		return
	}
	tasks := h.taskService.GetAll(status)
	sendSuccess(w, http.StatusOK, tasks, nil)
}
