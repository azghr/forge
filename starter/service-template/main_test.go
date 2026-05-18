package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	r := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	handleHealth(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %q", body["status"])
	}
}

func TestCreateTask(t *testing.T) {
	tasks = nil

	body := `{"id":"1","title":"test task"}`
	r := httptest.NewRequest("POST", "/api/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	handleCreateTask(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	var task Task
	if err := json.NewDecoder(w.Body).Decode(&task); err != nil {
		t.Fatal(err)
	}
	if task.Title != "test task" {
		t.Errorf("expected 'test task', got %q", task.Title)
	}
}

func TestCreateTask_emptyTitle(t *testing.T) {
	body := `{"id":"1","title":""}`
	r := httptest.NewRequest("POST", "/api/tasks", strings.NewReader(body))
	w := httptest.NewRecorder()
	handleCreateTask(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestCreateTask_badJSON(t *testing.T) {
	r := httptest.NewRequest("POST", "/api/tasks", strings.NewReader("not json"))
	w := httptest.NewRecorder()
	handleCreateTask(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListTasks(t *testing.T) {
	tasks = nil
	tasks = append(tasks, Task{ID: "1", Title: "existing"})

	r := httptest.NewRequest("GET", "/api/tasks", nil)
	w := httptest.NewRecorder()
	handleTasks(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var result []Task
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].Title != "existing" {
		t.Errorf("expected [existing], got %v", result)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"debug", "DEBUG"},
		{"info", "INFO"},
		{"warn", "WARN"},
		{"error", "ERROR"},
		{"unknown", "INFO"},
		{"", "INFO"},
	}
	for _, tt := range tests {
		got := parseLogLevel(tt.input)
		if got.String() != tt.want {
			t.Errorf("parseLogLevel(%q) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestConfigString(t *testing.T) {
	c := Config{Port: "9090", LogFormat: "json", LogLevel: "debug"}
	s := c.String()
	if s != "port=9090 log_format=json log_level=debug" {
		t.Errorf("unexpected config string: %s", s)
	}
}

func TestValidationError(t *testing.T) {
	err := &validationError{"test error"}
	if err.Error() != "test error" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestTaskValidate(t *testing.T) {
	if err := (Task{}).Validate(); err == nil {
		t.Error("expected error for empty task")
	}
	if err := (Task{Title: "ok"}).Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
