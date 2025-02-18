package task

import (
	"context"
	"nikandishan/structs"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

var (
	tasks = make(map[string]Task)
	mutex sync.RWMutex
	v     = validator.New()
)

type Task struct {
	ID          string    `json:"id" validate:"required,uuid4"`
	Title       string    `json:"title" validate:"required,min=3"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status" validate:"oneof=todo in_progress done"`
	CreatedAt   time.Time `json:"created_at"`
}

func AddTask(title, description, status string) Task {
	mutex.Lock()
	defer mutex.Unlock()

	id := uuid.New().String()

	task := Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   time.Now(),
	}

	tasks[id] = task
	return task
}

func GetTasks(ctx context.Context) ([]Task, error) {
	resultChan := make(chan []Task, 1)

	go func() {
		mutex.RLock()
		defer mutex.RUnlock()

		var taskList []Task
		for _, task := range tasks {
			taskList = append(taskList, task)
		}

		select {
		case resultChan <- taskList:
		case <-ctx.Done():
		}
	}()

	select {
	case taskList := <-resultChan:
		return taskList, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func GetTask(id string) (Task, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	task, exists := tasks[id]
	if !exists {
		return Task{}, structs.ErrTaskNotFound
	}

	return task, nil
}

func CreateTask(newTask Task) error {
	newTask.ID = uuid.New().String()
	newTask.CreatedAt = time.Now()

	if err := v.Struct(newTask); err != nil {
		return structs.ErrInvalidTaskFormat
	}

	mutex.Lock()
	tasks[newTask.ID] = newTask
	mutex.Unlock()

	return nil
}

func UpdateTask(id string, updatedTask Task) error {
	mutex.Lock()
	task, exists := tasks[id]
	mutex.Unlock()

	if !exists {
		return structs.ErrTaskNotFound
	}

	updatedTask.ID = task.ID
	updatedTask.CreatedAt = task.CreatedAt

	if err := v.Struct(updatedTask); err != nil {
		return structs.ErrInvalidTaskFormat
	}

	mutex.Lock()
	tasks[id] = updatedTask
	mutex.Unlock()

	return nil
}

func DeleteTask(id string) error {
	mutex.Lock()
	_, exists := tasks[id]
	if exists {
		delete(tasks, id)
	}
	mutex.Unlock()

	if !exists {
		return structs.ErrTaskNotFound
	}

	return nil
}
