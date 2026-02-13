package godoist

import "errors"

type Project struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Color        string          `json:"color"`
	ParentID     string          `json:"parent_id"`
	ChildOrder   int             `json:"child_order"`
	IsShared     bool            `json:"is_shared"`
	InboxProject bool            `json:"inbox_project"`
	IsFavorite   bool            `json:"is_favorite"`
	IsArchived   bool            `json:"is_archived"`
	IsCollapsed  bool            `json:"is_collapsed"`
	ViewStyle    string          `json:"view_style"`
	DefaultOrder int             `json:"default_order"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
	URL          string          `json:"url"`
	Manager      *ProjectManager `json:"-"`
}

func (p *Project) Update(key string, value interface{}) error {
	switch key {
	case "name", "Name":
		p.Name = value.(string)
	case "description", "Description":
		p.Description = value.(string)
	case "color", "Color":
		p.Color = value.(string)
	case "is_favorite", "IsFavorite":
		p.IsFavorite = value.(bool)
	case "view_style", "ViewStyle":
		p.ViewStyle = value.(string)
	default:
		return errors.New("unknown/unsupported Update")
	}
	return p.Manager.api.UpdateProject(p.ID, map[string]interface{}{key: value})
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
