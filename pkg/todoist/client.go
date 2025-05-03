package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// TodoistAPIBaseURL is the base URL for Todoist API
	TodoistAPIBaseURL = "https://api.todoist.com/rest/v2"
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 10 * time.Second
)

// Client represents a Todoist API client
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
	logger     *logrus.Logger
}

// NewClient creates a new Todoist API client
func NewClient(token string, logger *logrus.Logger) *Client {
	if token == "" {
		token = os.Getenv("TODOIST_API_TOKEN")
	}

	if logger == nil {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		token:   token,
		baseURL: TodoistAPIBaseURL,
		logger:  logger,
	}
}

// Task represents a Todoist task
type Task struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id"`
	ParentID    string `json:"parent_id,omitempty"`
	Priority    int    `json:"priority"`
	Due         *Due   `json:"due,omitempty"`
	URL         string `json:"url"`
}

// Due represents a due date for a task
type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime,omitempty"`
	String      string `json:"string,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

// Project represents a Todoist project
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	Order          int    `json:"order"`
	Color          string `json:"color"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"`
	URL            string `json:"url"`
	ParentID       string `json:"parent_id,omitempty"`
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	Content     string `json:"content"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	Order       int    `json:"order,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	DueDatetime string `json:"due_datetime,omitempty"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Content     string `json:"content,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	DueDatetime string `json:"due_datetime,omitempty"`
}

// GetTasks retrieves all active tasks
func (c *Client) GetTasks(projectID, filter string) ([]Task, error) {
	endpoint := "/tasks"

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters if provided
	q := req.URL.Query()
	if projectID != "" {
		q.Add("project_id", projectID)
	}
	if filter != "" {
		q.Add("filter", filter)
	}
	req.URL.RawQuery = q.Encode()

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return tasks, nil
}

// GetTask retrieves a specific task by ID
func (c *Client) GetTask(id string) (*Task, error) {
	endpoint := fmt.Sprintf("/tasks/%s", id)

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// GetProjects retrieves all projects
func (c *Client) GetProjects() ([]Project, error) {
	endpoint := "/projects"

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return projects, nil
}

// GetProject retrieves a specific project by ID
func (c *Client) GetProject(id string) (*Project, error) {
	endpoint := fmt.Sprintf("/projects/%s", id)

	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &project, nil
}

// CreateTask creates a new task
func (c *Client) CreateTask(req CreateTaskRequest) (*Task, error) {
	endpoint := "/tasks"

	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	httpReq, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	httpReq.Header.Add("Authorization", "Bearer "+c.token)
	httpReq.Header.Add("Content-Type", "application/json")

	// Execute the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// UpdateTask updates an existing task
func (c *Client) UpdateTask(id string, req UpdateTaskRequest) (*Task, error) {
	endpoint := fmt.Sprintf("/tasks/%s", id)

	// Convert request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	httpReq, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	httpReq.Header.Add("Authorization", "Bearer "+c.token)
	httpReq.Header.Add("Content-Type", "application/json")

	// Execute the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// CloseTask marks a task as completed
func (c *Client) CloseTask(id string) error {
	endpoint := fmt.Sprintf("/tasks/%s/close", id)

	// Create request
	req, err := http.NewRequest("POST", c.baseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// ReopenTask marks a task as not completed
func (c *Client) ReopenTask(id string) error {
	endpoint := fmt.Sprintf("/tasks/%s/reopen", id)

	// Create request
	req, err := http.NewRequest("POST", c.baseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(id string) error {
	endpoint := fmt.Sprintf("/tasks/%s", id)

	// Create request
	req, err := http.NewRequest("DELETE", c.baseURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+c.token)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Printf("Error closing response body: %v", err)
		}
	}()

	// Check response status
	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
