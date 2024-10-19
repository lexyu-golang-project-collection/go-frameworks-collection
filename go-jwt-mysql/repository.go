package main

import (
	"database/sql"
	"log"
)

type Repository interface {
	CreateUser() error

	CreateTask(t *Task) (*Task, error)

	GetTask(id string) (*Task, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser() error {
	return nil
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name, status, projectId, assignedToID) VALUES(?, ?, ?, ?)",
		t.Name, t.Status, t.ProjectID, t.AssignedToID)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	t.ID = id

	return t, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {

	var t Task

	stmt, err := s.db.Prepare("SELECT id, name, status, project_id, assigned_to, createAt FROM tasks WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID, &t.CreatedAt)

	return &t, err
}
