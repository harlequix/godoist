package godoist

import (
	"log/slog"
	"os"
)

type TaskManagerInterface interface {
    AddTask(task Task) error
    UpdateTask(task Task)
    Get(id string) *Task
    GetByName(name string) []*Task
    All() []*Task
}

type ProjectManagerInterface interface {
    AddProject(project Project)
    Update(projects []Project)
    Get(id string) *Project
    GetByName(name string) []*Project
    All() []*Project
}

type Todoist struct {
	Token    string
	logger   *slog.Logger
	API      *TodoistAPI
	Tasks    TaskManagerInterface
	Projects ProjectManagerInterface
}

// NewTodoist creates a new Todoist client
func NewTodoist(token string) *Todoist {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	aux := &Todoist{Token: token, logger: logger, API: NewDispatcher(token)}
	manager := Manager{}

	taskManager := NewTaskManager(aux.API)
	projectManager := NewProjectManager(aux.API)
	manager.Tasks = taskManager
    manager.Projects = projectManager
    taskManager.Manager = &manager
    projectManager.Manager = &manager

    aux.Tasks = taskManager
    aux.Projects = projectManager
	return aux
}

func (t *Todoist) Sync() error {
	response, err := t.API.Sync()
	if err != nil {
		t.logger.Error(err.Error())
		return err
	}
	for _, task := range response.Tasks {
		t.Tasks.AddTask(task)
	}
	t.Projects.Update(response.Projects)

	return nil
}
