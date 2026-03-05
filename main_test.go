package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseTasks []Task
	err := json.Unmarshal(w.Body.Bytes(), &responseTasks)
	assert.NoError(t, err)
	assert.NotEmpty(t, responseTasks)
}

func TestPostTask(t *testing.T) {
	router := setupRouter()

	newTask := Task{Title: "Test Task", Completed: false}
	jsonValue, _ := json.Marshal(newTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Message string `json:"message"`
		Task    Task   `json:"task"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test task created successfully", response.Message)
	assert.Equal(t, "Test Task", response.Task.Title)
	assert.Equal(t, false, response.Task.Completed)
}

func TestUpdateTask(t *testing.T) {
	router := setupRouter()

	updatedTask := Task{Title: "Updated Task", Completed: true}
	jsonValue, _ := json.Marshal(updatedTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var task Task
	err := json.Unmarshal(w.Body.Bytes(), &task)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Task", task.Title)
	assert.Equal(t, true, task.Completed)
}

func TestDeleteTask(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task deleted")
}

func TestNotFound(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
