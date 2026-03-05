package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks = []Task{
	{ID: 1, Title: "Learn Go", Completed: false},
	{ID: 2, Title: "Learn DevOps", Completed: false},
}

var lastID = 2

func setupRouter() *gin.Engine {
	r := gin.Default()

	// GET /tasks - List all tasks
	r.GET("/tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, tasks)
	})

	// POST /tasks - Create a new task
	r.POST("/tasks", func(c *gin.Context) {

		var newTask Task

		// Bind JSON
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		// Simple validation
		if newTask.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "title is required",
			})
			return
		}

		// Generate ID
		lastID++
		newTask.ID = lastID

		// Append task
		tasks = append(tasks, newTask)

		// Log task creation
		log.Printf("Task created: %+v\n", newTask)

		// Response
		c.JSON(http.StatusCreated, gin.H{
			"message": "task created successfully",
			"task":    newTask,
		})
	})

	// PUT /tasks/:id - Update a task
	r.PUT("/tasks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var updatedTask Task
		if err := c.ShouldBindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i, t := range tasks {
			if t.ID == id {
				tasks[i].Title = updatedTask.Title
				tasks[i].Completed = updatedTask.Completed
				c.JSON(http.StatusOK, tasks[i])
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
	})

	// DELETE /tasks/:id - Delete a task
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		for i, t := range tasks {
			if t.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
