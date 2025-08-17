package model

type Task struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"` // Pending, Completed, Failed, Error
	Retries int    `json:"retries"`
}
