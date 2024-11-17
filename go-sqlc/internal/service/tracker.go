package service

import (
	"sync"
	"sync/atomic"
)

type Tracker interface {
	TrackTask() func()
	WaitForTasks()
	TaskCount() int32
}

type TaskTracker struct {
	wg    sync.WaitGroup
	tasks int32
}

func NewTaskTracker() *TaskTracker {
	return &TaskTracker{}
}

func (t *TaskTracker) TrackTask() func() {
	t.wg.Add(1)
	atomic.AddInt32(&t.tasks, 1)
	return func() {
		atomic.AddInt32(&t.tasks, -1)
		t.wg.Done()
	}
}

func (t *TaskTracker) WaitForTasks() {
	t.wg.Wait()
}

func (t *TaskTracker) TaskCount() int32 {
	return atomic.LoadInt32(&t.tasks)
}
