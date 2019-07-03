package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func testApp() *gin.Engine {
	db, _ = gorm.Open("sqlite3", ":memory:")
	db.CreateTable(&tasks{})

	gin.SetMode(gin.TestMode)

	r := gin.Default()
	setupRoutes(r)
	return r
}

func requestTasks(t *testing.T, r *gin.Engine) []map[string]interface{} {

	req, err := http.NewRequest(http.MethodGet, "api/v1/tasks", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var dat []map[string]interface{}

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &dat); err != nil {
		panic(err)
	}

	return dat
}

func createTask(t *testing.T, r *gin.Engine, title string, list string, resultStatus int) map[string]interface{} {

	task := map[string]string{"Title": title, "List": list}
	jsonStr, _ := json.Marshal(task)
	req, err := http.NewRequest(http.MethodPost, "api/v1/tasks", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var dat map[string]interface{}

	if w.Code != resultStatus {
		t.Fatalf("Expected to get status %d but instead got %d\n", resultStatus, w.Code)
	}

	return dat

}

func completeTask(t *testing.T, r *gin.Engine, title string, list string) int {

	task := map[string]string{"Title": title, "List": list}
	jsonStr, _ := json.Marshal(task)
	req, err := http.NewRequest(http.MethodDelete, "api/v1/tasks", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

func TestGetTasks(t *testing.T) {

	r := testApp()

	dat := requestTasks(t, r)

	if len(dat) > 0 {
		t.Fatalf("Expected to have no tasks in the result. Got %d\n", len(dat))
	}
}

func TestPostTask(t *testing.T) {

	r := testApp()

	createTask(t, r, "Task Title A", "list-name", http.StatusCreated)
	request := requestTasks(t, r)

	if len(request) != 1 {
		t.Fatalf("Expected to have one task in the result. Got %d\n", len(request))
	}

}

func TestPostDuplicateTask(t *testing.T) {

	r := testApp()

	createTask(t, r, "Task Title A", "list-name", http.StatusCreated)
	createTask(t, r, "Task Title A", "list-name", http.StatusConflict)
	request := requestTasks(t, r)

	if len(request) != 1 {
		t.Fatalf("Expected to have one task in the result. Got %d\n", len(request))
	}

}
func TestPostMultipleTasks(t *testing.T) {

	r := testApp()

	createTask(t, r, "Task Title A", "list-name", http.StatusCreated)
	createTask(t, r, "Task Title B", "list-name", http.StatusCreated)
	request := requestTasks(t, r)

	if len(request) != 2 {
		t.Fatalf("Expected to have one task in the result. Got %d\n", len(request))
	}

}

func TestDeleteTask(t *testing.T) {
	r := testApp()

	title := "Task Title"
	list := "list-name"

	createTask(t, r, title, list, http.StatusCreated)

	responseCode := completeTask(t, r, title, list)

	if responseCode != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, responseCode)
	}

	request := requestTasks(t, r)

	if len(request) != 0 {
		t.Fatalf("Expected to have one task in the result. Got %d\n", len(request))
	}

}
