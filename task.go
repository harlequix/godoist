package godoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Deadline struct {
	Date       string    `json:"date"`
	Lang       string    `json:"lang"`
	ParsedDate time.Time `json:"-"`
}

type Due struct {
	Date        string    `json:"date"`
	Lang        string    `json:"lang"`
	String      string    `json:"string"`
	Timezone    string    `json:"timezone"`
	IsRecurring bool      `json:"is_recurring"`
	ParsedDate  time.Time `json:"-"`
}

func (d *Due) UnmarshalJSON(data []byte) error {
	type Alias Due
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if d.Date != "" {
		parsedDate, err := time.Parse("2006-01-02T15:04:05", d.Date)
		if err != nil {
			parsedDate, err := time.Parse("2006-01-02", d.Date)
			if err != nil {
				return err
			}
			d.ParsedDate = parsedDate
		} else {
			d.ParsedDate = parsedDate
		}
	}

	return nil
}

func (d *Deadline) UnmarshalJSON(data []byte) error {
	type Alias Deadline
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if d.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			return err
		}
		d.ParsedDate = parsedDate
	}

	return nil
}

type Task struct {
	ID        string         `json:"id"`
	Content   string         `json:"content"`
	ProjectID string         `json:"project_id"`
	SectionID string         `json:"section_id"`
	Order     int            `json:"order"`
	Priority  PRIORITY_LEVEL `json:"priority"`
	Deadline  *Deadline      `json:"deadline"`
	Due       *Due           `json:"due"`
	ParentID  string         `json:"parent_id"`
	Labels    []string       `json:"labels"`
	manager   *TaskManager   `json:"-"`
	TempID    string         `json:"-"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		Deadline *json.RawMessage `json:"deadline"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Deadline != nil {
		var deadline Deadline
		if err := json.Unmarshal(*aux.Deadline, &deadline); err != nil {
			return err
		}
		t.Deadline = &deadline
	} else {
		t.Deadline = nil
	}

	return nil
}

func (t *Task) GetChildren() []*Task {
	var tasks = make([]*Task, 0)
	for _, task := range t.manager.tasks {
		if task.ParentID == t.ID {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (t *Task) String() string {
	return t.Content
}

func (t *Task) AddLabel(label string) {
	for _, existingLabel := range t.Labels {
		if existingLabel == label {
			return
		}
	}
	new_label := append(t.Labels, label)
	t.Update("labels", new_label)
}

func (t *Task) RemoveLabel(label string) error {
	for i, existingLabel := range t.Labels {
		if existingLabel == label {
			new_labels := append(t.Labels[:i], t.Labels[i+1:]...)
			t.Update("labels", new_labels)
			return nil
		}
	}
	return fmt.Errorf("label not found: %s", label)
}

func (t *Task) Update(key string, value interface{}) error {
	switch key {
	case "content", "Content":
		t.Content = value.(string)
	case "project_id", "ProjectID":
		t.ProjectID = value.(string)
	case "section_id", "SectionID":
		t.SectionID = value.(string)
	case "order", "Order":
		t.Order = value.(int)
	case "priority", "Priority":
		t.Priority = value.(PRIORITY_LEVEL)
	case "deadline", "Deadline":
		t.Deadline = value.(*Deadline)
	case "due", "Due":
		t.Due = value.(*Due)
	case "parent_id", "ParentID":
		t.ParentID = value.(string)
	case "labels", "Labels":
		t.Labels = value.([]string)
	default:
		t.manager.api.logger.Error("Unknown/unsupported Update", "Command", key, "Task", t)
		return errors.New("unknown/unsupported Update")
	}
	t.manager.api.update("item_update", map[string]interface{}{"id": t.ID, key: value})
	return nil
}
