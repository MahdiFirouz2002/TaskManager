package server

import (
	"context"
	"errors"
	"net/http"
	"nikandishan/structs"
	"nikandishan/structs/task"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {
	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTask)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())

	protected.POST("/tasks", createTask)
	protected.PUT("/tasks/:id", updateTask)
	protected.DELETE("/tasks/:id", deleteTask)
}

func getTask(c *gin.Context) {
	id := c.Param("id")

	task, err := task.GetTask(id)
	if errors.Is(err, structs.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func getTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	tasks, err := task.GetTasks(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timed out"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "There is no Task"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var newTask task.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := task.CreateTask(newTask)
	if errors.Is(err, structs.ErrInvalidTaskFormat) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")

	var updatedTask task.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := task.UpdateTask(id, updatedTask)
	if errors.Is(err, structs.ErrInvalidTaskFormat) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Task Format"})
		return
	} else if errors.Is(err, structs.ErrTaskNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Task ID"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	err := task.DeleteTask(id)
	if errors.Is(err, structs.ErrTaskNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
