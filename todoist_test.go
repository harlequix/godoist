package godoist

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSync(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tasks := []Task{
			{ID: "1", Content: "Buy milk", ProjectID: "100", Order: 1, Priority: LOW},
			{ID: "2", Content: "Write tests", ProjectID: "100", Order: 2, Priority: HIGH, Labels: []string{"dev"}},
		}
		json.NewEncoder(w).Encode(tasks)
	})
	mux.HandleFunc("GET /projects", func(w http.ResponseWriter, r *http.Request) {
		projects := []Project{
			{ID: "100", Name: "Inbox", IsInboxProject: true, Order: 0},
			{ID: "200", Name: "Work", Order: 1},
		}
		json.NewEncoder(w).Encode(projects)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	orig := APIURL
	APIURL = srv.URL
	defer func() { APIURL = orig }()

	td := NewTodoist("test-token")
	if err := td.Sync(); err != nil {
		t.Fatalf("Sync() returned error: %v", err)
	}

	// Verify tasks
	if td.Tasks.Len() != 2 {
		t.Fatalf("expected 2 tasks, got %d", td.Tasks.Len())
	}
	task := td.Tasks.Get("1")
	if task == nil {
		t.Fatal("task '1' not found")
	}
	if task.Content != "Buy milk" {
		t.Errorf("expected 'Buy milk', got %q", task.Content)
	}
	if task.Priority != LOW {
		t.Errorf("expected priority LOW, got %v", task.Priority)
	}

	task2 := td.Tasks.Get("2")
	if task2 == nil {
		t.Fatal("task '2' not found")
	}
	if len(task2.Labels) != 1 || task2.Labels[0] != "dev" {
		t.Errorf("expected labels [dev], got %v", task2.Labels)
	}

	// Verify projects
	projects := td.Projects.All()
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
	inbox := td.Projects.Get("100")
	if inbox == nil {
		t.Fatal("project '100' not found")
	}
	if inbox.Name != "Inbox" {
		t.Errorf("expected 'Inbox', got %q", inbox.Name)
	}
	if !inbox.IsInboxProject {
		t.Error("expected IsInboxProject to be true")
	}
}
