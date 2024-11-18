package main

type MockRepository struct {
}

func (m *MockRepository) CreateUser() error {
	return nil
}

func (m *MockRepository) CreateTask(t *Task) (*Task, error) {
	return &Task{}, nil
}

func (m *MockRepository) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}
