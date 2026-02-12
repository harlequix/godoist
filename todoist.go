package godoist

import (
	"log/slog"
	"os"
	"sync"
)

type Todoist struct {
	Token    string
	logger   *slog.Logger
	API      *TodoistAPI
	Tasks    TaskManager
	Projects ProjectManager
}

// NewTodoist creates a new Todoist client
func NewTodoist(token string) *Todoist {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	aux := &Todoist{Token: token, logger: logger, API: NewDispatcher(token)}
	manager := Manager{}

	aux.Tasks = *NewTaskManager(aux.API)
	aux.Projects = *NewProjectManager(aux.API)
	manager.Tasks = &aux.Tasks
	manager.Projects = &aux.Projects
	aux.Tasks.Manager = &manager
	aux.Projects.Manager = &manager
	return aux
}

func (t *Todoist) Sync() error {
	var (
		tasks    []Task
		projects []Project
		taskErr  error
		projErr  error
		wg       sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		tasks, taskErr = t.API.GetTasks()
	}()
	go func() {
		defer wg.Done()
		projects, projErr = t.API.GetProjects()
	}()
	wg.Wait()

	if taskErr != nil {
		t.logger.Error(taskErr.Error())
		return taskErr
	}
	if projErr != nil {
		t.logger.Error(projErr.Error())
		return projErr
	}

	t.Tasks.Update(tasks)
	t.Projects.Update(projects)
	return nil
}

// Commit is a no-op kept for backwards compatibility.
// The API v1 executes operations immediately; there is nothing to commit.
func (t *Todoist) Commit() error {
	t.logger.Warn("Commit() is deprecated: API v1 executes operations immediately")
	return nil
}
