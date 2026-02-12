package godoist

import (
	"log/slog"
	"os"
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
	tasks, err := t.API.GetTasks()
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	t.Tasks.Update(tasks)

	projects, err := t.API.GetProjects()
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	t.Projects.Update(projects)

	return nil
}

// Commit is a no-op kept for backwards compatibility.
// The API v1 executes operations immediately; there is nothing to commit.
func (t *Todoist) Commit() error {
	t.logger.Warn("Commit() is deprecated: API v1 executes operations immediately")
	return nil
}
