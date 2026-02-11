package godoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var (
	APIURL = "https://api.todoist.com/rest/v2"
)

type TodoistAPI struct {
	Token  string
	logger *slog.Logger
}

// NewDispatcher creates a new Todoist API client
func NewDispatcher(token string) *TodoistAPI {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &TodoistAPI{Token: token, logger: logger}
}

func (t *TodoistAPI) doGet(path string, result interface{}) error {
	req, err := http.NewRequest("GET", APIURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error %s: %s", resp.Status, string(body))
	}

	return json.Unmarshal(body, result)
}

func (t *TodoistAPI) doPost(path string, payload interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", APIURL+path, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error %s: %s", resp.Status, string(body))
	}

	if result != nil && len(body) > 0 {
		return json.Unmarshal(body, result)
	}
	return nil
}

func (t *TodoistAPI) doPostNoBody(path string) error {
	req, err := http.NewRequest("POST", APIURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+t.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %s: %s", resp.Status, string(body))
	}
	return nil
}

func (t *TodoistAPI) GetTasks() ([]Task, error) {
	var tasks []Task
	err := t.doGet("/tasks", &tasks)
	return tasks, err
}

func (t *TodoistAPI) GetProjects() ([]Project, error) {
	var projects []Project
	err := t.doGet("/projects", &projects)
	return projects, err
}

func (t *TodoistAPI) CreateTask(fields map[string]interface{}) (*Task, error) {
	var task Task
	err := t.doPost("/tasks", fields, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *TodoistAPI) UpdateTask(id string, fields map[string]interface{}) error {
	return t.doPost("/tasks/"+id, fields, nil)
}

func (t *TodoistAPI) CloseTask(id string) error {
	return t.doPostNoBody("/tasks/" + id + "/close")
}

func (t *TodoistAPI) CreateProject(fields map[string]interface{}) (*Project, error) {
	var project Project
	err := t.doPost("/projects", fields, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (t *TodoistAPI) UpdateProject(id string, fields map[string]interface{}) error {
	return t.doPost("/projects/"+id, fields, nil)
}
