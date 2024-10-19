package main

import "time"

type Task struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	ProjectID    int64     `json:"project_id"`
	AssignedToID int64     `json:"assigned_to_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateTaskPayload struct {
	Name         string `json:"name"`
	ProjectID    int64  `json:"project_id"`
	AssignedToID int64  `json:"assigned_to_iD"`
}
