package godoist

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMoveTask(t *testing.T) {
	var (
		receivedPath string
		receivedBody map[string]interface{}
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "bad content type", http.StatusBadRequest)
			return
		}
		receivedPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedBody)
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	orig := APIURL
	APIURL = srv.URL
	defer func() { APIURL = orig }()

	api := NewDispatcher("test-token")

	t.Run("with project only", func(t *testing.T) {
		err := api.MoveTask("task-1", "proj-2", "")
		if err != nil {
			t.Fatalf("MoveTask() returned error: %v", err)
		}
		if receivedPath != "/tasks/task-1/move" {
			t.Errorf("expected path '/tasks/task-1/move', got %q", receivedPath)
		}
		if receivedBody["project_id"] != "proj-2" {
			t.Errorf("expected project_id 'proj-2', got %v", receivedBody["project_id"])
		}
		if _, ok := receivedBody["parent_id"]; ok {
			t.Error("expected parent_id to be absent when empty")
		}
	})

	t.Run("with parent", func(t *testing.T) {
		err := api.MoveTask("task-1", "proj-2", "parent-3")
		if err != nil {
			t.Fatalf("MoveTask() returned error: %v", err)
		}
		if receivedBody["parent_id"] != "parent-3" {
			t.Errorf("expected parent_id 'parent-3', got %v", receivedBody["parent_id"])
		}
	})
}

func TestMoveTaskAPIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks/{id}/move", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	orig := APIURL
	APIURL = srv.URL
	defer func() { APIURL = orig }()

	api := NewDispatcher("test-token")
	err := api.MoveTask("task-1", "proj-2", "")
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestDeleteTask(t *testing.T) {
	var deletedPath string

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		deletedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	orig := APIURL
	APIURL = srv.URL
	defer func() { APIURL = orig }()

	api := NewDispatcher("test-token")
	err := api.DeleteTask("task-99")
	if err != nil {
		t.Fatalf("DeleteTask() returned error: %v", err)
	}
	if deletedPath != "/tasks/task-99" {
		t.Errorf("expected path '/tasks/task-99', got %q", deletedPath)
	}
}
