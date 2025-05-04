package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetProjects retrieves all projects
func (c *Client) GetProjects() ([]Project, error) {
	endpoint := "/projects"

	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	bodyBytes, err := c.processResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse response
	var projects []Project
	if err := json.Unmarshal(bodyBytes, &projects); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return projects, nil
}

// GetProject retrieves a specific project by ID
func (c *Client) GetProject(id string) (*Project, error) {
	endpoint := fmt.Sprintf("/projects/%s", id)

	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	bodyBytes, err := c.processResponse(resp, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// Parse response
	var project Project
	if err := json.Unmarshal(bodyBytes, &project); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &project, nil
}
