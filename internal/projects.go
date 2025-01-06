package internal

type Project struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Color      string          `json:"color"`
	ParentID   string          `json:"parent_id"`
	ChildOrder int             `json:"child_order"`
	Shared     bool            `json:"shared"`
	ViewStyle  string          `json:"view_style"`
	Manager    *ProjectManager `json:"-"`
}

func (p *Project) GetTasks() []*Task {
	if p.Manager == nil {
		return nil
	}

	tasks := []*Task{}
	for _, task := range p.Manager.Manager.Tasks.All() {
		if task.ProjectID == p.ID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (p *Project) GetChildren() []*Project {
	if p.Manager == nil {
		return nil
	}

	projects := []*Project{}
	for _, project := range p.Manager.All() {
		if project.ParentID == p.ID {
			projects = append(projects, project)
		}
	}
	return projects
}
