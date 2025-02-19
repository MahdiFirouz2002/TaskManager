package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nikandishan/structs/task"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var mockTask = task.Task{
	ID:        uuid.New().String(),
	Title:     "Test Task",
	Status:    "pending",
	CreatedAt: time.Now(),
}

func mockGetTasks(ctx context.Context, status string) ([]task.Task, error) {
	return []task.Task{mockTask}, nil
}

func mockCreateTask(newTask task.Task) (error, task.Task) {
	newTask.ID = uuid.New().String()
	newTask.CreatedAt = time.Now()
	return nil, newTask
}

func TestGetTasks(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	r := gin.Default()
	r.GET("/tasks", func(c *gin.Context) {
		getTasks(c, mockGetTasks)
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []task.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.Nil(t, err)
	assert.Greater(t, len(tasks), 0)
	assert.Equal(t, "Test Task", tasks[0].Title)
}

func TestCreateTask(t *testing.T) {
	taskData := `{"title": "New Task", "status": "in_progress"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(taskData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.Default()
	r.POST("/tasks", func(c *gin.Context) {
		createTask(c, mockCreateTask)
	})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTask task.Task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.Nil(t, err)
	assert.Equal(t, "New Task", createdTask.Title)
	assert.Equal(t, "in_progress", createdTask.Status)
}
