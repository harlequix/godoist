package godoist

import (
	"log/slog"
	"os"

	"github.com/harlequix/godoist/internal"
)

type Todoist struct {
	Token    string
	logger   *slog.Logger
	API      *internal.TodoistAPI
	Tasks    internal.TaskManager
	Projects internal.ProjectManager
}

// NewTodoist creates a new Todoist client
func NewTodoist(token string) *Todoist {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	aux := &Todoist{Token: token, logger: logger, API: internal.NewDispatcher(token)}
	manager := internal.Manager{}

	aux.Tasks = *internal.NewTaskManager(aux.API)
	aux.Projects = *internal.NewProjectManager(aux.API)
	manager.Tasks = &aux.Tasks
	manager.Projects = &aux.Projects
	aux.Tasks.Manager = &manager
	aux.Projects.Manager = &manager
	return aux
}

func (t *Todoist) Sync() error {
	response, err := t.API.Sync()
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	t.Tasks.Update(response.Tasks)
	t.Projects.Update(response.Projects)

	return nil
}
