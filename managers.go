package godoist

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Manager struct {
	Tasks    *TaskManager
	Projects *ProjectManager
}

type TaskManager struct {
	api     *TodoistAPI
	tasks   map[string]*Task
	Manager *Manager
}

func NewTaskManager(api *TodoistAPI) *TaskManager {
	return &TaskManager{api: api, tasks: make(map[string]*Task)}
}

func (t *TaskManager) addTask(task Task) {
	t.tasks[task.ID] = &task
}

func (t *TaskManager) All() []*Task {
	var tasks = make([]*Task, 0, len(t.tasks))
	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (t *TaskManager) Update(tasks []Task) {
	for _, task := range tasks {
		task.manager = t
		t.addTask(task)
	}
}

func (t *TaskManager) Get(id string) *Task {
	task, exists := t.tasks[id]
	if !exists {
		return nil
	}
	return task
}

func (t *TaskManager) GetByName(name string) []*Task {
	var tasks = make([]*Task, 0)
	for _, task := range t.tasks {
		if task.Content == name {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (t *TaskManager) String() string {
	return "TaskManager"
}

func (t *TaskManager) Len() int {
	return len(t.tasks)
}

func (t *TaskManager) UpdateTask(task Task) {
	t.tasks[task.ID] = &task
}

func (t *TaskManager) AddTask(task Task) error {
	if _, exists := t.tasks[task.ID]; exists {
		return fmt.Errorf("Task with ID %s already exists", task.ID)
	}
	if task.ID == "" {
		task.TempID = uuid.New().String()
		taskJSON, err := json.Marshal(task)
		if err != nil {
			return err
		}
		var taskMap map[string]interface{}
		if err := json.Unmarshal(taskJSON, &taskMap); err != nil {
			return err
		}
		_, err = json.Marshal(task)
		if err != nil {
			return err
		}
		t.api.create("item_add", cleanRequests(taskMap), task.TempID)
	}
	task.manager = t
	t.tasks[task.ID] = &task
	return nil
}

func (t *TaskManager) Create(content string) *Task {
	task := Task{Content: content, manager: t}
	t.AddTask(task)
	return &task
}

func cleanRequests(requests map[string]interface{}) map[string]interface{} {
	for key, value := range requests {
		if value == nil {
			delete(requests, key)
		}
		if key == "TempID" {
			delete(requests, key)
		}
		if value == "" {
			delete(requests, key)
		}
		if value == 0 {
			delete(requests, key)
		}
		if value == 0.0 {
			delete(requests, key)
		}
	}
	return requests
}

type ProjectManager struct {
	api      *TodoistAPI
	projects map[string]*Project
	Manager  *Manager
}

func NewProjectManager(api *TodoistAPI) *ProjectManager {
	return &ProjectManager{api: api, projects: make(map[string]*Project)}
}

func (p *ProjectManager) AddProject(project Project) {
	project.Manager = p
	p.projects[project.ID] = &project
}

func (p *ProjectManager) Update(projects []Project) {
	for _, project := range projects {
		project.Manager = p
		p.projects[project.ID] = &project
	}
}

func (p *ProjectManager) All() []*Project {
	var projects = make([]*Project, 0, len(p.projects))
	for _, project := range p.projects {
		projects = append(projects, project)
	}
	return projects
}

func (p *ProjectManager) Get(id string) *Project {
	project, exists := p.projects[id]
	if !exists {
		return nil
	}
	return project
}

func (p *ProjectManager) GetByName(name string) []*Project {
	var projects = make([]*Project, 0)
	for _, project := range p.projects {
		if project.Name == name {
			projects = append(projects, project)
		}
	}
	return projects
}
