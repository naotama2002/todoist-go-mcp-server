package todoist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetTasks retrieves all active tasks
func (c *Client) GetTasks(ctx context.Context, projectID, filter string) ([]Task, error) {
	endpoint := "/tasks"

	// Add query parameters if provided
	if projectID != "" || filter != "" {
		// Create URL for query parameter handling
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse endpoint: %w", err)
		}

		// Use Query() to properly handle URL encoding
		q := u.Query()
		if projectID != "" {
			q.Add("project_id", projectID)
		}
		if filter != "" {
			q.Add("filter", filter)
		}
		u.RawQuery = q.Encode()

		endpoint = u.String()
	}

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	bodyBytes, err := c.processResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse response
	var tasks []Task
	if err := json.Unmarshal(bodyBytes, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return tasks, nil
}

// GetTask retrieves a specific task by ID
func (c *Client) GetTask(ctx context.Context, id string) (*Task, error) {
	endpoint := fmt.Sprintf("/tasks/%s", id)

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	bodyBytes, err := c.processResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse response
	var task Task
	if err := json.Unmarshal(bodyBytes, &task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// CreateTask creates a new task
func (c *Client) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
	endpoint := "/tasks"

	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	bodyBytes, err := c.processResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse response
	var task Task
	if err := json.Unmarshal(bodyBytes, &task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// UpdateTask updates an existing task
func (c *Client) UpdateTask(ctx context.Context, id string, req UpdateTaskRequest) (*Task, error) {
	endpoint := fmt.Sprintf("/tasks/%s", id)

	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	bodyBytes, err := c.processResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse response
	var task Task
	if err := json.Unmarshal(bodyBytes, &task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// CloseTask marks a task as completed
func (c *Client) CloseTask(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/tasks/%s/close", id)

	resp, err := c.doRequest(ctx, "POST", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to close task: %w", err)
	}

	_, err = c.processResponse(resp, http.StatusNoContent)
	return err
}

// ReopenTask marks a task as not completed
func (c *Client) ReopenTask(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/tasks/%s/reopen", id)

	resp, err := c.doRequest(ctx, "POST", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to reopen task: %w", err)
	}

	_, err = c.processResponse(resp, http.StatusNoContent)
	return err
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/tasks/%s", id)

	resp, err := c.doRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	_, err = c.processResponse(resp, http.StatusNoContent)
	return err
}
