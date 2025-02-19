package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	loadTestCount = 100
)

func TestLoadCreateAndGetTasks(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/tasks", func(c *gin.Context) {
		createTask(c, mockCreateTask)
	})
	r.GET("/tasks", func(c *gin.Context) {
		getTasks(c, mockGetTasks)
	})

	var wg sync.WaitGroup

	for i := 0; i < loadTestCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			taskData := fmt.Sprintf(`{"title": "Task %d", "status": "pending"}`, i)
			req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(taskData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
		}(i)
	}

	for i := 0; i < loadTestCount/10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req, _ := http.NewRequest("GET", "/tasks", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
		}()
	}

	wg.Wait()
}
