package godoist

type Comment struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	Content   string `json:"content"`
	PostedAt  string `json:"posted_at"`
	ProjectID string `json:"project_id"`
}

// GetComments retrieves all comments for a task
func (t *Task) GetComments() ([]Comment, error) {
	return t.manager.api.GetComments(t.ID)
}

// CreateComment creates a comment for a task
func (api *TodoistAPI) CreateComment(taskID, content string) (*Comment, error) {
	payload := map[string]interface{}{
		"task_id": taskID,
		"content": content,
	}

	var comment Comment
	err := api.doPost("/comments", payload, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetComments retrieves all comments for a task
func (api *TodoistAPI) GetComments(taskID string) ([]Comment, error) {
	var comments []Comment
	err := api.doGetPaginated("/comments?task_id="+taskID, &comments)
	return comments, err
}

// UpdateComment updates a comment by its ID
func (api *TodoistAPI) UpdateComment(commentID, content string) error {
	payload := map[string]interface{}{
		"content": content,
	}
	return api.doPost("/comments/"+commentID, payload, nil)
}

// DeleteComment deletes a comment by its ID
func (api *TodoistAPI) DeleteComment(commentID string) error {
	return api.doDelete("/comments/" + commentID)
}
