package main

import (
	"context"
	"errors"
	"lo/internal/handler"
	"lo/internal/repo"
	"lo/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	repo := repo.NewTaskStorage()

	logChan := make(chan string, 100)
	go service.LogWorker(logChan)

	srv := service.NewTaskService(repo, logChan)

	h := handler.NewTaskHandler(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.HandleTasks)
	mux.HandleFunc("/tasks/", h.HandleTaskByID)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Printf("Сервер запущен на :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Инициируется завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Принудительное завершение сервера: %v", err)
	}

	tasks := srv.GetAll("")
	log.Println("Состояние задач на момент завершения:")
	for _, task := range tasks {
		log.Printf("Task ID: %d, Title: %s, Status: %s, Retries: %d", task.Id, task.Title, task.Status, task.Retries)
	}

	close(logChan)
	time.Sleep(1 * time.Second)
	log.Println("Сервер успешно завершен")
}
