package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	concurrentlogger "nikandishan/concurrentLogger"
	"nikandishan/structs/task"
	"nikandishan/utils/customeError"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	r.GET("/tasks", func(c *gin.Context) {
		getTasks(c, task.GetTasks)
	})
	r.GET("/task/:id", getTask)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())

	protected.POST("/tasks", func(c *gin.Context) {
		createTask(c, task.CreateTask)
	})
	protected.PUT("/tasks/:id", updateTask)
	protected.DELETE("/tasks/:id", deleteTask)
}

func getTask(c *gin.Context) {
	id := c.Param("id")

	task, err := task.GetTask(id)
	if errors.Is(err, customeError.ErrTaskNotFound) {
		RespondWithError(c, http.StatusNotFound, "Task not found", "TASK_404")
		return
	} else if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Internal server error", "TASK_500")
		return
	}

	c.JSON(http.StatusOK, task)
}

func getTasks(c *gin.Context, getTasksFunc func(ctx context.Context, status string) ([]task.Task, error)) {
	status := c.Query("status")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	tasks, err := getTasksFunc(ctx, status)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			RespondWithError(c, http.StatusRequestTimeout, "Request timed out", "TASK_408")
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Internal server error", "TASK_500")
		}
		return
	}

	if len(tasks) == 0 {
		RespondWithError(c, http.StatusOK, "No tasks found", "TASK_204")
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context, createTaskFunc func(task.Task) (error, task.Task)) {
	var newTask task.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request body", "TASK_400")
		return
	}

	err, createdTask := createTaskFunc(newTask)
	if errors.Is(err, customeError.ErrInvalidTaskFormat) {
		RespondWithError(c, http.StatusBadRequest, err.Error(), "TASK_400")
		return
	} else if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to create task", "TASK_500")
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")

	var inputTask task.Task
	if err := c.ShouldBindJSON(&inputTask); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request body", "TASK_400")
		return
	}

	err, updatedTask := task.UpdateTask(id, inputTask)
	if errors.Is(err, customeError.ErrInvalidTaskFormat) {
		RespondWithError(c, http.StatusBadRequest, "Invalid task format", "TASK_400")
		return
	} else if errors.Is(err, customeError.ErrTaskNotFound) {
		RespondWithError(c, http.StatusNotFound, "Task not found", "TASK_404")
		return
	} else if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to update task", "TASK_500")
		return
	}

	select {
	case concurrentlogger.UpdateTaskChan <- updatedTask:
	default:
		log.Println("Warning: Task update dropped due to full queue")
	}
	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	err := task.DeleteTask(id)
	if errors.Is(err, customeError.ErrTaskNotFound) {
		RespondWithError(c, http.StatusNotFound, "Task not found", "TASK_404")
		return
	} else if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to delete task", "TASK_500")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
