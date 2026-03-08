package godoist

// MoveTask moves a task to a different project and/or parent.
func (t *TodoistAPI) MoveTask(taskID, projectID, parentID string) error {
	fields := map[string]interface{}{
		"project_id": projectID,
	}
	if parentID != "" {
		fields["parent_id"] = parentID
	}
	return t.doPost("/tasks/"+taskID+"/move", fields, nil)
}

// DeleteTask deletes a task by ID.
func (t *TodoistAPI) DeleteTask(id string) error {
	return t.doDelete("/tasks/" + id)
}
