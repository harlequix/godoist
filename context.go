package godoist

import (
	"encoding/json"
	"fmt"
	"strings"
)

const ContextPrefix = "[CONTEXT]"

// getContextComment retrieves the existing context comment for a task, if any
func (t *Task) getContextComment() (*Comment, error) {
	comments, err := t.manager.api.GetComments(t.ID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if strings.HasPrefix(comment.Content, ContextPrefix) {
			return &comment, nil
		}
	}
	return nil, nil
}

// GetContext retrieves the context data for a task
func (t *Task) GetContext() (map[string]interface{}, error) {
	comment, err := t.getContextComment()
	if err != nil {
		return nil, err
	}

	if comment == nil {
		return make(map[string]interface{}), nil
	}

	contextJSON := strings.TrimPrefix(comment.Content, ContextPrefix+" ")
	var contextData map[string]interface{}
	if err := json.Unmarshal([]byte(contextJSON), &contextData); err != nil {
		return nil, fmt.Errorf("failed to parse context: %w", err)
	}

	return contextData, nil
}

// SetContext sets or updates the context data for a task
func (t *Task) SetContext(contextData map[string]interface{}) error {
	contextJSON, err := json.Marshal(contextData)
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	content := fmt.Sprintf("%s %s", ContextPrefix, string(contextJSON))

	existingComment, err := t.getContextComment()
	if err != nil {
		return err
	}

	if existingComment != nil {
		return t.manager.api.UpdateComment(existingComment.ID, content)
	}

	_, err = t.manager.api.CreateComment(t.ID, content)
	return err
}

// UpdateContext updates specific fields in the context without replacing everything
func (t *Task) UpdateContext(updates map[string]interface{}) error {
	currentContext, err := t.GetContext()
	if err != nil {
		return err
	}

	for key, value := range updates {
		currentContext[key] = value
	}

	return t.SetContext(currentContext)
}

// DeleteContext removes the context comment entirely
func (t *Task) DeleteContext() error {
	comment, err := t.getContextComment()
	if err != nil {
		return err
	}

	if comment == nil {
		return nil
	}

	return t.manager.api.DeleteComment(comment.ID)
}

// DeleteContextField removes a specific field from the context
func (t *Task) DeleteContextField(key string) error {
	currentContext, err := t.GetContext()
	if err != nil {
		return err
	}

	delete(currentContext, key)

	if len(currentContext) == 0 {
		return t.DeleteContext()
	}

	return t.SetContext(currentContext)
}
